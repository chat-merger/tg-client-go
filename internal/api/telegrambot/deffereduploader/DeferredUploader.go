package deffereduploader

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/service/blobstore"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
	"sync"
	"time"
)

type IDeferredUploader interface {
	Upload(msg gotgbot.Message) error
}

type Kind uint8

const (
	Unknown Kind = iota
	GroupMedia
	Media
	Texted
	Forward
)

type MsgWithKind struct {
	kind     Kind
	original gotgbot.Message
	msg      merger.CreateMessage
}

// IMPLEMENTATION:

var _ IDeferredUploader = (*DeferredUploader)(nil)

type DeferredUploader struct {
	releaseTimeout time.Duration
	mu             *sync.Mutex
	mm             kvstore.MessagesMap
	con            IConvertor
	queue          *Queue
}

func NewDeferredUploader(mm kvstore.MessagesMap, files blobstore.TempBlobStore, bot *gotgbot.Bot) *DeferredUploader {
	return &DeferredUploader{
		releaseTimeout: time.Millisecond * 100,
		mu:             new(sync.Mutex),
		mm:             mm,
		con:            NewConvertor(mm, files, bot),
		queue:          NewQueue(make(chan MsgWithKind, 500)),
	}
}

func (d *DeferredUploader) Upload(original gotgbot.Message) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	msg, err := d.con.Convert(original)
	if err != nil {
		return fmt.Errorf("convert original gotgbot msg to merger: %s", err)
	}
	mwk := MsgWithKind{
		kind:     defineKind(original),
		original: original,
		msg:      msg,
	}
	d.queue.Ch() <- mwk
	return nil
}

func defineKind(msg gotgbot.Message) Kind {
	switch {
	case msg.ForwardDate != 0:
		return Forward
	case msg.MediaGroupId != "":
		return GroupMedia
	case isMedia(msg):
		return Media
	case msg.Text != "":
		return Texted
	default:
		return Unknown
	}
}

func isMedia(msg gotgbot.Message) bool {
	return isPhoto(msg) || isVideo(msg) || isAudio(msg) || isDocument(msg) || isSticker(msg)
}

func isPhoto(msg gotgbot.Message) bool {
	return len(msg.Photo) > 0
}

func isVideo(msg gotgbot.Message) bool {
	return msg.Video != nil
}

func isAudio(msg gotgbot.Message) bool {
	return msg.Audio != nil
}

func isDocument(msg gotgbot.Message) bool {
	return msg.Document != nil
}

func isSticker(msg gotgbot.Message) bool {
	return msg.Sticker != nil
}
