package deffereduploader

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
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
	case GroupMedia:
		if n == GroupMedia && tghelper.HasSameMediaGroup(prevMsg.original, nextMsg.original) {
			return PutMerged
		}
	case Texted:
		if n == GroupMedia || n == Forward {
			log.Println("[Compare] case Texted > PutMerged")
			return PutMerged
		}
	case Forward:
		if n == Forward {
			return PutMerged
		}
	default:
	}

	return PutNext
}

func DefineKind(msg gotgbot.Message) Kind {
	switch {
	case tghelper.IsForward(msg):
		return Forward
	case tghelper.IsPartOfMediaGroup(msg):
		return GroupMedia
	case tghelper.IsMedia(msg):
		return Media
	case tghelper.HasText(msg):
		return Texted
	default:
		return Unknown
	}
}
