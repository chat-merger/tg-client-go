package telegrambot

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
	"net/http"
	"strconv"
	"time"
)

func replyMergerIdFromMsg(msg *gotgbot.Message, mm kvstore.MessagesMap) *string {
	if msg.ReplyToMessage != nil {
		id, exists, err := mm.ByContID(mmScope, toContID(msg.ReplyToMessage.MessageId))
		if err != nil {
			log.Printf("[ERROR] msg from message map: %s", err)
		}
		if exists {
			return (*string)(id)
		}
	}
	return nil
}

func (c *Client) onTelegramMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	replyTo := replyMergerIdFromMsg(ctx.Message, c.messagesMap)
	author := ctx.EffectiveUser.Username
	ctx.EffectiveMessage.
		medias := make([]merger.Media, 0, len(ctx.Message.Photo))
	for _, ps := range ctx.Message.Photo {
		file, err := c.bot.GetFile(ps.FileId, nil)
		if err != nil {
			log.Printf("[ERROR] get file from blobstore: %s", err)
			continue
		}

		get, err := http.Get(file.URL(b, nil))
		if err != nil {
			log.Printf("[ERROR] http get: %s", err)
			continue
		}

		uri, err := c.files.Save(get.Body)
		if err != nil {
			log.Printf("[ERROR] uri file to blobstore: %s", err)
			continue
		}
		err = get.Body.Close()
		if err != nil {
			log.Printf("[ERROR] close http body: %s", err)
			return err
		}
		medias = append(medias, merger.Media{
			Kind:    merger.Photo,
			Spoiler: ctx.Message.HasMediaSpoiler,
			Url:     *uri,
		})
	}

	msg := merger.CreateMessage{
		ReplyId:   (*merger.ID)(replyTo),
		Date:      time.Unix(ctx.Message.Date, 0),
		Username:  &author,
		Silent:    false, // where prop??
		Text:      &ctx.Message.Text,
		Media:     medias,
		Forwarded: nil,
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
	err = c.messagesMap.Save(mmScope, msg.Id, toContID(tgMsg.MessageId))
	if err != nil {
		log.Printf("[ERROR] msg id to MessageMap: %s", err)
		return
	}
}

func replyTgIdFromMsg(msg *merger.Message, mm kvstore.MessagesMap) int64 {
	if msg.ReplyId != nil {
		id, exists, err := mm.ByMergedID(mmScope, *msg.ReplyId)
		if err != nil {
			log.Printf("[ERROR] msg from message map: %s", err)
		}
		if exists {
			return toInt64(*id)
		}
	}
	return 0
}
