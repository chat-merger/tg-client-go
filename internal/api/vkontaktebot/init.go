package vkontaktebot

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

func InitClient(deps Deps) (*Client, error) {

	vk := api.NewVK(deps.Token)
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		return nil, fmt.Errorf("GroupsGetByID: %s", err)
	}

	conn, err := deps.Server.Register(deps.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("register client: %s", err)
	}

	// Initializing Long Poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		return nil, fmt.Errorf("NewLongPoll: %s", err)
	}
	c := &Client{
		vk:          vk,
		lp:          lp,
		conn:        conn,
		peerID:      deps.PeerId,
		messagesMap: deps.MessagesMap,
	}

	c.gotgbotSetup()

	return c, err
}
