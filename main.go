package main

import (
	"gharpeti/cmd/db"
	"gharpeti/handlers"
	"gharpeti/routes"
	"log"

	_ "github.com/joho/godotenv/autoload"
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

	// here we create geospatial index
	err := db.DB.Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error
	if err != nil {
		log.Fatal("failed to create PostGIS extension:", err)
	}

	// Run the CREATE INDEX command
	errr := db.DB.Exec("CREATE INDEX idx_properties_location ON properties USING GIST (ST_GeographyFromText('SRID=4326;POINT(' || longitude || ' ' || latitude || ')'))").Error
	if errr != nil {
		log.Fatal("failed to create index:", err)
	}
	routes.UserRoutes(e)
	routes.PropertyRoutes(e)
	routes.ApplicationRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
