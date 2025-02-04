package router

import "github.com/labstack/echo/v4"

func (r *Router) ApplyMenuRoutes(e *echo.Echo) {
	menu := e.Group("/menu")
	menu.POST("/category", r.h.HandleCreateCategory)
	menu.PATCH("/category/:id", r.h.HandleEditCategory)
}
