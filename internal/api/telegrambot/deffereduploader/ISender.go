package deffereduploader

import (
	"fmt"
	"log"
	"merger-adapter/internal/api/telegrambot/tghelper"
	"merger-adapter/internal/component/debug"
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
	log.Println("[Sender.Send] s.conn.Send(msg.msg)")
	debug.Print(msg.msg)
	if err != nil {
		return fmt.Errorf("send message to Server: %s", err)
	}
	err = s.mm.Save(tghelper.KvStoreScope, mMsg.Id, toContID(msg.original.MessageId))
	//log.Println("[Sender.Send] s.mm.Save")
	if err != nil {
		return fmt.Errorf("save msg id to MessageMap: %s", err)
	}
	return nil
}
