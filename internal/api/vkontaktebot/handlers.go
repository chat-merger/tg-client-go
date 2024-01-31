package vkontaktebot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
	"log"
	"merger-adapter/internal/service/blobstore"
	"merger-adapter/internal/service/kvstore"
	"merger-adapter/internal/service/merger"
	"strconv"
	"time"
)

func (c *Client) gotgbotSetup() {
	// New message event
	c.lp.MessageNew(c.onMessage)
}

func (c *Client) onMessage(_ context.Context, obj events.MessageNewObject) {
	if obj.Message.PeerID != c.peerID {
		return
	}
	var replyTo *string
	if obj.Message.ReplyMessage != nil {
		id, exists, err := c.messagesMap.ByContID(mmScope, toContID(obj.Message.ReplyMessage.ConversationMessageID))
		if err != nil {
			log.Printf("[ERROR] msg from message map: %s", err)
		}
		if exists {
			replyTo = (*string)(id)
		}
	}

	var author *string

	usrs, _ := c.vk.UsersGet(api.Params{
		"user_ids": obj.Message.FromID,
	})
	if len(usrs) > 0 {
		fname := usrs[0].FirstName + " " + usrs[0].LastName
		author = &fname
	}
	msg := merger.CreateMessage{
		ReplyId:   (*merger.ID)(replyTo),
		Date:      time.Unix(int64(obj.Message.Date), 0),
		Username:  author,
		Silent:    bool(obj.Message.IsSilent),
		Text:      &obj.Message.Text,
		Media:     nil,
		Forwarded: nil,
	}
	mMsg, err := c.conn.Send(msg)
	if err != nil {
		log.Fatalf("send message to Server: %s", err)
	}
	err = c.messagesMap.Save(mmScope, mMsg.Id, toContID(obj.Message.ConversationMessageID))
	if err != nil {
		log.Printf("[ERROR] vkbot onMessage: save msg id to MessageMap: %s", err)
	}
}

const mmScope = kvstore.Scope("VkBotScope")

func toContID(id int) kvstore.ContMsgID {
	return kvstore.ContMsgID(strconv.Itoa(id))
}

func toInt(id kvstore.ContMsgID) int {
	vkMsgId, err := strconv.Atoi(string(id))
	if err != nil {
		log.Fatalf("[ERROR] convert kvstore.ContMsgID to int: %s", err)
	}
	return vkMsgId
}

func (c *Client) listenServerMessages() error {
	for {
		msg, err := c.conn.Update()
		if err != nil {
			return fmt.Errorf("receive update: %s", err)
		}

		b := params.NewMessagesSendBuilder()
		b.Message(msg.FormatShort())
		b.RandomID(0)
		b.PeerIDs([]int{c.peerID})
		// reply
		addReplyIfExists(b, msg, c.messagesMap, c.peerID)
		//https://vk.com/album-224192083_303771670
		//https://vk.com/album-224192083_303406730
		addAttachmentsIfExists(b, c.vk, msg, c.peerID, c.my.ID, c.files)
		vkMsg, err := c.vk.MessagesSendPeerIDs(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		err = c.messagesMap.Save(mmScope, msg.Id, toContID(vkMsg[0].ConversationMessageID))
		if err != nil {
			return fmt.Errorf("save msg id to MessageMap: %s", err)
		}
	}
}
func addAttachmentsIfExists(b *params.MessagesSendBuilder, vk *api.VK, msg *merger.Message, albumId int, myId int, files blobstore.TempBlobStore) {
	attachmentsString := ""

	photos := make([]merger.Media, 0, len(msg.Media))
	for _, media := range msg.Media {
		if media.Kind == merger.Photo {
			photos = append(photos, media)
		}
	}
	if len(photos) != 0 {
		savedPhotos := make(api.PhotosSaveResponse, 0, len(photos))
		// load
		for _, photo := range photos {
			readCloser, err := files.Get(photo.Url)
			if err != nil {
				log.Printf("[ERROR] readCloser from files: %s", err)
				continue

			}
			if readCloser == nil {
				log.Printf("[WARNING] readCloser is nil by uri: %s", photo.Url)
				continue
			}
			uploadServer, err := vk.PhotosGetUploadServer(nil)
			if err != nil {
				log.Printf("[ERROR] photos get upload server: %s", err)
				continue
			}
			bodyContent, err := vk.UploadFile(uploadServer.UploadURL, readCloser, "file1", "file1.jpeg")
			if err != nil {
				return
			}
			if err != nil {
				log.Printf("[ERROR] upload file: %s", err)
				continue
			}
			var handler object.PhotosPhotoUploadResponse

			err = json.Unmarshal(bodyContent, &handler)
			if err != nil {
				return
			}
			if err != nil {
				log.Printf("[ERROR] unmarshal: %s", err)
				continue
			}
			saved, err := vk.PhotosSave(api.Params{
				"server":      handler.Server,
				"photos_list": handler.PhotosList,
				"aid":         handler.AID,
				"hash":        handler.Hash,
				"album_id":    albumId,
			})
			if err != nil {
				log.Printf("[ERROR] photo save: %s", err)
				continue
			}
			savedPhotos = append(savedPhotos, saved...)

		}
		// attach
		for _, photo := range savedPhotos {
			attachmentsString += fmt.Sprintf("photo%d_%d_%s", myId, photo.ID, photo.AccessKey)
		}
	}

	b.Attachment(attachmentsString)
}

func addReplyIfExists(b *params.MessagesSendBuilder, msg *merger.Message, messagesMap kvstore.MessagesMap, peerId int) {
	if msg.ReplyId != nil {
		id, exists, err := messagesMap.ByMergedID(mmScope, *msg.ReplyId)
		if err != nil {
			log.Printf("[ERROR] msg from message map: %s", err)
		}
		log.Printf("messagesMap.ByMergedID: id=%v, exists=%v, err=%s", id, exists, err)
		if exists {
			repPar := map[string]any{
				"conversation_message_ids": toInt(*id),
				"peer_id":                  peerId,
				"is_reply":                 true,
			}
			jsonString, err := json.Marshal(repPar)
			if err != nil {
				log.Printf("[ERROR] marshal params: %s", err)
			}
			b.Forward(string(jsonString))
		}
	}
}
