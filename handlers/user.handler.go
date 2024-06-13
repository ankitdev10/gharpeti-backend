// handlers/user.go

package handlers

import (
	"errors"
	"fmt"
	"gharpeti/cmd/db"
	"gharpeti/dto"
	"gharpeti/models"
	"gharpeti/utils"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(c echo.Context) error {
	dto := c.Get("dto").(*dto.CreateUserDTO)

	result := db.DB.Where("email = ?", dto.Email).First(&models.User{})

	if result.RowsAffected > 0 {
		return utils.SendError(c, 400, "User already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return utils.SendError(c, http.StatusInternalServerError, "Error hashing password")
	}

	newUser := models.User{
		FullName: dto.FullName,
		Email:    dto.Email,
		Password: string(hashedPassword),
		Location: dto.Location,
		Phone:    dto.Phone,
		Type:     dto.Type,
	}

	newUser.Password = string(hashedPassword)

	if err := db.DB.Create(&newUser).Error; err != nil {
		fmt.Println(err)
		return utils.SendError(c, http.StatusInternalServerError, "Error creating User")
	}

	newUser.Password = ""

	return c.JSON(http.StatusCreated, newUser)
}
func GetUser(c echo.Context) error {
	var users []models.User
	result := db.DB.Find(&users)
	if result.Error != nil {
		fmt.Println(result.Error)
		return utils.SendError(c, http.StatusInternalServerError, "Error fetching users")
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

			return utils.SendError(c, http.StatusInternalServerError, "Error fetching user")
		} else {
			return utils.SendError(c, http.StatusNotFound, "User not found")
		}
	}
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// BUG: Try binding this
func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	dto := c.Get("dto").(*dto.UpdateUserDTO)
	db := db.DB
	var existingUser models.User

	findUser := db.Find(&existingUser, id)
	if findUser.Error != nil {
		if errors.Is(findUser.Error, gorm.ErrRecordNotFound) {
			return utils.SendError(c, http.StatusNotFound, "User not found")
		}
	}

	updatedUser := existingUser

	if err := c.Bind(&dto); err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Can not bind request body")
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
		return utils.SendError(c, 403, "Invalid token")
	}

	userID := uint(userIDFloat64)

	var activeCustomer models.User
	result := db.DB.First(&activeCustomer, userID)

	if result.Error != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to fetch active customer")
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
