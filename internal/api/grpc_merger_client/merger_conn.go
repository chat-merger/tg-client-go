package grpc_side

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"merger-adapter/internal/api/pb"
	"merger-adapter/internal/debug"
	"merger-adapter/internal/service/merger"
)

type mergerConn struct {
	conn pb.BaseService_ConnectClient
}

func (s *GrpcMergerClient) Register(xApiKey string) (merger.Conn, error) {
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(authHeader, xApiKey),
	)
	connect, err := s.pbClient.Connect(ctx)
	debug.Print(connect)
	if err != nil {
		return nil, fmt.Errorf("client connect: %s", err)
	}
	return &mergerConn{conn: connect}, nil
}

func (m *mergerConn) Send(data merger.CreateMessage) error {
	req, err := createMessageToRequest(data)
	if err != nil {
		return fmt.Errorf("convertation create message to request: %s", err)
	}
	err = m.conn.Send(req)
	if err != nil {
		return fmt.Errorf("send message to conn: %s", err)
	}
	return nil
}

func (m *mergerConn) Update() (*merger.Message, error) {
	response, err := m.conn.Recv()
	if err != nil {
		return nil, fmt.Errorf("receive from conn: %s", err)
	}
	message, err := responseToMessage(response)
	if err != nil {
		return nil, fmt.Errorf("convertation response to message: %s", err)
	}
	return message, nil
}
