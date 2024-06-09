package handlers

import (
	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	return c.String(200, "Successfully connected to the server")
}
