package vkontaktebot

import (
	"context"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"log"
	"merger-adapter/internal/debug"
	"merger-adapter/internal/service/merger"
	"strconv"
	"time"
)

func (c *Client) gotgbotSetup() {
	// New message event
	c.lp.MessageNew(c.onMessage)
}

func (c *Client) filter(msg *gotgbot.Message) bool {
	return msg.Chat.Id == c.chatID && msg.Text != ""
}

func (c *Client) onMessage(_ context.Context, obj events.MessageNewObject) {
	debug.Print(obj)
	var replyedId *string
	if obj.Message.ReplyMessage != nil {
		id := strconv.Itoa(obj.Message.ReplyMessage.ID)
		replyedId = &id
	}
	author := "vk user"

	msg := merger.CreateMessage{
		ReplyId: (*merger.ID)(replyedId),
		Date:    time.Unix(int64(obj.Message.Date), 0),
		Author:  &author,
		Silent:  bool(obj.Message.IsSilent),
		Body: &merger.BodyText{
			Format: "", // todo
			Value:  obj.Message.Text,
		},
	}

	err := c.conn.Send(msg)
	if err != nil {
		log.Fatalf("send message to Server: %s", err)
	}
}

func (c *Client) listenServerMessages() error {
	for {
		msg, err := c.conn.Update()
		if err != nil {
			return fmt.Errorf("receive update: %s", err)
		}
		b := params.NewMessagesSendBuilder()
		b.Message(fmt.Sprintf(
			"- %v\n- %v\n- %v\n- %v\n- %v\n- %v\n- %v\n",
			msg.Id,
			msg.ReplyId,
			msg.Date,
			msg.Author,
			msg.From,
			msg.Silent,
			msg.Body,
		))
		b.RandomID(0)
		b.PeerID(int(c.chatID))

		_, err = c.vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
	}
}
