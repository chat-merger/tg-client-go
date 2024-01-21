package telegrambot

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
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
	//debug.Print(ctx.Message)
	var replyedId *string
	if ctx.Message.ReplyToMessage != nil {
		id := strconv.FormatInt(ctx.Message.ReplyToMessage.MessageId, 10)
		replyedId = &id
	}
	author := ctx.EffectiveUser.Username

	msg := merger.CreateMessage{
		ReplyId: (*merger.ID)(replyedId),
		Date:    time.Unix(ctx.Message.Date, 0),
		Uername: &author,
		Silent:  false, // where prop??
		Body: &merger.BodyText{
			Format: merger.Plain,
			Value:  ctx.Message.Text,
		},
	}

	_, err := c.conn.Send(msg)
	if err != nil {
		return fmt.Errorf("send message to Server: %s", err)
	}
	return nil
}

func (c *Client) listenServerMessages() error {
	for {
		msg, err := c.conn.Update()
		if err != nil {
			return fmt.Errorf("receive update: %s", err)
		}
		_, err = c.bot.SendMessage(c.chatID, msg.FormatShort(), nil)
		if err != nil {
			return fmt.Errorf("send message to tg: %s", err)
		}
	}
}
