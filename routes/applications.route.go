package routes

import (
	"gharpeti/dto"
	"gharpeti/handlers"
	"gharpeti/middlewares"

	"github.com/labstack/echo/v4"
)

func ApplicationRoutes(e *echo.Echo) {
	e.POST("applications/create", handlers.CreateApplication, middlewares.VerifyAuth, middlewares.ValidateDTO(&dto.CreateApplicationDTO{}))
	e.GET("applications/user", handlers.GetUserApplications, middlewares.VerifyAuth)
	e.GET("applications/owner", handlers.GetOwnerApplications, middlewares.VerifyAuth)
	e.PUT("applications/respond/:appId", handlers.RespondToApplication, middlewares.VerifyAuth, middlewares.ValidateDTO(&dto.RespondToApplication{}))
}
