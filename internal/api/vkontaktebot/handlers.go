package vkontaktebot

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"log"
	"merger-adapter/internal/service/merger"
	"strconv"
	"time"
)

func (c *Client) gotgbotSetup() {
	// New message event
	c.lp.MessageNew(c.onMessage)
}

func (c *Client) onMessage(_ context.Context, obj events.MessageNewObject) {
	if obj.Message.PeerID != c.peerID {
		return
	}
	//debug.Print(obj)
	var replyedId *string
	if obj.Message.ReplyMessage != nil {
		id := strconv.Itoa(obj.Message.ReplyMessage.ID)
		replyedId = &id
	}
	var author *string

	usrs, _ := c.vk.UsersGet(api.Params{
		"user_ids": obj.Message.FromID,
	})
	if len(usrs) > 0 {
		fname := usrs[0].FirstName + " " + usrs[0].LastName
		author = &fname
	}
	msg := merger.CreateMessage{
		ReplyId: (*merger.ID)(replyedId),
		Date:    time.Unix(int64(obj.Message.Date), 0),
		Author:  author,
		Silent:  bool(obj.Message.IsSilent),
		Body: &merger.BodyText{
			Format: merger.Plain,
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
		b.Message(msg.FormatShort())
		b.RandomID(0)
		b.PeerID(c.peerID)

		_, err = c.vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
	}
}
