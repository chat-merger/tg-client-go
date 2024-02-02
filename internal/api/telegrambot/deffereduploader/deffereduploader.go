package deffereduploader

import "github.com/PaulSonOfLars/gotgbot/v2"

type DeferredUploader interface {
	Upload(msg gotgbot.Message) error
}

type Kind uint8

const (
	Unknown Kind = iota
	GroupMedia
	Media
	Texted
	Forward
)
