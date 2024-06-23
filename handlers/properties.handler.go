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
	"strings"

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
		return utils.SendError(c, http.StatusInternalServerError, "Internal server error")
	}
	return utils.SendSuccessResponse(c, http.StatusOK, "Found Propery", property)
}

func Search(c echo.Context) error {

	var properties []models.Property

	// filter params
	latStr := c.QueryParam("lat")
	lngStr := c.QueryParam("lng")
	lat, latErr := strconv.ParseFloat(latStr, 64)
	lng, lngErr := strconv.ParseFloat(lngStr, 64)
	radius, err := strconv.ParseInt(c.QueryParam("radius"), 10, 64)
	rooms, err := strconv.ParseInt(c.QueryParam("rooms"), 10, 64)
	minPrice, minPriceErr := strconv.ParseInt(c.QueryParam("minPrice"), 10, 64)
	maxPrice, maxPriceErr := strconv.ParseInt(c.QueryParam("maxPrice"), 10, 64)

	if err != nil || radius <= 0 {
		radius = 10000 // 10 Km
	}
	// pagination
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	var query strings.Builder
	query.WriteString("Select * from properties")

	var conditions []string

	if latErr == nil && lngErr == nil {

		conditions = append(conditions, fmt.Sprintf(`
		ST_DWithin(
			ST_GeographyFromText('POINT(' || longitude || ' ' || latitude || ')'),
			ST_GeographyFromText('POINT(%f %f)'),
			%d
		)
	`, lng, lat, radius))
	}

	if rooms >= 0 {
		conditions = append(conditions, fmt.Sprintf("rooms >= %d", rooms))
	}

	if minPriceErr == nil && maxPriceErr == nil {
		conditions = append(conditions, fmt.Sprintf("price BETWEEN %d AND %d", minPrice*100, maxPrice*100))
	}
	if len(conditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	// adding paginaton

	query.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset))

	if err := db.DB.Raw(query.String()).Find(&properties).Error; err != nil {
		log.Fatalf("Failed to query properties: %v", err)
	}

	return utils.SendSuccessResponse(c, http.StatusOK, "Properties fetched", properties)
}
