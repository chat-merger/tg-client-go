package deffereduploader

import "log"

type IComparator interface {
	Compare(prev *MsgWithKind, next MsgWithKind) (current MsgWithKind, result CompareResult)
}

type CompareResult uint8

const (
	PutNext CompareResult = iota
	Merged
)

type Comparator struct{}

func NewComparatorImpl() *Comparator {
	return &Comparator{}
}

func (c *Comparator) Compare(prev *MsgWithKind, next MsgWithKind) (current MsgWithKind, result CompareResult) {
	if prev == nil {
		log.Println("[Comparator.Compare...] prev == nil")
		return next, PutNext
	}
	switch prev.kind {
	case GroupMedia:
		//log.Println("[Comparator.Compare...] case GroupMedia")
		if next.kind == GroupMedia { // even when tgMsg.media_group_id not equals
			return mergeCommon(*prev, next), Merged
		}
	case Texted:
		//log.Println("[Comparator.Compare...] case Texted")
		if next.kind == GroupMedia || next.kind == Forward {
			return mergeCommon(*prev, next), Merged
		}
	case Forward:
		//log.Println("[Comparator.Compare...] case Forward")
		if next.kind == Forward {
			return mergeCommon(*prev, next), Merged
		}
	default:
		//log.Println("[Comparator.Compare...] default")
	}

	return next, PutNext
}

func mergeCommon(a, b MsgWithKind) MsgWithKind {
	a.msg.Forwarded = append(a.msg.Forwarded, b.msg.Forwarded...)
	a.msg.Media = append(a.msg.Media, b.msg.Media...)
	return a
}
