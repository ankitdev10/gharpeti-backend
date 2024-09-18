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
	e.PUT("/property/:id", handlers.UpdateProperty, middlewares.ValidateDTO(&dto.UpdatePropertyDTO{}), middlewares.VerifyAuth)
	e.GET("/property/search", handlers.Search)
	e.DELETE("/property/:id", handlers.DeletePropery, middlewares.VerifyAuth)
	e.GET("/property/owner/all", handlers.GetPropertyOfOwner, middlewares.VerifyAuth)
	e.GET("/property/latest", handlers.GetLatestProperty)
	e.POST("/property/upload/:id", handlers.UploadImage, middlewares.VerifyAuth)
}
