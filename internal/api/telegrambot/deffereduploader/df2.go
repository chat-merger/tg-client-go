package deffereduploader

import (
	"context"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
	"sync"
	"time"
)

var _ DeferredUploader = (*DeferredUploader2)(nil)

type DeferredUploader2 struct {
	releaseTimeout time.Duration
	mu             *sync.Mutex
	mm             kvstore.MessagesMap
}

type Runner struct {
	//current *merger.Message
	//prev    Kind
	prev MsgWithKind
	// todo + shared mutex
	queue chan MsgWithKind
	com   Comparator
}

func (r *Runner) ReadNext(ctx context.Context) {
	var mwk MsgWithKind
	select {
	case mwk = <-r.queue:
	case <-ctx.Done():
		return
	}

	r.prev = r.com.CompareAndMbSendAndReturnCurrent(r.prev, mwk)
}

func (r *Runner) Put() chan<- MsgWithKind {
	return r.queue
}

type MsgWithKind struct {
	kind Kind
	msg  merger.Message
}

// type StrategyBuckets map[Kind]StrategyBucket

//type UserBuckets map[int64]UserBucket
//
//type StrategyBucket struct {
//	msg       *merger.Message
//	releaseIn time.Time
//}
//
//type UserBucket struct {
//	lastStrategy Kind
//	releaseIn    time.Time
//
//	groupMedia *merger.Message
//	media      *merger.Message
//	texted     *merger.Message
//	forward    *merger.Message
//}

func (d *DeferredUploader2) Upload(msg gotgbot.Message) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	uid := msg.GetSender().Id()
	ds := defineStrategy(msg)
	err := d.sendToHandler(msg, ds)

	return nil
}
func (d *DeferredUploader2) sendToHandler(msg gotgbot.Message, ds Kind) error {
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

func (d *DeferredUploader2) handleGroupMedia(msg gotgbot.Message, prev *merger.Message) (*merger.Message, error) {
	if prev == nil {
		return d.msgToMerger(msg)
	}
}

func defineStrategy(msg gotgbot.Message) Kind {
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
