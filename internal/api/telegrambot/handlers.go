package telegrambot

import (
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

func (c *Callback) SendMediaGroup(msgs []merger.Message) ([]gotgbot.Message, error) {
	c.bot.SendMediaGroup()
}

func (c *Callback) SendSticker(msg merger.Message) (*gotgbot.Message, error) {

}
func (c *Callback) SendPhoto(msg merger.Message) (*gotgbot.Message, error) {

	//if len(msg.Media) == 0 {
	//	return nil, errors.New("media is empty")
	//}
	photo := msg.Media[0]
	//if photo.Kind != merger.Photo {
	//	return nil, errors.New("media is not photo")
	//}
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
			Caption:         stringOrEmpty(msg.Text),
			ReplyParameters: replyParams,
			HasSpoiler:      photo.Spoiler,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("send message to tg: %s", err)
	}
	return tgMsg, nil
}

func (c *Callback) SendAudio(msg merger.Message) (*gotgbot.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Callback) SendVideo(msg merger.Message) (*gotgbot.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Callback) SendDocument(msg merger.Message) (*gotgbot.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Callback) SendForward(msg merger.Message) (*gotgbot.Message, error) {
	//TODO implement me
	panic("implement me")
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
