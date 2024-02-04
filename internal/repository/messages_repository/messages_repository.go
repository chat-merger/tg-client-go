package messages_repository

import "merger-adapter/internal/service/merger"

type MessagesRepository interface {
	Get(filter Filter) ([]Message, error)
	Add(message Message) error
}

type Filter struct {
	MergerMsgId      *merger.ID
	ReplyMergerMsgId *merger.ID
	FindNullReply    bool
	ChatId           *int64
	MsgId            *int64
	SenderId         *int64
	SenderFirstName  *string
	Kind             *Kind
	HasMedia         *bool
	CreatedAt        *int64
}

type Message struct {
	MergerMsgId      merger.ID
	ReplyMergerMsgId *merger.ID
	ChatId           int64
	MsgId            int64
	SenderId         int64
	SenderFirstName  string
	Kind             Kind
	HasMedia         bool
	CreatedAt        int64
}

type Kind uint8

const (
	Unknown    Kind = 0
	GroupMedia Kind = 1
	Media      Kind = 2
	Texted     Kind = 3
	Forward    Kind = 4
)
