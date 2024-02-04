package deffereduploader

import "github.com/PaulSonOfLars/gotgbot/v2"

type IQueue interface {
	Next() <-chan gotgbot.Message
}

type Queue struct {
	ch chan gotgbot.Message
}

func NewQueue(ch chan gotgbot.Message) *Queue {
	return &Queue{
		ch: ch,
	}
}

func (q *Queue) Next() <-chan gotgbot.Message {
	return q.ch
}

func (q *Queue) Ch() chan<- gotgbot.Message {
	return q.ch
}

func (q *Queue) Close() {
	for {
		select {
		case _, ok := <-q.ch:
			if !ok {
				return
			}
		default:
			close(q.ch)
			return
		}
	}
}
