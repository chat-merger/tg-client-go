package deffereduploader

import (
	"fmt"
	"log"
	"merger-adapter/internal/component/debug"
	mrepo "merger-adapter/internal/repository/messages_repository"
	"merger-adapter/internal/service/merger"
)

type ISender interface {
	Send(msg MsgWithKind) error
}

type Sender struct {
	conn merger.Conn
	repo mrepo.MessagesRepository
}

func NewSender(conn merger.Conn, repo mrepo.MessagesRepository) *Sender {
	return &Sender{conn: conn, repo: repo}
}

func (s *Sender) Send(msg MsgWithKind) error {
	msgFromMerger, err := s.conn.Send(msg.msg)
	log.Println("[Sender.Send] s.conn.Send(msg.msg)")
	debug.Print(msg.msg)
	if err != nil {
		return fmt.Errorf("send message to Server: %s", err)
	}
	err = s.repo.Add(mrepo.Message{
		ReplyMergerMsgId: msgFromMerger.ReplyId,
		MergerMsgId:      msgFromMerger.Id,
		ChatId:           msg.original.Chat.Id,
		MsgId:            msg.original.MessageId,
		SenderId:         msg.original.GetSender().Id(),
		SenderFirstName:  msg.original.GetSender().FirstName(),
		Kind:             mrepo.Kind(msg.kind),
		HasMedia:         msg.kind == Media || msg.kind == GroupMedia,
		CreatedAt:        msg.original.Date,
	})

	if err != nil {
		return fmt.Errorf("add message to repo: %s", err)
	}
	return nil
}
