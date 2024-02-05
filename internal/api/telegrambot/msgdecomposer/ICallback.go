package msgdecomposer

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/service/merger"
)

type ISender interface {
	SendTexted(msg merger.Message) (*gotgbot.Message, error)
	SendMediaGroup(msgs []merger.Message) ([]gotgbot.Message, error)
	SendPhoto(msg merger.Message) (*gotgbot.Message, error)
	SendSticker(msg merger.Message) (*gotgbot.Message, error)
	SendAudio(msg merger.Message) (*gotgbot.Message, error)
	SendVideo(msg merger.Message) (*gotgbot.Message, error)
	SendDocument(msg merger.Message) (*gotgbot.Message, error)
	SendForward(msg merger.Message) (*gotgbot.Message, error)
}
