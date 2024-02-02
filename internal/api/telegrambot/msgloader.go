package telegrambot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
	"sync"
	"time"
)

type DeferredUploader interface {
	Upload(msg gotgbot.Message, ds DeferringStrategy) error
}

type DeferringStrategy uint8

const (
	Unknown DeferringStrategy = iota
	GroupMedia
	Media
	Texted
	Forward
)

type DeferredUploaderImpl struct {
	checkPeriod    time.Duration
	releaseTimeout time.Duration
	userBuckets    UserBuckets
	mu             *sync.Mutex

	mm kvstore.MessagesMap
}

// type StrategyBuckets map[DeferringStrategy]StrategyBucket

type UserBuckets map[int64]UserBucket

type StrategyBucket struct {
	msg       *merger.Message
	releaseIn time.Time
}

type UserBucket struct {
	lastStrategy DeferringStrategy
	releaseIn    time.Time

	groupMedia *merger.Message
	media      *merger.Message
	texted     *merger.Message
	forward    *merger.Message
}

func (d *DeferredUploaderImpl) Upload(msg gotgbot.Message) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	uid := msg.GetSender().Id()
	ds := defineStrategy(msg)
	err := d.sendToHandler(msg, ds)

	return nil
}
func (d *DeferredUploaderImpl) sendToHandler(msg gotgbot.Message, ds DeferringStrategy) error {
	ubkt := d.userBuckets[msg.GetSender().Id()]
	var newMsg *merger.Message
	var err error
	switch ds {
	case GroupMedia:
		newMsg, err = d.handleGroupMedia(msg, ubkt.groupMedia)
	case Media:
	case Texted:
	case Forward:
	case Unknown:
	}
	if err != nil {

	}
	ubkt.releaseIn = time.Now().Add(d.releaseTimeout)
	return nil
}

func (d *DeferredUploaderImpl) handleGroupMedia(msg gotgbot.Message, prev *merger.Message) (*merger.Message, error) {
	if prev == nil {
		return d.msgToMerger(msg)
	}
}

func defineStrategy(msg gotgbot.Message) DeferringStrategy {
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

func isSticker(msg *gotgbot.Message) bool {
	return msg.Sticker != nil
}
