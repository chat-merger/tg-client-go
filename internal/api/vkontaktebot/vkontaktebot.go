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
	chatID int64
}

type Config struct {
	Token  string
	ApiKey string
	ChatID int64
	Server merger.MergerServer
}
