package telegrambot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/api/telegrambot/tghelper"
)

func (c *Client) filter(msg *gotgbot.Message) bool {
	return tghelper.InSpecificChat(*msg, c.chatID) && (tghelper.HasText(*msg) || tghelper.IsMedia(*msg))
}
