package deffereduploader

import (
	"context"
	"log"
	"merger-adapter/internal/api/telegrambot/tghelper"
	"merger-adapter/internal/service/merger"
	"time"
)

type IRunner interface {
	Run(
		ctx context.Context,
		queue IQueue,
		sender ISender,
		conv IConvertor,
		releasePeriodFun func(*MsgWithKind) time.Duration,
	)
}

type Runner struct {
	prev *MsgWithKind
}

func (r *Runner) Run(
	ctx context.Context,
	queue IQueue,
	sender ISender,
	conv IConvertor,
	releasePeriodFun func(*MsgWithKind) time.Duration,
) {
	// Возобновляемый таймер.
	// Здесь используется чтобы отпралять последнее сохраненное сообщение (r.prev)
	// если после Compare оно не было отправлено (т.е. сообщения были склеены).
	timer := time.NewTimer(0)
	stopTimer(timer)
	defer stopTimer(timer)

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
				log.Printf("[ERROR] conv.Convert: %s", err)
				continue
			}
			next := &MsgWithKind{
				kind:     tghelper.DefineKind(orig),
				original: orig,
				msg:      nextMsg,
			}
			compRes := Compare(r.prev, *next)

			switch compRes {
			case PutNext: // предыдущее было заменено, а не склеено
				if r.prev != nil {
					sendAndCleanPrev(*r.prev)
				}
				r.prev = next // put next to prev

			case PutMerged:
				r.prev.msg = merge(r.prev.msg, next.msg)
				r.prev.kind = next.kind
				// but except:
				// `r.prev.original = next.original`
				// then first receive relation merger id
			}
			// Если последующее сообщение не поступит до окончания таймера,
			// то r.prev будет отправлено.
			stopTimer(timer).Reset(releasePeriodFun(r.prev))

		case <-timer.C:
			if r.prev != nil {
				log.Printf("[<-timer.C] send %v", r.prev.kind)
				sendAndCleanPrev(*r.prev)
			}

		case <-ctx.Done():
			return
		}
	}
}

func stopTimer(timer *time.Timer) *time.Timer {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	return timer
}

func merge(a, b merger.CreateMessage) merger.CreateMessage {
	if (a.Text == nil || *a.Text == "") && (b.Text != nil && *b.Text != "") {
		a.Text = b.Text
	}
	a.Forwarded = append(a.Forwarded, b.Forwarded...)
	a.Media = append(a.Media, b.Media...)
	return a
}
