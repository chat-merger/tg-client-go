package deffereduploader

type IQueue interface {
	Next() <-chan MsgWithKind
}

type Queue struct {
	ch chan MsgWithKind
}

func NewQueue(ch chan MsgWithKind) *Queue {
	return &Queue{ch: ch}
}

func (q *Queue) Next() <-chan MsgWithKind {
	return q.ch
}

func (q *Queue) Ch() chan<- MsgWithKind {
	return q.ch
}
