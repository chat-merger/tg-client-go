package tghelper

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/service/kvstore"
)

const KvStoreScope = kvstore.Scope("TgBotScope")

func HasSameMediaGroup(a, b gotgbot.Message) bool {
	return a.MediaGroupId == b.MediaGroupId
}

func IsForward(msg gotgbot.Message) bool {
	return msg.ForwardOrigin != nil
}

func IsMedia(msg gotgbot.Message) bool {
	return IsPhoto(msg) || IsVideo(msg) || IsAudio(msg) || IsDocument(msg) || IsSticker(msg)
}

func IsPartOfMediaGroup(msg gotgbot.Message) bool {
	return msg.MediaGroupId != ""
}

func IsPhoto(msg gotgbot.Message) bool {
	return len(msg.Photo) > 0
}

func IsVideo(msg gotgbot.Message) bool {
	return msg.Video != nil
}

func IsAudio(msg gotgbot.Message) bool {
	return msg.Audio != nil
}

func IsDocument(msg gotgbot.Message) bool {
	return msg.Document != nil
}

func IsSticker(msg gotgbot.Message) bool {
	return msg.Sticker != nil
}

func InSpecificChat(msg gotgbot.Message, spec int64) bool {
	return msg.Chat.Id == spec
}

func HasText(msg gotgbot.Message) bool {
	return msg.Text != ""
}
