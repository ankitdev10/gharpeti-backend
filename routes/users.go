package routes

import (
	"gharpeti/handlers"
	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo) {
	e.GET("/", handlers.Home)
	e.POST("/users/create", handlers.CreateUser)
	e.GET("/users", handlers.GetUser)
	e.GET("/users/:id", handlers.GetOneUser)
	e.PUT("/users/update/:id", handlers.UpdateUser)
	e.POST("/auth/login", handlers.Login)
	e.GET("/activeUser", handlers.ActiveUser, handlers.ValidateToken)
}
