package telegrambot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"merger-adapter/internal/service/merger"
	"time"
)

func (d *DeferredUploaderImpl) msgToMerger(msg *gotgbot.Message, text *string, media []merger.Media, forward []merger.Forward) *merger.Message {
	uname := msg.GetSender().FirstName()
	return &merger.Message{
		ReplyId:   (*merger.ID)(replyMergerIdFromMsg(msg, d.mm)),
		Date:      time.Unix(msg.Date, 0),
		Username:  &uname,
		Silent:    false, // where prop??
		Text:      text,
		Media:     media,
		Forwarded: forward,
	}
}
