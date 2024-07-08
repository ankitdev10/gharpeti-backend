package middlewares

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ValidateDTO(dto interface{}) echo.MiddlewareFunc {
	validate := validator.New()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := c.Bind(dto); err != nil {
				fmt.Println(err)
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
			}

			if err := validate.Struct(dto); err != nil {
				validationErrors := make(map[string]string)
				for _, err := range err.(validator.ValidationErrors) {
					validationErrors[err.Field()] = getValidationErrorMessage(err)
				}

				return c.JSON(http.StatusBadRequest, validationErrors)
			}

			c.Set("dto", dto)

			return next(c)
		}
	}
}

func getValidationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email address"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters long"
	default:
		return err.Field() + " validation failed on " + err.Tag()
	}
}
