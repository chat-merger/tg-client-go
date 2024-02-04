package telegrambot

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"merger-adapter/internal/api/telegrambot/deffereduploader"
	mrepo "merger-adapter/internal/repository/messages_repository"
	"merger-adapter/internal/service/merger"
	"time"
)

var lastReq = time.Now()

func (c *Client) onTelegramCreatedMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Printf("onTelegramCreatedMessage: after last %d ms ", time.Now().UnixMilli()-lastReq.UnixMilli())
	lastReq = time.Now()
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
		go c.onMergerMessage(msg)
	}
}

func (c *Client) onMergerMessage(msg *merger.Message) {
	// reply
	replyTo, err := replyTgIdFromMsg(msg, c.repo, c.chatID)
	if err != nil {
		log.Printf("[ERROR] replyTgIdFromMsg: %s", err)
		replyTo = 0
	}

	tgMsg, err := c.bot.SendMessage(
		c.chatID,
		msg.FormatShort(),
		&gotgbot.SendMessageOpts{ReplyParameters: &gotgbot.ReplyParameters{ChatId: replyTo}},
	)
	if err != nil {
		log.Printf("[ERROR] send message to tg: %s", err)
		return
	}
	kind := deffereduploader.DefineKind(*tgMsg)
	err = c.repo.Add(mrepo.Message{
		ReplyMergerMsgId: msg.ReplyId,
		MergerMsgId:      msg.Id,
		ChatId:           tgMsg.Chat.Id,
		MsgId:            tgMsg.MessageId,
		SenderId:         tgMsg.GetSender().Id(),
		SenderFirstName:  tgMsg.GetSender().FirstName(),
		Kind:             mrepo.Kind(kind),
		HasMedia:         kind == deffereduploader.Media || kind == deffereduploader.GroupMedia,
		CreatedAt:        tgMsg.Date,
	})
	if err != nil {
		log.Printf("[ERROR] add msg to repos: %s", err)
		return
	}
}

func replyTgIdFromMsg(msg *merger.Message, repo mrepo.MessagesRepository, chatId int64) (int64, error) {
	if msg.ReplyId != nil {
		messages, err := repo.Get(mrepo.Filter{ChatId: &chatId, MergerMsgId: msg.ReplyId})
		if err != nil {
			return 0, fmt.Errorf("messages from repo: %s", err)
		}
		if len(messages) == 0 {
			return 0, nil
		}
		return messages[0].MsgId, nil
	}
	return 0, nil
}
