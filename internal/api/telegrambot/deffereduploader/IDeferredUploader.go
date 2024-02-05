package deffereduploader

import (
	"context"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/api/telegrambot/tghelper"
	mrepo "merger-adapter/internal/repository/messages_repository"
	"merger-adapter/internal/service/blobstore"
	"merger-adapter/internal/service/merger"
	"sync"
	"time"
)

type IDeferredUploader interface {
	Upload(msg gotgbot.Message) error
}

type MsgWithKind struct {
	kind     tghelper.Kind
	original gotgbot.Message
	msg      merger.CreateMessage
}

// IMPLEMENTATION:

var _ IDeferredUploader = (*DeferredUploader)(nil)

type DeferredUploader struct {
	mu         *sync.RWMutex
	con        IConvertor
	usersQueue map[int64]*Queue
	runner     IRunner
	sender     ISender
}

func NewDeferredUploader(repo mrepo.MessagesRepository, files blobstore.TempBlobStore, bot *gotgbot.Bot, conn merger.Conn) *DeferredUploader {
	s := NewSender(conn, repo)
	return &DeferredUploader{
		mu:         new(sync.RWMutex),
		con:        NewConvertor(repo, files, bot),
		usersQueue: make(map[int64]*Queue),
		runner:     new(Runner),
		sender:     s,
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
		go d.runner.Run(
			context.Background(),
			q,
			d.sender,
			d.con,
			func(msg *MsgWithKind) time.Duration {
				// форварды приходят с большой задержкой,
				// поэтому текст надо долго держать и ждать возможных форвардов
				if msg.kind == tghelper.Texted {
					return 400 * time.Millisecond
				} else {
					return 70 * time.Millisecond
				}
			},
		)
	}
	q.Ch() <- original
	return nil
}
