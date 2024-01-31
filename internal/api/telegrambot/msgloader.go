package telegrambot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/service/merger"
)

type DefferedMessage interface {
	PutForward(forward *gotgbot.Message)
	PutForward1(forward)
}

type