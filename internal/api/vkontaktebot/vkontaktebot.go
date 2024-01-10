package vkontaktebot

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"merger-adapter/internal/service/merger"
)

type Client struct {
	vk     *api.VK
	lp     *longpoll.LongPoll
	conn   merger.Conn
	peerID int
}

type Config struct {
	Token  string
	ApiKey string
	PeerId int
	Server merger.MergerServer
}
