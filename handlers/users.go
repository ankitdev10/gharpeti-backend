// handlers/user.go

package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gharpeti/cmd/db"
	"gharpeti/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(c echo.Context) error {
	u := new(models.User)

	if err := c.Bind(u); err != nil {
		fmt.Println(err)
		return err
	}

	// Check if the user already exists
	var existingUser models.User
	result := db.DB.Where("email = ?", u.Email).First(&existingUser)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {

			fmt.Println(result.Error)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Internal server error",
			})
		}
	} else {

		if result.RowsAffected > 0 {

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Email or username is already in use",
			})
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}

	u.Password = string(hashedPassword)

	err = db.DB.Create(u).Error
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}

	u.Password = ""

	return c.JSON(http.StatusCreated, u)
}

func GetUser(c echo.Context) error {
	var users []models.User
	result := db.DB.Find(&users)
	if result.Error != nil {
		fmt.Println(result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}
	return c.JSON(http.StatusOK, users)
}

func GetOneUser(c echo.Context) error {
	var user models.User
	id := c.Param("id")
	db := db.DB
	result := db.First(&user, id)

	if result.Error != nil {

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {

			return c.JSON(http.StatusInternalServerError, result.Error)
		} else {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
	}
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	db := db.DB
	var existingUser models.User

	findUser := db.Find(&existingUser, id)
	if findUser.Error != nil {
		if errors.Is(findUser.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "User not found")
		}
	}

	updatedUser := existingUser

	if err := c.Bind(&updatedUser); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}

	updateResult := db.Model(&existingUser).Updates(&updatedUser)
	if updateResult.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}
	return c.JSON(http.StatusOK, existingUser)
}

func ActiveUser(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userIDFloat64, ok := claims["id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch active customer"})
	}

	userID := uint(userIDFloat64)

	var activeCustomer models.User
	result := db.DB.First(&activeCustomer, userID)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch active customer"})
	}

	activeCustomer.Password = ""

	return c.JSON(http.StatusOK, activeCustomer)
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token format"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !jwtToken.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		c.Set("user", jwtToken)

		return next(c)
	}
}
