package grpcclient

import (
	"log"

	authPB "github.com/thejixer/jixifood/generated/auth"
	menuPB "github.com/thejixer/jixifood/generated/menu"
	"github.com/thejixer/jixifood/services/gateway/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	AuthClient authPB.AuthServiceClient
	authConn   *grpc.ClientConn
	MenuClient menuPB.MenuServiceClient
	menuConn   *grpc.ClientConn
}

func NewGRPCClient(cfg *config.GatewayConfig) *GrpcClient {
	authConn, err := grpc.NewClient(cfg.AuthServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can not initiate the grpc client")
	}
	authClient := authPB.NewAuthServiceClient(authConn)

	menuConn, err := grpc.NewClient(cfg.MenuServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can not initiate the grpc client")
	}
	menuClient := menuPB.NewMenuServiceClient(menuConn)

	return &GrpcClient{
		AuthClient: authClient,
		authConn:   authConn,
		MenuClient: menuClient,
		menuConn:   menuConn,
	}
}

func (gc *GrpcClient) Shutdown() {
	if gc.authConn != nil {
		gc.authConn.Close()
	}

	if gc.menuConn != nil {
		gc.menuConn.Close()
	}
}
