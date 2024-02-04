package deffereduploader

import (
	"context"
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
	mu             *sync.RWMutex
	mm             kvstore.MessagesMap
	con            IConvertor
	usersQueue     map[int64]*Queue
	runner         IRunner
	sender         ISender
}

func NewDeferredUploader(mm kvstore.MessagesMap, files blobstore.TempBlobStore, bot *gotgbot.Bot, conn merger.Conn) *DeferredUploader {
	s := NewSender(conn, mm)
	return &DeferredUploader{
		releaseTimeout: time.Millisecond * 70,
		mu:             new(sync.RWMutex),
		mm:             mm,
		con:            NewConvertor(mm, files, bot),
		usersQueue:     make(map[int64]*Queue),
		runner:         new(Runner),
		sender:         s,
	}
}

func (d *DeferredUploader) Upload(original gotgbot.Message) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	q, ok := d.usersQueue[original.GetSender().Id()]
	if !ok {
		d.mu.RUnlock()
		d.mu.Lock()
		q = NewQueue(make(chan gotgbot.Message, 50))
		d.usersQueue[original.GetSender().Id()] = q
		d.mu.Unlock()
		d.mu.RLock()
		go d.runner.Run(context.Background(), q, d.sender, d.con, d.releaseTimeout)
	}
	q.Ch() <- original
	return nil
}
