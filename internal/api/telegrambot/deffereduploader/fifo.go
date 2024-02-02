package deffereduploader

//
//type Fifo[T any] interface {
//	Pub(e T)
//	Next() <-chan T
//}
//
//var _ Fifo[merger.CreateMessage] = (*FifoMwt)(nil)
//
//type FifoMwt struct {
//	ch  chan merger.CreateMessage
//}
//
//func NewFifoMwt(cap int) *FifoMwt {
//	return &FifoMwt{
//		ch: make(chan merger.CreateMessage,cap),
//	}
//}
//
//func (f *FifoMwt) Pub(msg merger.CreateMessage) {
//	f.ch <- msg
//}
//
//func (f *FifoMwt) Next() <-chan merger.CreateMessage {
//	return f.ch
//}

//func createFifoChan[T any]() chan T {
//	return make(chan T, )
//}
//type Fifo chan merger.CreateMessage
