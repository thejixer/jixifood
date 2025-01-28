package server

import (
	"github.com/labstack/echo/v4"
	"github.com/thejixer/jixifood/services/gateway/internal/handlers"
	"github.com/thejixer/jixifood/services/gateway/internal/server/router"
)

type APIServer struct {
	listenAddr     string
	handlerService *handlers.HandlerService
}

func NewAPIServer(listenAddr string, handlerService *handlers.HandlerService) *APIServer {

	return &APIServer{
		listenAddr:     listenAddr,
		handlerService: handlerService,
	}
}

func (s *APIServer) Run() {
	e := echo.New()
	r := router.NewRouter(e, s.handlerService)

	r.ApplyMiddlewares(e)
	r.ApplyRoutes(e)

	e.Logger.Fatal(e.Start(s.listenAddr))
}
