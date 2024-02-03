package deffereduploader

import (
	"context"
	"log"
	"time"
)

type IRunner interface {
	Run(
		ctx context.Context,
		queue IQueue,
		comp IComparator,
		sender ISender,
		releasePeriod time.Duration,
	)
}

type Runner struct {
	prev *MsgWithKind
}

func (r *Runner) Run(
	ctx context.Context,
	queue IQueue,
	comp IComparator,
	sender ISender,
	releasePeriod time.Duration,
) {
	// Возобновляемый таймер.
	// Здесь используется чтобы отпралять последнее сохраненное сообщение (r.prev)
	// если после Compare оно не было отправлено (т.е. сообщения были склеены).
	timer := time.NewTimer(releasePeriod)
	timer.Stop()
	defer timer.Stop()

	sendPrev := func() {
		if r.prev != nil {
			err := sender.Send(*r.prev)
			if err != nil {
				log.Printf("[ERROR] sendPrev to sender: %s", err)
			}
			r.prev = nil // Очистить, т.к. сообщение считается отправленным
		}
	}

	for {
		select {
		case mwk, ok := <-queue.Next():
			if !ok {
				// Закрыли канал извне. Код достижим.
				log.Printf("[ERROR] chan of Runner queue is closed or empty")
				return
			}
			prev, compRes := comp.Compare(r.prev, mwk)
			// предыдущее было заменено, а не склеено
			if compRes == PutNext {
				sendPrev()
			}
			r.prev = &prev
			// Если последующее сообщение не поступит до окончания таймера,
			// то r.prev будет отправлено.
			timer.Reset(releasePeriod)

		case _, ok := <-timer.C:
			if !ok {
				// Недостижимый код. Вроде.
				log.Printf("[ERROR] chan of Runner timer is closed or empty")
				continue
			}
			sendPrev()

		case <-ctx.Done():
			return
		}
	}
}
