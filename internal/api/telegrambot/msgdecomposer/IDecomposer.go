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

type DeferredCallback func(callback ISender) (*gotgbot.Message, error)

func (d *MessageDecomposer) Decompose(msg merger.Message, sender ISender) error {

	for {
		callbacks := make([]DeferredCallback, 0)

		if len(msg.Media) == 1 && len(callbacks) == 0 {
			singleMediaMsg, poorMsg := pullSingleMedia(msg)
			msg = poorMsg
			callbacks = append(callbacks, sendSingleMedia(singleMediaMsg))

		} else if len(msg.Media) > 1 {
			mediaMsgs, poorMsg := pullMediaGroups(msg)
			msg = poorMsg
			callbacks = append(callbacks, sendSingleMedia(singleMediaMsg))

		} else if len(msg.Forwarded) > 0 {
			// todo if reply + text -=> send text+repl
			// todo if text -> sendTexted
			// todo final send each forward
		} else if msg.Text != nil && *msg.Text != "" {

		} else if msg.ReplyId != nil {

		} else if msg.Username != nil {

		} else {
			return nil
		}
	}
}

func pullMediaGroups(msg merger.Message) (extracted map[merger.MediaType][]merger.Message, poor merger.Message) {
	poor = msg
	poor.Media = nil

	for _, media := range msg.Media {
		msgs, ok := extracted[media.Kind]
		if !ok {
			msgs = make([]merger.Message, 0)
		}
		ex := msg
		ex.Media = []merger.Media{media}
		ex.Forwarded = nil
		msgs = append(msgs, ex)
		extracted[media.Kind] = msgs
	}
}

func sendMediaGroups(msgMap map[merger.MediaType][]merger.Message) []DeferredCallback {
	cbks := make([]DeferredCallback, 0)

	for kind, msgSlice := range msgMap {
		//isStickerGroup := false
		if kind == merger.Sticker {
			for _, sticker := range msgSlice {
				cbks = append(cbks, func(callback ISender) (*gotgbot.Message, error) {
					return callback.SendSticker(sticker)
				})
			}
		} else {
			cbks = append(cbks, func(callback ISender) (*gotgbot.Message, error) {
				group, err := callback.SendMediaGroup(msgSlice)
				if err != nil {
					return nil, fmt.Errorf("send media group: %s", err)
				}
				media := msgSlice[0].Media
				switch media {
				case merger.Audio:
					return callback.SendAudio(msgSlice)
				case merger.Video:
					return callback.SendVideo(msgSlice)
				case merger.File:
					return callback.SendDocument(msgSlice)
				case merger.Photo:
					return callback.SendPhoto(msgSlice)
				}
				return nil, ErrUnknownMediaType
			},
			)
		}
	}
	return cbks
}

func pullSingleMedia(msg merger.Message) (extracted merger.Message, poor merger.Message) {
	extracted = msg
	extracted.Media = extracted.Media[:1]
	extracted.Forwarded = nil

	poor = msg
	poor.Media = poor.Media[1:]

	return extracted, msg
}

func sendSingleMedia(msg merger.Message) DeferredCallback {
	return func(callback ISender) (*gotgbot.Message, error) {
		media := msg.Media[0]
		switch media.Kind {
		case merger.Audio:
			return callback.SendAudio(msg)
		case merger.Video:
			return callback.SendVideo(msg)
		case merger.File:
			return callback.SendDocument(msg)
		case merger.Photo:
			return callback.SendPhoto(msg)
		case merger.Sticker:
			return callback.SendSticker(msg)
		}
		return nil, ErrUnknownMediaType
	}
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
