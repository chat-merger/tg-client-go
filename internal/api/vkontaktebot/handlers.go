package vkontaktebot

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"log"
	"math/rand"
	"merger-adapter/internal/component/debug"
	"merger-adapter/internal/service/kvstore"
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
		Uername: author,
		Silent:  bool(obj.Message.IsSilent),
		Body: &merger.BodyText{
			Format: merger.Plain,
			Value:  obj.Message.Text,
		},
	}

	mMsg, err := c.conn.Send(msg)
	if err != nil {
		log.Fatalf("send message to Server: %s", err)
	}
	err = c.messagesMap.Save(mmScope, mMsg.Id, toContID(obj.Message.ID))
	if err != nil {
		log.Printf("[ERROR] vkbot onMessage: save msg id to MessageMap: %s", err)
	}
}

const mmScope = kvstore.Scope("VkBotScope")

func toContID(id int) kvstore.ContMsgID {
	return kvstore.ContMsgID(strconv.Itoa(id))
}

func toInt(id kvstore.ContMsgID) int {
	vkMsgId, err := strconv.Atoi(string(id))
	if err != nil {
		log.Fatalf("[ERROR] convert kvstore.ContMsgID to int: %s", err)
	}
	return vkMsgId
}

func (c *Client) listenServerMessages() error {
	for {
		msg, err := c.conn.Update()
		if err != nil {
			return fmt.Errorf("receive update: %s", err)
		}
		vkMsgId := rand.New(rand.NewSource(time.Now().Unix())).Int()
		b := params.NewMessagesSendBuilder()
		b.Message(msg.FormatShort())
		b.RandomID(vkMsgId)
		b.PeerID(c.peerID)
		// reply
		if msg.ReplyId != nil {
			id, exists, err := c.messagesMap.ByMergedID(mmScope, *msg.ReplyId)
			if err != nil {
				log.Printf("[ERROR] msg from message map: %s", err)
			}
			if exists {
				b.ReplyTo(toInt(*id))
			}
		}

		_, err = c.vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		debug.Print(vkMsgId)
		err = c.messagesMap.Save(mmScope, msg.Id, toContID(vkMsgId))
		if err != nil {
			return fmt.Errorf("save msg id to MessageMap: %s", err)
		}
	}
}
