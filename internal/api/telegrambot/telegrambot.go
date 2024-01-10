package telegrambot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"merger-adapter/internal/service/merger"
)

type Client struct {
	bot        *gotgbot.Bot
	dispatcher *ext.Dispatcher
	updater    *ext.Updater
	conn       merger.Conn
	chatID     int64
}

type Config struct {
	Token  string
	ApiKey string
	ChatID int64
	Server merger.MergerServer
}
