package telegrambot

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"merger-adapter/internal/api/telegrambot/tghelper"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
	"strconv"
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
		go c.onMergerMessage(msg)
	}
}

func (c *Client) onMergerMessage(msg *merger.Message) {
	// reply
	replyTo := replyTgIdFromMsg(msg, c.messagesMap)

	tgMsg, err := c.bot.SendMessage(
		c.chatID,
		msg.FormatShort(),
		&gotgbot.SendMessageOpts{ReplyToMessageId: replyTo},
	)
	if err != nil {
		log.Printf("[ERROR] send message to tg: %s", err)
		return
	}
	err = c.messagesMap.Save(tghelper.KvStoreScope, msg.Id, toContID(tgMsg.MessageId))
	if err != nil {
		log.Printf("[ERROR] msg id to MessageMap: %s", err)
		return
	}
}

func toInt64(id kvstore.ContMsgID) int64 {
	vkMsgId, err := strconv.Atoi(string(id))
	if err != nil {
		log.Fatalf("[ERROR] convert kvstore.ContMsgID to int: %s", err)
	}
	return int64(vkMsgId)
}

func replyTgIdFromMsg(msg *merger.Message, mm kvstore.MessagesMap) int64 {
	if msg.ReplyId != nil {
		id, exists, err := mm.ByMergedID(tghelper.KvStoreScope, *msg.ReplyId)
		if err != nil {
			log.Printf("[ERROR] msg from message map: %s", err)
		}
		if exists {
			return toInt64(*id)
		}
	}
	return 0
}

func toContID(id int64) kvstore.ContMsgID {
	return kvstore.ContMsgID(strconv.FormatInt(id, 10))
}
