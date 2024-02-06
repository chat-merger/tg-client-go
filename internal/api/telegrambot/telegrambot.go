package telegrambot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"merger-adapter/internal/api/telegrambot/deffereduploader"
	mrepo "merger-adapter/internal/repository/messages_repository"
	"merger-adapter/internal/service/blobstore"
	"merger-adapter/internal/service/merger"
)

type Client struct {
	bot     *gotgbot.Bot
	updater *ext.Updater
	conn    merger.Conn
	chatID  int64
	repo    mrepo.MessagesRepository
	du      deffereduploader.IDeferredUploader
	files   blobstore.TempBlobStore
}

type Deps struct {
	Token  string
	ApiKey string
	ChatID int64
	Server merger.MergerServer
	Files  blobstore.TempBlobStore
	Repo   mrepo.MessagesRepository
}
