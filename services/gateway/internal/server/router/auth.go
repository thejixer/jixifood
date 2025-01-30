package router

import (
	"github.com/labstack/echo/v4"
)

func (r *Router) ApplyAuthRoutes(e *echo.Echo) {
	auth := e.Group("/auth")
	auth.POST("/request-otp", r.h.HandleRequestOTP)
	auth.POST("/verify-otp", r.h.HandleVerifyOTP)
	auth.POST("/me", r.h.HandleME)
	auth.POST("/create-user", r.h.HandleCreateUser)
	auth.POST("/change-userrole", r.h.HandleChangeUserRole)
}
