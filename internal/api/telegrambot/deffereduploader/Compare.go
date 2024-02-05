package deffereduploader

import (
	"log"
	"merger-adapter/internal/api/telegrambot/tghelper"
)

type CompareResult uint8

const (
	PutNext CompareResult = iota
	PutMerged
)

func Compare(prevMsg *MsgWithKind, nextMsg MsgWithKind) (result CompareResult) {
	if prevMsg == nil {
		log.Println("[Comparator.Compare...] prev == nil")
		return PutNext
	}
	p := prevMsg.kind
	n := nextMsg.kind
	switch p {
	case tghelper.GroupMedia:
		if n == tghelper.GroupMedia && tghelper.HasSameMediaGroup(prevMsg.original, nextMsg.original) {
			return PutMerged
		}
	case tghelper.Texted:
		if n == tghelper.GroupMedia || n == tghelper.Forward {
			log.Println("[Compare] case Texted > PutMerged")
			return PutMerged
		}
	case tghelper.Forward:
		if n == tghelper.Forward {
			return PutMerged
		}
	default:
	}

	return PutNext
}
