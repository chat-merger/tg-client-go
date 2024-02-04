package deffereduploader

import (
	"context"
	"log"
	"merger-adapter/internal/service/merger"
	"time"
)

type IRunner interface {
	Run(
		ctx context.Context,
		queue IQueue,
		sender ISender,
		conv IConvertor,
		releasePeriod time.Duration,
	)
}

//type MsgWithStrat struct {
//	original gotgbot.Message
//	strat    CompareResultStrat
//}

type Runner struct {
	prev *MsgWithKind
	//// Deprecated. use prev
	//prevf *MsgWithStrat
}

func (r *Runner) Run(
	ctx context.Context,
	queue IQueue,
	sender ISender,
	conv IConvertor,
	releasePeriod time.Duration,
) {
	// Возобновляемый таймер.
	// Здесь используется чтобы отпралять последнее сохраненное сообщение (r.prev)
	// если после Compare оно не было отправлено (т.е. сообщения были склеены).
	timer := time.NewTimer(releasePeriod)
	timer.Stop()
	defer timer.Stop()

	sendAndCleanPrev := func(msg MsgWithKind) {
		err := sender.Send(msg)
		if err != nil {
			log.Printf("[ERROR] sender.Send: %s", err)
		}
		r.prev = nil // Очистить, т.к. сообщение считается отправленным
	}

	for {
		select {
		case orig, ok := <-queue.Next():
			if !ok {
				// Закрыли канал извне. Код достижим.
				log.Printf("[ERROR] chan of Runner queue is closed or empty")
				return
			}
			nextMsg, err := conv.Convert(orig)
			if err != nil {
				log.Printf("conv.Convert: %s", err)
				continue
			}
			next := &MsgWithKind{
				kind:     defineKind(orig),
				original: orig,
				msg:      nextMsg,
			}
			compRes := Compare(r.prev, *next)
			// предыдущее было заменено, а не склеено
			switch compRes {
			case PutNext:
				if r.prev != nil {
					sendAndCleanPrev(*r.prev)
				}
				r.prev = next // put next to prev
			case PutMerged:
				r.prev.msg = merge(r.prev.msg, next.msg)
			}
			// Если последующее сообщение не поступит до окончания таймера,
			// то r.prev будет отправлено.
			timer.Reset(releasePeriod)

		case _, ok := <-timer.C:
			if !ok {
				// Недостижимый код. Вроде.
				log.Printf("[ERROR] chan of Runner timer is closed or empty")
				continue
			}
			if r.prev != nil {
				sendAndCleanPrev(*r.prev)
			}

		case <-ctx.Done():
			return
		}
	}
}

func merge(a, b merger.CreateMessage) merger.CreateMessage {
	a.Forwarded = append(a.Forwarded, b.Forwarded...)
	a.Media = append(a.Media, b.Media...)
	return a
}
