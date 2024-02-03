package telegrambot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func (c *Client) filter(msg *gotgbot.Message) bool {
	return inSpecificChat(msg, c.chatID) && (isWithText(msg) || isMedia(msg))
}

func inSpecificChat(msg *gotgbot.Message, spec int64) bool {
	return msg.Chat.Id == spec
}

func isWithText(msg *gotgbot.Message) bool {
	return msg.Text != ""
}

func isMedia(msg *gotgbot.Message) bool {
	return isPhoto(msg) || isVideo(msg) || isAudio(msg) || isDocument(msg) || isSticker(msg)
}

func isPhoto(msg *gotgbot.Message) bool {
	return len(msg.Photo) > 0
}

func isVideo(msg *gotgbot.Message) bool {
	return msg.Video != nil
}

func isAudio(msg *gotgbot.Message) bool {
	return msg.Audio != nil
}

func isDocument(msg *gotgbot.Message) bool {
	return msg.Document != nil
}

func isSticker(msg *gotgbot.Message) bool {
	return msg.Sticker != nil
}
