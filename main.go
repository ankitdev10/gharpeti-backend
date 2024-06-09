package main

import (
	"gharpeti/cmd/db"
	"gharpeti/handlers"
	"gharpeti/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.GET("/", handlers.Home)
	db.InitDB()

	routes.UserRoutes(e)
	routes.PropertyRoutes(e)

	e.Logger.Fatal(e.Start(":4000"))
}
