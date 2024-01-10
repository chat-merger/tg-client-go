package vkontaktebot

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

func InitClient(cfg Config) (*Client, error) {

	vk := api.NewVK(cfg.Token)
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		return nil, fmt.Errorf("GroupsGetByID: %s", err)
	}

	conn, err := cfg.Server.Register(cfg.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("register client: %s", err)
	}

	// Initializing Long Poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		return nil, fmt.Errorf("NewLongPoll: %s", err)
	}
	c := &Client{
		vk:     vk,
		lp:     lp,
		conn:   conn,
		chatID: cfg.ChatID,
	}

	c.gotgbotSetup()

	return c, err
}
