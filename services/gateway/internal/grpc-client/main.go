package grpcclient

import (
	"log"

	authPB "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/services/gateway/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	AuthClient authPB.AuthServiceClient
	authConn   *grpc.ClientConn
}

func NewGRPCClient(cfg *config.GatewayConfig) *GrpcClient {
	authConn, err := grpc.NewClient(cfg.AuthServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can not initiate the server")
	}
	authClient := authPB.NewAuthServiceClient(authConn)

	return &GrpcClient{
		AuthClient: authClient,
		authConn:   authConn,
	}
}

func (gc *GrpcClient) Shutdown() {
	if gc.authConn != nil {
		gc.authConn.Close()
	}
}
