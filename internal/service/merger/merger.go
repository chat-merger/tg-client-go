package merger

type MergerServer interface {
	Register(xApiKey string) (Conn, error)
}

type Conn interface {
	Send(data CreateMessage) (*Message, error)
	Update() (*Message, error)
}
