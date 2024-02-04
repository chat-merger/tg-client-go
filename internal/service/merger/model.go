package merger

import "time"

type Message struct {
	Id        ID
	ReplyId   *ID
	Date      time.Time
	Username  *string
	From      string // client name
	Silent    bool
	Text      *string
	Media     []Media
	Forwarded []Forward
}

type CreateMessage struct {
	ReplyId   *ID
	Date      time.Time
	Username  *string
	Silent    bool
	Text      *string
	Media     []Media
	Forwarded []Forward
}

type Media struct {
	Kind    MediaType
	Spoiler bool
	Url     string
}

type Forward struct {
	Id       *ID
	Date     time.Time
	Username *string
	Text     *string
	Media    []Media
}

type ID string

type MediaType string

const (
	Audio   MediaType = "Audio"
	Video   MediaType = "Video"
	File    MediaType = "File"
	Photo   MediaType = "Photo"
	Sticker MediaType = "Sticker"
)
