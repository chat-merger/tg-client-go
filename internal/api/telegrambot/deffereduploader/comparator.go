package deffereduploader

type Comparator interface {
	CompareAndMbSendAndReturnCurrent(prev, next MsgWithKind) (current MsgWithKind)
}

type ComparatorImpl struct {
}

func (c *ComparatorImpl) CompareAndMbSendAndReturnCurrent(prev, next MsgWithKind) (current MsgWithKind) {
	// prev?
}
