package telegrambot

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"log"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
	"strconv"
	"time"
)

func (c *Client) gotgbotSetup() {
	c.dispatcher.AddHandler(handlers.NewMessage(c.filter, c.onMessage))
}

func (c *Client) filter(msg *gotgbot.Message) bool {
	return msg.Chat.Id == c.chatID && msg.Text != ""
}

func (c *Client) onMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	var replyTo *string
	if ctx.Message.ReplyToMessage != nil {
		id, exists, err := c.messagesMap.ByContID(mmScope, toContID(ctx.Message.ReplyToMessage.MessageId))
		if err != nil {
			log.Printf("[ERROR] msg from message map: %s", err)
		}
		if exists {
			replyTo = (*string)(id)
		}
	}
	author := ctx.EffectiveUser.Username

	msg := merger.CreateMessage{
		ReplyId: (*merger.ID)(replyTo),
		Date:    time.Unix(ctx.Message.Date, 0),
		Uername: &author,
		Silent:  false, // where prop??
		Body: &merger.BodyText{
			Format: merger.Plain,
			Value:  ctx.Message.Text,
		},
	}

	mMsg, err := c.conn.Send(msg)
	if err != nil {
		return fmt.Errorf("send message to Server: %s", err)
	}
	err = c.messagesMap.Save(mmScope, mMsg.Id, toContID(ctx.Message.MessageId))
	if err != nil {
		return fmt.Errorf("save msg id to MessageMap: %s", err)
	}
	return nil
}

const mmScope = kvstore.Scope("TgBotScope")

func toContID(id int64) kvstore.ContMsgID {
	return kvstore.ContMsgID(strconv.FormatInt(id, 10))
}

func toInt64(id kvstore.ContMsgID) int64 {
	vkMsgId, err := strconv.Atoi(string(id))
	if err != nil {
		log.Fatalf("[ERROR] convert kvstore.ContMsgID to int: %s", err)
	}
	return int64(vkMsgId)
}

func (c *Client) listenServerMessages() error {
	for {
		msg, err := c.conn.Update()
		if err != nil {
			return fmt.Errorf("receive update: %s", err)
		}

		// reply
		var replyTo int64
		if msg.ReplyId != nil {
			id, exists, err := c.messagesMap.ByMergedID(mmScope, *msg.ReplyId)
			if err != nil {
				log.Printf("[ERROR] msg from message map: %s", err)
			}
			if exists {
				replyTo = toInt64(*id)
			}
		}

		tgMsg, err := c.bot.SendMessage(
			c.chatID,
			msg.FormatShort(),
			&gotgbot.SendMessageOpts{ReplyToMessageId: replyTo},
		)
		if err != nil {
			return fmt.Errorf("send message to tg: %s", err)
		}
		err = c.messagesMap.Save(mmScope, msg.Id, toContID(tgMsg.MessageId))
		if err != nil {
			return fmt.Errorf("save msg id to MessageMap: %s", err)
		}
	}
}
