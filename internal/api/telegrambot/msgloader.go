package telegrambot

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/service/merger"
	"sync"
	"time"
)

//type DeferredMessage interface {
//	PutForward(forward *gotgbot.Message)
//}
//
//type DeferredSender interface {
//	addGroupMedia(msg *gotgbot.Message)
//	addMedia(msg *gotgbot.Message)
//	addTexted(msg *gotgbot.Message)
//	addForward(msg *gotgbot.Message)
//}

type DeferredUploader interface {
	Upload(msg gotgbot.Message, ds DeferringStrategy) error
}

type DeferringStrategy uint8

const (
	_ DeferringStrategy = iota
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
}

type StrategyBuckets map[DeferringStrategy]StrategyBucket
type UserBuckets map[int64]StrategyBuckets

func (ub UserBuckets) addToStrategy(msg gotgbot.Message, ds DeferringStrategy) {
	ub.
}

type StrategyBucket struct {
	msg       *merger.Message
	releaseIn time.Time
}

func (d *DeferredUploaderImpl) Upload(msg gotgbot.Message, ds DeferringStrategy) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	uid := msg.GetSender().Id()
	x, ok := d.userBuckets[uid]
	if !ok {
		x = make(StrategyBuckets)
		d.userBuckets[uid] = x
	}
	bst, ok := x[ds]
	if !ok {
		bst = StrategyBucket{
			msg:       nil,
			releaseIn: time.Now().Add(d.releaseTimeout),
		}
		x[ds] = bst
	}
	err := bst.Add(msg, ds)
	if err != nil {
		return fmt.Errorf("add to strategy bucket: %s", err)
	}

	return nil
}

func (sb *StrategyBucket) Add(msg gotgbot.Message, ds DeferringStrategy) error {
	switch ds {
	case GroupMedia:
	case Media:
	case Texted:
		//
	case Forward:
	}
}
