package deffereduploader

import "context"

type IRunner interface {
	Run(ctx context.Context, queue IQueue, comp Comparator) error
}

type Runner struct {
	prev *MsgWithKind
}

func (r *Runner) Run(ctx context.Context, queue IQueue, comp Comparator) {
	for {
		var mwk MsgWithKind
		select {
		case mwk = <-queue.Next():
		case <-ctx.Done():
			return
		}
		prev := comp.CompareAndMbSendAndReturnCurrent(r.prev, mwk)
		r.prev = &prev
	}
}
