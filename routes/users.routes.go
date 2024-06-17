package routes

import (
	"gharpeti/dto"
	"gharpeti/handlers"
	"gharpeti/middlewares"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo) {
	e.GET("/", handlers.Home)
	e.POST("/users/create", handlers.CreateUser, middlewares.ValidateDTO(&dto.CreateUserDTO{}))
	e.GET("/users", handlers.GetUser)
	e.GET("/users/:id", handlers.GetOneUser)
	e.PUT("/users/update/:id", handlers.UpdateUser, middlewares.VerifyAuth, middlewares.ValidateDTO(&dto.UpdateUserDTO{}))
	e.POST("/auth/login", handlers.Login)
	e.GET("/me", handlers.Me, middlewares.VerifyAuth)
}
