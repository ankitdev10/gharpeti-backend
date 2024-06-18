package routes

import (
	"gharpeti/dto"
	"gharpeti/handlers"
	"gharpeti/middlewares"

	"github.com/labstack/echo/v4"
)

func PropertyRoutes(e *echo.Echo) {
	e.POST("/property/create", handlers.CreateProperty, middlewares.ValidateDTO(&dto.CreatePropertyDTO{}), middlewares.VerifyAuth)
	e.GET("/property", handlers.GetProperties)
	e.GET("/property/:id", handlers.GetProperty)
	e.GET("/property/search", handlers.Search)
}
