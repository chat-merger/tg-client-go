package telegrambot

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log"
)

func InitClient(cfg Config) (*Client, error) {
	bot, err := gotgbot.NewBot(
		cfg.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("building telegram bot: %s", err)
	}

	conn, err := cfg.Server.Register(cfg.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("register client: %s", err)
	}

	disp := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a onMessage, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(disp, nil)

	c := &Client{
		bot:        bot,
		dispatcher: disp,
		updater:    updater,
		conn:       conn,
		chatID:     cfg.ChatID,
	}
	c.gotgbotSetup()

	return c, err
}
