package handlers

import (
	"fmt"
	"gharpeti/cmd/db"
	"gharpeti/dto"
	"gharpeti/models"
	"gharpeti/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CreateProperty(c echo.Context) error {
	dto := c.Get("dto").(*dto.CreatePropertyDTO)

	user := c.Get("user").(models.User)

	fmt.Println(user)

	if user.Type != "gharpeti" {
		return utils.SendError(c, http.StatusForbidden, "You are not authorized to create a property")
	}

	property := models.Property{
		Title:       dto.Title,
		Description: dto.Description,
		Price:       dto.Price,
		Rooms:       dto.Rooms,
		Location:    dto.Location,
		OwnerID:     user.ID,
		Latitude:    dto.Latitude,
		Longitude:   dto.Longitude,
	}

	if err := db.DB.Create(&property).Error; err != nil {
		fmt.Println(err)
		return utils.SendError(c, http.StatusInternalServerError, "Error creating Property")
	}

	return utils.SendSuccessResponse(c, http.StatusCreated, "Property created", property)
}

func GetProperties(c echo.Context) error {
	var properties []models.Property

	result := db.DB.Preload("Owner").Find(&properties)
	if result.Error != nil {
		fmt.Println(result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, properties)
}

func GetProperty(c echo.Context) error {
	var property models.Property
	id := c.Param("id")
	result := db.DB.Preload("Owner").First(&property, id)
	if result.Error != nil {
		fmt.Println(result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, property)
}

func Search(c echo.Context) error {

	var properties []models.Property
	latStr := c.QueryParam("lat")
	lngStr := c.QueryParam("lng")
	lat, _ := strconv.ParseFloat(latStr, 64)
	lng, _ := strconv.ParseFloat(lngStr, 64)

	radius := 100000
	fmt.Println(lat, lng, radius)

	query := fmt.Sprintf(`
		SELECT *
		FROM properties
		WHERE ST_DWithin(
			ST_GeographyFromText('POINT(' || longitude || ' ' || latitude || ')'),
			ST_GeographyFromText('POINT(%f %f)'),
			%d
		)
	`, lng, lat, radius)

	if err := db.DB.Raw(query).Find(&properties).Error; err != nil {
		log.Fatalf("Failed to query properties: %v", err)
	}

	fmt.Println(query)
	return utils.SendSuccessResponse(c, http.StatusOK, "Properties fetched", properties)
}
