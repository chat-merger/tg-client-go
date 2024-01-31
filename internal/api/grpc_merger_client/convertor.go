package grpc_side

import (
	"merger-adapter/internal/api/pb"
	"merger-adapter/internal/service/merger"
	"time"
)

func msgToDomain(request *pb.Message) (*merger.Message, error) {

	msg := &merger.Message{
		Id:        merger.ID(request.Id),
		ReplyId:   (*merger.ID)(request.ReplyMsgId),
		Date:      time.Unix(request.CreatedAt, 0),
		Username:  request.Username,
		From:      request.Client,
		Silent:    request.Silent,
		Text:      request.Text,
		Media:     make([]merger.Media, 0, len(request.Media)),
		Forwarded: make([]merger.Forward, 0, len(request.Forwarded)),
	}
	// add media
	for _, it := range request.Media {
		msg.Media = append(msg.Media, mediaToDomain(it))
	}
	for _, it := range request.Forwarded {
		msg.Forwarded = append(msg.Forwarded, forwardToDomain(it))
	}

	return msg, nil
}

func newMsgToPb(msg merger.CreateMessage) (*pb.NewMessageBody, error) {
	var replyMsgId *string
	if msg.ReplyId != nil {
		replyMsgId = (*string)(msg.ReplyId)
	}
	// response
	response := &pb.NewMessageBody{
		CreatedAt:  msg.Date.Unix(),
		Silent:     msg.Silent,
		ReplyMsgId: replyMsgId,
		Username:   msg.Username,
		Text:       msg.Text,
		Media:      make([]*pb.Media, 0, len(msg.Media)),
		Forwarded:  make([]*pb.Forwarded, 0, len(msg.Forwarded)),
	}
	// add media
	for _, it := range msg.Media {
		response.Media = append(response.Media, mediaToPb(it))
	}
	// add forwarded
	for _, it := range msg.Forwarded {
		response.Forwarded = append(response.Forwarded, forwardToPb(it))
	}

	return response, nil
}

func forwardToPb(bm merger.Forward) *pb.Forwarded {
	media := make([]*pb.Media, 0, len(bm.Media))
	for _, it := range bm.Media {
		media = append(media, mediaToPb(it))
	}
	return &pb.Forwarded{
		Id:        (*string)(bm.Id),
		CreatedAt: bm.Date.Unix(),
		Username:  bm.Username,
		Text:      bm.Text,
		Media:     media,
	}
}

func forwardToDomain(bm *pb.Forwarded) merger.Forward {
	media := make([]merger.Media, 0, len(bm.Media))
	for _, it := range bm.Media {
		media = append(media, mediaToDomain(it))
	}
	return merger.Forward{
		Id:       (*merger.ID)(bm.Id),
		Date:     time.Unix(bm.CreatedAt, 0),
		Username: bm.Username,
		Text:     bm.Text,
		Media:    media,
	}
}

func mediaToDomain(bm *pb.Media) merger.Media {
	return merger.Media{
		Kind:    mediaTypeToDomain(bm.Type),
		Spoiler: bm.Spoiler,
		Url:     bm.Url,
	}
}

func mediaToPb(bm merger.Media) *pb.Media {
	return &pb.Media{
		Type:    mediaTypeToPb(bm.Kind),
		Spoiler: bm.Spoiler,
		Url:     bm.Url,
	}
}

func mediaTypeToDomain(kind pb.Media_Type) merger.MediaType {
	var tf merger.MediaType
	switch kind {
	case pb.Media_AUDIO:
		tf = merger.Audio
	case pb.Media_VIDEO:
		tf = merger.Video
	case pb.Media_FILE:
		tf = merger.File
	case pb.Media_PHOTO:
		tf = merger.Photo
	case pb.Media_STICKER:
		tf = merger.Sticker
	}
	return tf
}

func mediaTypeToPb(kind merger.MediaType) pb.Media_Type {
	var tf pb.Media_Type
	switch kind {
	case merger.Audio:
		tf = pb.Media_AUDIO
	case merger.Video:
		tf = pb.Media_VIDEO
	case merger.File:
		tf = pb.Media_FILE
	case merger.Photo:
		tf = pb.Media_PHOTO
	case merger.Sticker:
		tf = pb.Media_STICKER
	}
	return tf
}
