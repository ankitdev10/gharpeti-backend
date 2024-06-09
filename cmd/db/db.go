package db

import (
	"fmt"
	"log"
	"os"

	"gharpeti/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

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
}
