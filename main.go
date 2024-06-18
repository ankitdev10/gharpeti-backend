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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"auth-token"},
	}))
	e.GET("/", handlers.Home)
	db.InitDB()

	routes.UserRoutes(e)
	routes.PropertyRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
