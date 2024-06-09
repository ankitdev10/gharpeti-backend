package routes

import (
	"gharpeti/handlers"

	"github.com/labstack/echo/v4"
)

func PropertyRoutes(e *echo.Echo) {
	e.POST("/property/create", handlers.CreateProperty)
	e.GET("/property", handlers.GetProperties)
	e.GET("/property/:id", handlers.GetProperty)
	e.GET("/property/search", handlers.Search)
}
