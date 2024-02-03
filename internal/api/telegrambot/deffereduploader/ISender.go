package deffereduploader

import (
	"fmt"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
)

type ISender interface {
	Send(msg MsgWithKind) error
}

type Sender struct {
	conn merger.Conn
	mm   kvstore.MessagesMap
}

func NewSender(conn merger.Conn, mm kvstore.MessagesMap) *Sender {
	return &Sender{conn: conn, mm: mm}
}

func (s *Sender) Send(msg MsgWithKind) error {
	mMsg, err := s.conn.Send(msg.msg)
	if err != nil {
		return fmt.Errorf("send message to Server: %s", err)
	}
	err = s.mm.Save(mmScope, mMsg.Id, toContID(msg.original.MessageId))
	if err != nil {
		return fmt.Errorf("save msg id to MessageMap: %s", err)
	}
	return nil
}
