package telegrambot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"merger-adapter/internal/service/blobstore"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
)

type Client struct {
	bot         *gotgbot.Bot
	dispatcher  *ext.Dispatcher
	updater     *ext.Updater
	conn        merger.Conn
	chatID      int64
	messagesMap kvstore.MessagesMap
	files       blobstore.TempBlobStore
}

type Deps struct {
	Token       string
	ApiKey      string
	ChatID      int64
	Server      merger.MergerServer
	MessagesMap kvstore.MessagesMap
	Files       blobstore.TempBlobStore
}
