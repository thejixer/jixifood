package router

import (
	"github.com/labstack/echo/v4"
	"github.com/thejixer/jixifood/services/gateway/internal/handlers"
)

type Router struct {
	e *echo.Echo
	h *handlers.HandlerService
}

func NewRouter(e *echo.Echo, h *handlers.HandlerService) *Router {
	return &Router{
		e: e,
		h: h,
	}
}

func (r *Router) ApplyRoutes(e *echo.Echo) {
	e.GET("/", r.h.HandleHelloWorld)
}
