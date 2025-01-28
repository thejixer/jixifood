package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	grpcclient "github.com/thejixer/jixifood/services/gateway/internal/grpc-client"
	"github.com/thejixer/jixifood/shared/models"
)

type HandlerService struct {
	gc *grpcclient.GrpcClient
}

func NewHandlerService(gc *grpcclient.GrpcClient) *HandlerService {
	return &HandlerService{
		gc: gc,
	}
}

func WriteReponse(c echo.Context, s int, m string) error {
	return c.JSON(s, models.ResponseDTO{Msg: m, StatusCode: s})
}

func (h *HandlerService) HandleHelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World, from Jixifood")
}
