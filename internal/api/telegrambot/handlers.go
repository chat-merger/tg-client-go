package telegrambot

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"merger-adapter/internal/api/telegrambot/deffereduploader"
	mrepo "merger-adapter/internal/repository/messages_repository"
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

func onMergerMessage(msg *merger.Message, bot *gotgbot.Bot, repo mrepo.MessagesRepository, chatID int64) {
	// reply
	replyTo, err := replyTgIdFromMsg(msg, repo, chatID)
	if err != nil {
		log.Printf("[ERROR] replyTgIdFromMsg: %s", err)
		replyTo = 0
	}

	tgMsg, err := bot.SendMessage(
		chatID,
		msg.FormatShort(),
		&gotgbot.SendMessageOpts{ReplyParameters: &gotgbot.ReplyParameters{ChatId: replyTo}},
	)
	if err != nil {
		log.Printf("[ERROR] send message to tg: %s", err)
		return
	}
	err = saveToRepo(repo, msg, tgMsg)
	if err != nil {
		log.Printf("[ERROR] save message to repo: %s", err)
		return
	}
}

func saveToRepo(repo mrepo.MessagesRepository, msg *merger.Message, tgMsg *gotgbot.Message) error {
	kind := deffereduploader.DefineKind(*tgMsg)
	return repo.Add(mrepo.Message{
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
