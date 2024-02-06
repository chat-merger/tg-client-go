package telegrambot

import (
	"errors"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"merger-adapter/internal/api/telegrambot/msgdecomposer"
	mrepo "merger-adapter/internal/repository/messages_repository"
	"merger-adapter/internal/service/blobstore"
	"merger-adapter/internal/service/merger"
)

func (c *Client) onTelegramCreatedMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	err := c.du.Upload(*ctx.Message)
	if err != nil {
		return fmt.Errorf("upload msg to IDeferredUploader: %s", err)
	}
	return nil
}

func (c *Client) listenServerMessages() error {
	for {
		msg, err := c.conn.Update()
		if err != nil {
			return fmt.Errorf("receive update: %s", err)
		}
		go onMergerMessage(msg, c.bot, c.repo, c.chatID)
	}
}

type Callback struct {
	bot    *gotgbot.Bot
	chatId int64
	files  blobstore.TempBlobStore
	repo   mrepo.MessagesRepository
}

func (c *Callback) SendTexted(msg merger.Message) (*gotgbot.Message, error) {
	replyTo := replyParametersOrNil(msg, c.repo, c.chatId)
	tgMsg, err := c.bot.SendMessage(
		c.chatId,
		msg.FormatShort(),
		&gotgbot.SendMessageOpts{
			ReplyParameters:     replyTo,
			DisableNotification: msg.Silent,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("send message to tg: %s", err)
	}
	return tgMsg, nil
}

func (c *Callback) SendMediaGroup(msg merger.Message) ([]gotgbot.Message, error) {
	input := make([]gotgbot.InputMedia, 0)
	for _, media := range msg.Media {
		var b []byte
		reader, err := c.files.Get(media.Url)
		if err != nil {
			return nil, fmt.Errorf("get photo from files: %s", err)
		}
		_, err = reader.Read(b)
		if err != nil {
			return nil, fmt.Errorf("read bytes from reader: %s", err)
		}
		var imedia gotgbot.InputMedia
		switch media.Kind {
		case merger.Audio:
			imedia = gotgbot.InputMediaAudio{
				Media:   b,
				Caption: stringOrEmpty(msg.Text),
				Title:   media.Url,
			}
		case merger.Video:
			imedia = gotgbot.InputMediaVideo{
				Media:      b,
				Caption:    stringOrEmpty(msg.Text),
				HasSpoiler: media.Spoiler,
			}
		case merger.File:
			imedia = gotgbot.InputMediaDocument{
				Media:   b,
				Caption: stringOrEmpty(msg.Text),
			}
		case merger.Photo:
			imedia = gotgbot.InputMediaPhoto{
				Media:      b,
				Caption:    stringOrEmpty(msg.Text),
				HasSpoiler: media.Spoiler,
			}
		default:
			return nil, errors.New("unknown media type")
		}
		input = append(input, imedia)
	}
	replyParams := replyParametersOrNil(msg, c.repo, c.chatId)
	tgMsg, err := c.bot.SendMediaGroup(
		c.chatId,
		input,
		&gotgbot.SendMediaGroupOpts{
			ReplyParameters:     replyParams,
			DisableNotification: msg.Silent,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("send media group to tg: %s", err)
	}
	return tgMsg, nil
}

func (c *Callback) SendSticker(msg merger.Message) (*gotgbot.Message, error) {

	sticker := msg.Media[0]

	var b []byte
	reader, err := c.files.Get(sticker.Url)
	if err != nil {
		return nil, fmt.Errorf("get sticker from files: %s", err)
	}
	_, err = reader.Read(b)
	if err != nil {
		return nil, fmt.Errorf("read bytes from reader: %s", err)
	}
	replyParams := replyParametersOrNil(msg, c.repo, c.chatId)
	tgMsg, err := c.bot.SendSticker(
		c.chatId,
		b,
		&gotgbot.SendStickerOpts{
			ReplyParameters:     replyParams,
			DisableNotification: msg.Silent,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("send sticker to tg: %s", err)
	}
	return tgMsg, nil
}

func (c *Callback) SendPhoto(msg merger.Message) (*gotgbot.Message, error) {

	photo := msg.Media[0]

	var b []byte
	reader, err := c.files.Get(photo.Url)
	if err != nil {
		return nil, fmt.Errorf("get photo from files: %s", err)
	}
	_, err = reader.Read(b)
	if err != nil {
		return nil, fmt.Errorf("read bytes from reader: %s", err)
	}
	replyParams := replyParametersOrNil(msg, c.repo, c.chatId)
	tgMsg, err := c.bot.SendPhoto(
		c.chatId,
		b,
		&gotgbot.SendPhotoOpts{
			Caption:             stringOrEmpty(msg.Text),
			ReplyParameters:     replyParams,
			HasSpoiler:          photo.Spoiler,
			DisableNotification: msg.Silent,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("send photo to tg: %s", err)
	}
	return tgMsg, nil
}

func (c *Callback) SendAudio(msg merger.Message) (*gotgbot.Message, error) {
	audio := msg.Media[0]

	var b []byte
	reader, err := c.files.Get(audio.Url)
	if err != nil {
		return nil, fmt.Errorf("get audio from files: %s", err)
	}
	_, err = reader.Read(b)
	if err != nil {
		return nil, fmt.Errorf("read bytes from reader: %s", err)
	}
	replyParams := replyParametersOrNil(msg, c.repo, c.chatId)
	tgMsg, err := c.bot.SendAudio(
		c.chatId,
		b,
		&gotgbot.SendAudioOpts{
			Caption:             stringOrEmpty(msg.Text),
			ReplyParameters:     replyParams,
			DisableNotification: msg.Silent,
			Title:               audio.Url,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("send audio to tg: %s", err)
	}
	return tgMsg, nil
}

func (c *Callback) SendVideo(msg merger.Message) (*gotgbot.Message, error) {
	video := msg.Media[0]

	var b []byte
	reader, err := c.files.Get(video.Url)
	if err != nil {
		return nil, fmt.Errorf("get video from files: %s", err)
	}
	_, err = reader.Read(b)
	if err != nil {
		return nil, fmt.Errorf("read bytes from reader: %s", err)
	}
	replyParams := replyParametersOrNil(msg, c.repo, c.chatId)
	tgMsg, err := c.bot.SendVideo(
		c.chatId,
		b,
		&gotgbot.SendVideoOpts{
			Caption:             stringOrEmpty(msg.Text),
			ReplyParameters:     replyParams,
			DisableNotification: msg.Silent,
			HasSpoiler:          video.Spoiler,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("send video to tg: %s", err)
	}
	return tgMsg, nil
}

func (c *Callback) SendDocument(msg merger.Message) (*gotgbot.Message, error) {
	document := msg.Media[0]

	var b []byte
	reader, err := c.files.Get(document.Url)
	if err != nil {
		return nil, fmt.Errorf("get document from files: %s", err)
	}
	_, err = reader.Read(b)
	if err != nil {
		return nil, fmt.Errorf("read bytes from reader: %s", err)
	}
	replyParams := replyParametersOrNil(msg, c.repo, c.chatId)
	tgMsg, err := c.bot.SendDocument(
		c.chatId,
		b,
		&gotgbot.SendDocumentOpts{
			Caption:             stringOrEmpty(msg.Text),
			ReplyParameters:     replyParams,
			DisableNotification: msg.Silent,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("send document to tg: %s", err)
	}
	return tgMsg, nil
}

func (c *Callback) SendSingleForward(msg merger.Message) (*gotgbot.Message, error) {
	forward := msg.Forwarded[0]

	var b []byte
	reader, err := c.files.Get(forward.Url)
	if err != nil {
		return nil, fmt.Errorf("get forward from files: %s", err)
	}
	_, err = reader.Read(b)
	if err != nil {
		return nil, fmt.Errorf("read bytes from reader: %s", err)
	}
	replyParams := replyParametersOrNil(msg, c.repo, c.chatId)
	tgMsg, err := c.bot.ForwardMessage(
		c.chatId,
		b,
		&gotgbot.ForwardMessageOptss{
			Caption:             stringOrEmpty(msg.Text),
			ReplyParameters:     replyParams,
			DisableNotification: msg.Silent,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("send forward to tg: %s", err)
	}
	return tgMsg, nil
}

func onMergerMessage(msg *merger.Message, bot *gotgbot.Bot, repo mrepo.MessagesRepository, callback msgdecomposer.ISender, decomposer msgdecomposer.IMessageDecomposer, chatID int64) {
	err := decomposer.Decompose(*msg, callback)
	if err != nil {
		log.Printf("[ERROR] decompose msg: %s", err)
		return
	}
}

func replyParametersOrNil(msg merger.Message, repo mrepo.MessagesRepository, chatId int64) *gotgbot.ReplyParameters {
	if msg.ReplyId != nil {
		messages, err := repo.Get(mrepo.Filter{ChatId: &chatId, MergerMsgId: msg.ReplyId})
		if err != nil {
			log.Printf("[ERROR] messages from repo: %s", err)
			return nil
		}
		if len(messages) == 0 {
			return nil
		}
		return &gotgbot.ReplyParameters{
			MessageId: messages[0].MsgId,
		}
	}
	return nil
}

func stringOrEmpty(str *string) string {
	if str == nil {
		return ""
	} else {
		return *str
	}
}
