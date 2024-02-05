package tghelper

import "github.com/PaulSonOfLars/gotgbot/v2"

type Kind uint8

const ( // db stored kind, then val is fixed
	Unknown    Kind = 0
	GroupMedia Kind = 1
	Media      Kind = 2
	Texted     Kind = 3
	Forward    Kind = 4
)

func (k Kind) IsMedia() bool {
	return k == Media || k == GroupMedia
}

func DefineKind(msg gotgbot.Message) Kind {
	switch {
	case IsForward(msg):
		return Forward
	case IsPartOfMediaGroup(msg):
		return GroupMedia
	case IsMedia(msg):
		return Media
	case HasText(msg):
		return Texted
	default:
		return Unknown
	}
}
