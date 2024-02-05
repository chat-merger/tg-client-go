package messages_repository

import (
	"merger-adapter/internal/api/telegrambot/tghelper"
	"merger-adapter/internal/service/merger"
)

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
	Kind             *tghelper.Kind
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
	Kind             tghelper.Kind
	HasMedia         bool
	CreatedAt        int64
}
