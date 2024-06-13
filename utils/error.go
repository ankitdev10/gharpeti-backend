package utils

import "github.com/labstack/echo/v4"

type ErroResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SendError(c echo.Context, status int, message string) error {
	err := ErroResponse{
		Status:  status,
		Message: message,
	}

	return c.JSON(status, err)
}
