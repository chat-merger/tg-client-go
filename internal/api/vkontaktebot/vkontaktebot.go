package vkontaktebot

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/SevereCloud/vksdk/v2/object"
	"merger-adapter/internal/service/blobstore"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
)

type Client struct {
	vk          *api.VK
	lp          *longpoll.LongPoll
	conn        merger.Conn
	peerID      int
	messagesMap kvstore.MessagesMap
	my          object.GroupsGroup
	files       blobstore.TempBlobStore
}

type Deps struct {
	Token       string
	ApiKey      string
	PeerId      int
	Server      merger.MergerServer
	MessagesMap kvstore.MessagesMap
	Files       blobstore.TempBlobStore
}
