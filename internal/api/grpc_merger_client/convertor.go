package grpc_side

import (
	"errors"
	"log"
	"merger-adapter/internal/api/pb"
	"merger-adapter/internal/service/merger"
	"time"
)

func responseToMessage(response *pb.Response) (*merger.Message, error) {
	var body merger.Body
	switch response.Body.(type) {
	case *pb.Response_Text:
		var rt = response.Body.(*pb.Response_Text).Text
		body = &merger.BodyText{
			Format: pbTextFormatToModel(rt.Format),
			Value:  rt.Value,
		}
	case *pb.Response_Media:
		var rm = response.Body.(*pb.Response_Media).Media
		body = &merger.BodyMedia{
			Kind:    pbMediaTypeToModel(rm.Type),
			Caption: rm.Caption,
			Spoiler: rm.Spoiler,
			Url:     rm.Url,
		}
	default:
		return nil, errors.New("response body not match with ResponseBody interface")
	}

	return &merger.Message{
		Id:       merger.ID(response.Id),
		ReplyId:  (*merger.ID)(response.ReplyMsgId),
		Date:     time.Unix(response.CreatedAt, 0),
		Username: response.Username,
		From:     response.Client,
		Silent:   response.Silent,
		Body:     body,
	}, nil
}

func createMessageToRequest(msg merger.CreateMessage) (*pb.Request, error) {
	// response
	response := &pb.Request{
		ReplyMsgId: (*string)(msg.ReplyId),
		CreatedAt:  msg.Date.Unix(),
		Username:   msg.Uername,
		Silent:     msg.Silent,
		Body:       nil, // WithoutBody!!!!!
	}
	// add body
	switch msg.Body.(type) {
	case *merger.BodyText:
		text := msg.Body.(*merger.BodyText)
		response.Body = modelBodyTextToPb(*text)
	case *merger.BodyMedia:
		media := msg.Body.(*merger.BodyMedia)
		response.Body = modelBodyMediaToPb(*media)
	default:
		log.Fatalf("unknown msg.Body:  %#v", msg.Body)
	}
	return response, nil
}

func modelBodyTextToPb(bt merger.BodyText) *pb.Request_Text {
	return &pb.Request_Text{
		Text: &pb.Text{
			Format: modelTextFormatToPbTextFormat(bt.Format),
			Value:  bt.Value,
		},
	}
}

func modelBodyMediaToPb(bm merger.BodyMedia) *pb.Request_Media {
	return &pb.Request_Media{
		Media: &pb.Media{
			Type:    modelMediaTypeToPbMediaType(bm.Kind),
			Caption: bm.Caption,
			Spoiler: bm.Spoiler,
			Url:     bm.Url,
		},
	}
}

func pbTextFormatToModel(format pb.Text_Format) merger.TextFormat {
	var tf merger.TextFormat
	switch format {
	case pb.Text_MARKDOWN:
		tf = merger.Markdown
	case pb.Text_PLAIN:
		tf = merger.Plain
	}
	return tf
}

func pbMediaTypeToModel(kind pb.Media_Type) merger.MediaType {
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

func modelTextFormatToPbTextFormat(format merger.TextFormat) pb.Text_Format {
	var tf pb.Text_Format
	switch format {
	case merger.Markdown:
		tf = pb.Text_MARKDOWN
	case merger.Plain:
		tf = pb.Text_PLAIN
	}
	return tf
}

func modelMediaTypeToPbMediaType(kind merger.MediaType) pb.Media_Type {
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
