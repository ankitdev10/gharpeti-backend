package middlewares

import (
	"fmt"
	"gharpeti/cmd/db"
	"gharpeti/models"
	"gharpeti/utils"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func VerifyAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		SECRET := os.Getenv("JWT_SECRET")
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return utils.SendError(c, http.StatusUnauthorized, "Missing or malformed JWT")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.SendError(c, http.StatusUnauthorized, "Invalid authorization header format")
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
			}
			return []byte(SECRET), nil
		})

		if err != nil || !token.Valid {
			return utils.SendError(c, http.StatusUnauthorized, "Invalid or expired JWT")
		}

		claims := token.Claims.(jwt.MapClaims)
		userID, ok := claims["user"].(map[string]interface{})["id"].(float64)
		fmt.Println(userID)
		if !ok {
			return utils.SendError(c, http.StatusUnauthorized, "Invalid token claims")
		}

		var user models.User
		result := db.DB.First(&user, uint(userID))
		if result.Error != nil {
			return utils.SendError(c, http.StatusUnauthorized, "User not found")
		}

		if user.Token != tokenStr {
			return utils.SendError(c, http.StatusUnauthorized, "Invalid token")
		}

		c.Set("user", user)
		return next(c)
	}
}

