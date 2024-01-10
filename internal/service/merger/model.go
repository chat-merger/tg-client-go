package merger

import "time"

type ID string

type Message struct {
	Id      ID
	ReplyId *ID
	Date    time.Time
	Author  *string
	From    string // client name
	Silent  bool
	Body    Body
}

type Body interface{ IsBody() }

type BodyText struct {
	Format TextFormat
	Value  string
}

func (b *BodyText) IsBody() {}

type TextFormat string

const (
	Plain    TextFormat = "Plain"
	Markdown TextFormat = "Markdown"
)

type BodyMedia struct {
	Kind    MediaType
	Caption *string
	Spoiler bool
	Url     string
}

func (b *BodyMedia) IsBody() {}

type MediaType string

const (
	Audio   MediaType = "Audio"
	Video   MediaType = "Video"
	File    MediaType = "File"
	Photo   MediaType = "Photo"
	Sticker MediaType = "Sticker"
)

// create message

type CreateMessage struct {
	ReplyId *ID
	Date    time.Time
	Author  *string
	Silent  bool
	Body    Body
}
