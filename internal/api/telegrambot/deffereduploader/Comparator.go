package deffereduploader

type Comparator interface {
	CompareAndMbSendAndReturnCurrent(prev *MsgWithKind, next MsgWithKind) (current MsgWithKind)
}

type ComparatorImpl struct {
	sender ISender
}

func (c *ComparatorImpl) CompareAndMbSendAndReturnCurrent(prev *MsgWithKind, next MsgWithKind) (current MsgWithKind) {
	if prev == nil {
		return next
	}
	switch prev.kind {
	case GroupMedia:
		if next.kind == GroupMedia { // even when tgMsg.media_group_id not equals
			return mergeCommon(*prev, next)
		}
	case Texted:
		if next.kind == GroupMedia || next.kind == Forward {
			return mergeCommon(*prev, next)
		}
	case Forward:
		if next.kind == Forward {
			return mergeCommon(*prev, next)
		}
	default:
	}

	c.sender.Send(*prev)
	return next
}

func mergeCommon(a, b MsgWithKind) MsgWithKind {
	a.msg.Forwarded = append(a.msg.Forwarded, b.msg.Forwarded...)
	a.msg.Media = append(a.msg.Media, b.msg.Media...)
	return a
}
