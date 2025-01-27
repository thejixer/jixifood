package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HandlerService struct{}

func NewHandlerService() *HandlerService {
	return &HandlerService{}
}

func (h *HandlerService) HandleHelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World, from Jixifood")
}
