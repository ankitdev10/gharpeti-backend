package db

import (
	"fmt"
	"os"

	"gharpeti/models"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	ssl := os.Getenv("SSL")
	fmt.Println(ssl)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbHost, dbUser, dbPass, dbName, dbPort, ssl)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	fmt.Println(dsn)
	if err != nil {
		panic(err.Error())
	}
	DB = db

	fmt.Println("Successfully connected to database")

	fmt.Println("-------Running Migrations-----------")
	migrations()
}

func migrations() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println(err)
	}

	err = DB.AutoMigrate(&models.Picture{})
	if err != nil {
		fmt.Println(err)
	}

	err = DB.AutoMigrate(&models.Property{})
	if err != nil {
		fmt.Println(err)
	}

	err = DB.AutoMigrate(&models.Application{})
	if err != nil {
		fmt.Println(err)
	}
}
