package main

import (
	"github.com/joho/godotenv"
	"github.com/thejixer/jixifood/services/gateway/internal/config"
	grpcclient "github.com/thejixer/jixifood/services/gateway/internal/grpc-client"
	"github.com/thejixer/jixifood/services/gateway/internal/handlers"
	"github.com/thejixer/jixifood/services/gateway/internal/server"
)

func init() {
	godotenv.Load()
}

func main() {

	cfg := config.NewGatewayConfig()

	gc := grpcclient.NewGRPCClient(cfg)
	defer gc.Shutdown()

	handlerService := handlers.NewHandlerService(gc)

	s := server.NewAPIServer(cfg.ListenAddr, handlerService)
	s.Run()

}
