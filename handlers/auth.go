package handlers

import (
	"errors"
	"fmt"
	"gharpeti/cmd/db"
	"gharpeti/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c echo.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	SECRET := os.Getenv("JWT_SECRET")

	var user models.User
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&creds); err != nil {
		return err
	}

	if creds.Email == "" || creds.Password == "" {
		return c.JSON(403, "missing credentials")
	}

	result := db.DB.Where("email = ?", creds.Email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Email is not associated with any account.")
		}
		return c.JSON(http.StatusInternalServerError, result.Error)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, "Wrong credentials")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenStr, err := token.SignedString([]byte(SECRET))

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, "Failed to generate token")
	}

	c.Response().Header().Set(echo.HeaderAuthorization, tokenStr)

	return c.JSON(200, user)
}
