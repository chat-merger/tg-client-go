package deffereduploader

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"log"
	"merger-adapter/internal/api/telegrambot/tghelper"
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
	var cm *merger.CreateMessage
	var err error
	if tghelper.IsForward(msg) { // forward:
		cm, err = buildMsgAsForward(msg, c.bot, c.files)
	} else { // not forward:
		cm, err = buildMsgAsOriginal(msg, c.bot, c.files, c.mm)
	}
	if err != nil {
		return merger.CreateMessage{}, fmt.Errorf("buildMsg: %s", err)
	}
	return *cm, nil
}

func buildMsgAsOriginal(msg gotgbot.Message, bot *gotgbot.Bot, files blobstore.TempBlobStore, mm kvstore.MessagesMap) (*merger.CreateMessage, error) {
	replyTo := replyMergerIdFromMsg(msg, mm)
	firstName := msg.GetSender().FirstName()
	createMsg := merger.CreateMessage{
		ReplyId:   (*merger.ID)(replyTo),
		Date:      time.Unix(msg.Date, 0),
		Username:  &firstName,
		Silent:    false, // where prop??
		Text:      &msg.Text,
		Media:     make([]merger.Media, 0),
		Forwarded: make([]merger.Forward, 0),
	}
	// add media
	media, err := getMediaOrNil(msg, bot, files)
	if err != nil {
		log.Printf("[ERROR] getMediaOrNil: %s", err)
	} else if media != nil {
		createMsg.Media = append(createMsg.Media, *media)
		createMsg.Text = &msg.Caption
	}
	return &createMsg, nil
}

func buildMsgAsForward(msg gotgbot.Message, bot *gotgbot.Bot, files blobstore.TempBlobStore) (*merger.CreateMessage, error) {
	firstName := msg.GetSender().FirstName()
	createMsg := merger.CreateMessage{
		ReplyId:   nil,
		Date:      time.Unix(msg.Date, 0),
		Username:  &firstName,
		Silent:    false, // where prop??
		Text:      nil,
		Media:     make([]merger.Media, 0),
		Forwarded: make([]merger.Forward, 0),
	}

	fwd := merger.Forward{
		Id:       nil,
		Date:     time.Unix(msg.ForwardDate, 0),
		Username: &msg.ForwardFrom.FirstName,
		Text:     &msg.Text,
		Media:    make([]merger.Media, 0),
	}
	// add media
	media, err := getMediaOrNil(msg, bot, files)
	if err != nil {
		log.Printf("[ERROR] getMediaOrNil: %s", err)
	} else if media != nil {
		// add fwd to msg
		fwd.Media = append(fwd.Media, *media)
		fwd.Text = &msg.Caption
	}
	createMsg.Forwarded = append(createMsg.Forwarded, fwd)
	return &createMsg, nil
}

func getMediaOrNil(original gotgbot.Message, bot *gotgbot.Bot, files blobstore.TempBlobStore) (*merger.Media, error) {
	var err error
	var media *merger.Media
	if len(original.Photo) != 0 {
		media, err = downloadMedia(original.Photo[len(original.Photo)-1].FileId, original.HasMediaSpoiler, merger.Photo, bot, files)
	} else if original.Video != nil {
		media, err = downloadMedia(original.Video.FileId, original.HasMediaSpoiler, merger.Video, bot, files)
	} else if original.Audio != nil {
		media, err = downloadMedia(original.Audio.FileId, original.HasMediaSpoiler, merger.Audio, bot, files)
	} else if original.Document != nil {
		media, err = downloadMedia(original.Document.FileId, original.HasMediaSpoiler, merger.File, bot, files)
	} else if original.Sticker != nil {
		media, err = downloadMedia(original.Sticker.FileId, original.HasMediaSpoiler, merger.Sticker, bot, files)
	}
	if err != nil {
		return nil, fmt.Errorf("downloadMedia: %s", err)
	}
	return media, nil
}

func toContID(id int64) kvstore.ContMsgID {
	return kvstore.ContMsgID(strconv.FormatInt(id, 10))
}

func replyMergerIdFromMsg(msg gotgbot.Message, mm kvstore.MessagesMap) *string {
	if msg.ReplyToMessage != nil {
		id, exists, err := mm.ByContID(tghelper.KvStoreScope, toContID(msg.ReplyToMessage.MessageId))
		if err != nil {
			log.Printf("[ERROR] msg from message map: %s", err)
		}
		if exists {
			return (*string)(id)
		}
	}
	return nil
}

func downloadMedia(fileId string, hasMediaSpoiler bool, mtype merger.MediaType, bot *gotgbot.Bot, files blobstore.TempBlobStore) (*merger.Media, error) {
	file, err := bot.GetFile(fileId, nil)
	if err != nil {
		return nil, fmt.Errorf(" get file from blobstore: %s", err)
	}
	get, err := http.Get(file.URL(bot, nil))
	if err != nil {
		return nil, fmt.Errorf("http get: %s", err)
	}

	uri, err := files.Save(get.Body)
	if err != nil {
		return nil, fmt.Errorf("uri file to blobstore: %s", err)
	}
	err = get.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("close http body: %s", err)
	}
	return &merger.Media{
		Kind:    mtype,
		Spoiler: hasMediaSpoiler,
		Url:     *uri,
	}, nil
}
