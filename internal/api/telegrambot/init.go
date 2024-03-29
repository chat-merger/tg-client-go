package telegrambot

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"log"
	"merger-adapter/internal/api/telegrambot/deffereduploader"
)

func InitClient(deps Deps) (*Client, error) {
	bot, err := gotgbot.NewBot(
		deps.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("building telegram bot: %s", err)
	}

	conn, err := deps.Server.Register(deps.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("register client: %s", err)
	}

	disp := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a onTelegramCreatedMessage, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(disp, nil)

	c := &Client{
		bot:     bot,
		updater: updater,
		conn:    conn,
		chatID:  deps.ChatID,
		repo:    deps.Repo,
		du: deffereduploader.NewDeferredUploader(
			deps.Repo,
			deps.Files,
			bot,
			conn,
		),
		files: deps.Files,
	}
	disp.AddHandler(handlers.NewMessage(c.filter, c.onTelegramCreatedMessage))

	return c, err
}
