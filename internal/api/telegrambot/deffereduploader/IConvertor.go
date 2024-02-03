package deffereduploader

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"log"
	"merger-adapter/internal/service/blobstore"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
	"net/http"
	"strconv"
	"time"
)

type IConvertor interface {
	Convert(msg gotgbot.Message) (merger.CreateMessage, error)
}
type Convertor struct {
	mm    kvstore.MessagesMap
	files blobstore.TempBlobStore
	bot   *gotgbot.Bot
}

func NewConvertor(mm kvstore.MessagesMap, files blobstore.TempBlobStore, bot *gotgbot.Bot) *Convertor {
	return &Convertor{mm: mm, files: files, bot: bot}
}

func (c *Convertor) Convert(msg gotgbot.Message) (merger.CreateMessage, error) {
	replyTo := replyMergerIdFromMsg(msg, c.mm)
	author := msg.GetSender().Username()
	medias := make([]merger.Media, 0, len(msg.Photo))
	for _, ps := range msg.Photo {
		file, err := c.bot.GetFile(ps.FileId, nil)
		if err != nil {
			log.Printf("[ERROR] get file from blobstore: %s", err)
			continue
		}

		get, err := http.Get(file.URL(c.bot, nil))
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
			return merger.CreateMessage{}, err
		}
		medias = append(medias, merger.Media{
			Kind:    merger.Photo,
			Spoiler: msg.HasMediaSpoiler,
			Url:     *uri,
		})
	}

	createMsg := merger.CreateMessage{
		ReplyId:   (*merger.ID)(replyTo),
		Date:      time.Unix(msg.Date, 0),
		Username:  &author,
		Silent:    false, // where prop??
		Text:      &msg.Text,
		Media:     medias,
		Forwarded: nil,
	}
	return createMsg, nil
}

const mmScope = kvstore.Scope("TgBotScope")

func toContID(id int64) kvstore.ContMsgID {
	return kvstore.ContMsgID(strconv.FormatInt(id, 10))
}

func replyMergerIdFromMsg(msg gotgbot.Message, mm kvstore.MessagesMap) *string {
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
