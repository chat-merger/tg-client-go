package telegrambot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"merger-adapter/internal/api/telegrambot/deffereduploader"
	"merger-adapter/internal/service/blobstore"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
)

type Client struct {
	bot         *gotgbot.Bot
	updater     *ext.Updater
	conn        merger.Conn
	chatID      int64
	messagesMap kvstore.MessagesMap
	du          deffereduploader.IDeferredUploader
}

type Deps struct {
	Token       string
	ApiKey      string
	ChatID      int64
	Server      merger.MergerServer
	MessagesMap kvstore.MessagesMap
	Files       blobstore.TempBlobStore
}
