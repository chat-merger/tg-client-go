package msgdecomposer

import (
	"errors"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/api/telegrambot/tghelper"
	mrepo "merger-adapter/internal/repository/messages_repository"
	"merger-adapter/internal/service/merger"
)

type IMessageDecomposer interface {
	Decompose(msg merger.Message, callback ISender) error
}

type MessageDecomposer struct {
}

var (
	ErrEmptyMessage     = errors.New("empty message")
	ErrUnknownMediaType = errors.New("unknown media type")
)

//type DeferredCallback func(callback ISender) ([]gotgbot.Message, error)

type Pair struct {
	orig gotgbot.Message
	msg  merger.Message
}

func (d *MessageDecomposer) Decompose(msg merger.Message, sender ISender) error {

	for {
		pairs := make([]Pair, 0)

		if len(msg.Media) == 1 && len(pairs) == 0 {
			singleMediaMsg, poorMsg := pullSingleMedia(msg)
			msg = poorMsg
			out, err := sendSingleMedia(singleMediaMsg, sender)
			if err != nil {
				return fmt.Errorf("send single media: %s", err)
			}
			pairs = append(pairs, out)

		} else if len(msg.Media) > 1 {
			mediaMsgs, poorMsg := pullMediaGroups(msg)
			msg = poorMsg
			group, err := sendMediaGroups(mediaMsgs, sender)
			if err != nil {
				return fmt.Errorf("send media group: %s", err)
			}
			pairs = append(pairs, group...)

		} else if len(msg.Forwarded) > 0 {
			msg.Forwarded = nil
			// todo if reply + text -=> send text+repl
			// todo if text -> sendTexted
			// todo final send each forward

		} else if msg.Text != nil && *msg.Text != "" && len(pairs) == 0 {
			out, err := sender.SendTexted(msg)
			if err != nil {
				return fmt.Errorf("send texted: %s", err)
			}
			pairs = append(pairs, Pair{
				orig: *out,
				msg:  msg,
			})
			msg.Text = nil

		} else if msg.ReplyId != nil && len(pairs) == 0 {
			msg.ReplyId = nil
			// todo

		} else if msg.Username != nil && len(pairs) == 0 {
			msg.Username = nil
			// todo

		} else {
			return nil
		}
	}
}

func pullMediaGroups(msg merger.Message) (extracted map[merger.MediaType]merger.Message, poor merger.Message) {
	poor = msg
	poor.Media = nil

	for _, media := range msg.Media {
		ex, ok := extracted[media.Kind]
		if !ok {
			ex = msg
			ex.Forwarded = nil
			ex.Media = make([]merger.Media, 0)
			ex.Text = nil
			if len(extracted) == 0 && msg.Text != nil {
				ex.Text = msg.Text
			}
		}
		ex.Media = append(ex.Media, media)
		extracted[media.Kind] = ex
	}
	return extracted, poor
}

func sendMediaGroups(msgMap map[merger.MediaType]merger.Message, sender ISender) ([]Pair, error) {
	msgs := make([]Pair, 0)

	for kind, msg := range msgMap {
		//isStickerGroup := false
		if kind == merger.Sticker {
			orig, err := sender.SendSticker(msg)
			if err != nil {
				return nil, fmt.Errorf("send sticker: %s", err)
			}
			msgs = append(msgs, Pair{
				orig: *orig,
				msg:  msg,
			})
		} else {
			group, err := sender.SendMediaGroup(msg)
			if err != nil {
				return nil, fmt.Errorf("send media group: %s", err)
			}
			for _, origin := range group {
				msgs = append(msgs, Pair{
					orig: origin,
					msg:  msg,
				})
			}
		}
	}
	return msgs, nil
}

func pullSingleMedia(msg merger.Message) (extracted merger.Message, poor merger.Message) {
	extracted = msg
	extracted.Media = extracted.Media[:1]
	extracted.Forwarded = nil

	poor = msg
	poor.Media = poor.Media[1:]

	return extracted, poor
}

func sendSingleMedia(msg merger.Message, sender ISender) (Pair, error) {
	media := msg.Media[0]
	var orig *gotgbot.Message
	var err error
	switch media.Kind {
	case merger.Audio:
		orig, err = sender.SendAudio(msg)
	case merger.Video:
		orig, err = sender.SendVideo(msg)
	case merger.File:
		orig, err = sender.SendDocument(msg)
	case merger.Photo:
		orig, err = sender.SendPhoto(msg)
	case merger.Sticker:
		orig, err = sender.SendSticker(msg)
	default:
		return Pair{}, ErrUnknownMediaType
	}
	if err != nil {
		return Pair{}, fmt.Errorf("send some single file: %s", err)
	}
	return Pair{
		orig: *orig,
		msg:  msg,
	}, err
}

func saveToRepo(repo mrepo.MessagesRepository, msg *merger.Message, tgMsg *gotgbot.Message) error {
	kind := tghelper.DefineKind(*tgMsg)
	return repo.Add(mrepo.Message{
		ReplyMergerMsgId: msg.ReplyId,
		MergerMsgId:      msg.Id,
		ChatId:           tgMsg.Chat.Id,
		MsgId:            tgMsg.MessageId,
		SenderId:         tgMsg.GetSender().Id(),
		SenderFirstName:  tgMsg.GetSender().FirstName(),
		Kind:             kind,
		HasMedia:         kind.IsMedia(),
		CreatedAt:        tgMsg.Date,
	})
}
