package main

import (
	"github.com/joho/godotenv"
	"github.com/thejixer/jixifood/services/gateway/internal/config"
	"github.com/thejixer/jixifood/services/gateway/internal/handlers"
	"github.com/thejixer/jixifood/services/gateway/internal/server"
)

func init() {
	godotenv.Load()
}

func main() {

	cfg := config.NewGatewayConfig()

	handlerService := handlers.NewHandlerService()

	s := server.NewAPIServer(cfg.ListenAddr, handlerService)
	s.Run()

}
