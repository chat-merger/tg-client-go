package grpc_side

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"merger-adapter/internal/api/pb"
	"merger-adapter/internal/service/merger"
)

type Config struct {
	Host string
}

const authHeader = "X-Api-Key"

var _ merger.MergerServer = (*GrpcMergerClient)(nil)

type GrpcMergerClient struct {
	client pb.BaseServiceClient
}

func InitGrpcMergerClient(cfg Config) (*GrpcMergerClient, error) {
	var cc, err = grpc.Dial(cfg.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("create dial: %s", err)
	}
	cc.ResetConnectBackoff()
	client := pb.NewBaseServiceClient(cc)
	return &GrpcMergerClient{
		client: client,
	}, nil
}
