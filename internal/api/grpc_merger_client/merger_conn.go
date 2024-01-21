package grpc_side

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"merger-adapter/internal/api/pb"
	"merger-adapter/internal/debug"
	"merger-adapter/internal/service/merger"
)

type mergerConn struct {
	updates pb.BaseService_UpdatesClient
	send    func(req *pb.Request) (*pb.Response, error)
}

func (s *GrpcMergerClient) Register(xApiKey string) (merger.Conn, error) {
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(authHeader, xApiKey),
	)
	updates, err := s.client.Updates(ctx, &emptypb.Empty{})
	debug.Print(updates)
	if err != nil {
		return nil, fmt.Errorf("client updates: %s", err)
	}
	return &mergerConn{updates: updates, send: func(req *pb.Request) (*pb.Response, error) {
		return s.client.SendMessage(ctx, req)
	}}, nil
}

func (m *mergerConn) Send(data merger.CreateMessage) (*merger.Message, error) {
	req, err := createMessageToRequest(data)
	if err != nil {
		return nil, fmt.Errorf("convertation create message to request: %s", err)
	}
	response, err := m.send(req)
	if err != nil {
		return nil, fmt.Errorf("send message to updates: %s", err)
	}
	message, err := responseToMessage(response)
	if err != nil {
		return nil, fmt.Errorf("convertation response to message: %s", err)
	}
	return message, nil
}

func (m *mergerConn) Update() (*merger.Message, error) {
	response, err := m.updates.Recv()
	if err != nil {
		return nil, fmt.Errorf("receive from updates: %s", err)
	}
	message, err := responseToMessage(response)
	if err != nil {
		return nil, fmt.Errorf("convertation response to message: %s", err)
	}
	return message, nil
}
