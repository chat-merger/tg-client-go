package deffereduploader

import (
	"context"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/api/telegrambot/tghelper"
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
	usersQueue     map[int64]*Queue
	runner         IRunner
	comp           IComparator
	sender         ISender
}

func NewDeferredUploader(mm kvstore.MessagesMap, files blobstore.TempBlobStore, bot *gotgbot.Bot, conn merger.Conn) *DeferredUploader {
	s := NewSender(conn, mm)
	return &DeferredUploader{
		releaseTimeout: time.Millisecond * 600,
		mu:             new(sync.Mutex),
		mm:             mm,
		con:            NewConvertor(mm, files, bot),
		queue:          InitQueue(make(chan MsgWithKind, 500)),
		usersQueue:     make(map[int64]*Queue),
		runner:         new(Runner),
		comp:           NewComparatorImpl(),
		sender:         s,
	}
}

func (d *DeferredUploader) Upload(original gotgbot.Message) error {

	msg, err := d.con.Convert(original)
	if err != nil {
		return fmt.Errorf("convert original gotgbot msg to merger: %s", err)
	}
	mwk := MsgWithKind{
		kind:     defineKind(original),
		original: original,
		msg:      msg,
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	q, ok := d.usersQueue[original.GetSender().Id()]
	if !ok {
		q = InitQueue(make(chan MsgWithKind, 50))
		d.usersQueue[original.GetSender().Id()] = q
		go d.runner.Run(context.Background(), q, d.comp, d.sender, d.releaseTimeout)
	}
	q.Ch() <- mwk
	return nil
}

func defineKind(msg gotgbot.Message) Kind {
	switch {
	case tghelper.IsForward(msg):
		return Forward
	case tghelper.IsPartOfMediaGroup(msg):
		return GroupMedia
	case tghelper.IsMedia(msg):
		return Media
	case tghelper.HasText(msg):
		return Texted
	default:
		return Unknown
	}
}
