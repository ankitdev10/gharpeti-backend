package handlers

import (
	"errors"
	"fmt"
	"gharpeti/cmd/db"
	"gharpeti/models"
	"gharpeti/utils"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
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
		return utils.SendError(c, http.StatusBadRequest, "Invalid request body")
	}

	result := db.DB.Where("email = ?", creds.Email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.SendError(c, http.StatusNotFound, "Email is not associated with any account.")
		}
		return utils.SendError(c, http.StatusInternalServerError, "Internal server error")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid credentials")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	user.Token = ""
	user.Password = ""

	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenStr, err := token.SignedString([]byte(SECRET))

	if err != nil {
		fmt.Println(err)
		return utils.SendError(c, http.StatusInternalServerError, "Internal server error")
	}
	user.Token = tokenStr
	db.DB.Save(&user)
	c.Response().Header().Set("Auth-Token", tokenStr)
	user.Password = ""
	user.Token = ""
	return c.JSON(200, user)
}
