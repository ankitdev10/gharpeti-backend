package handlers

import (
	"errors"
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
	"gorm.io/gorm"
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
		Price:       dto.Price * 100,
		Rooms:       dto.Rooms,
		Location:    dto.Location,
		OwnerID:     user.ID,
		Latitude:    dto.Latitude,
		Longitude:   dto.Longitude,
		Attributes:  dto.Attributes,
	}

	if err := db.DB.Create(&property).Error; err != nil {
		fmt.Println(err)
		return utils.SendError(c, http.StatusInternalServerError, "Error creating Property")
	}

	return utils.SendSuccessResponse(c, http.StatusCreated, "Property created", property)
}

func UpdateProperty(c echo.Context) error {
	dto := c.Get("dto").(*dto.UpdatePropertyDTO)
	id := c.Param("id")

	var property models.Property
	if err := db.DB.First(&property, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.SendError(c, http.StatusNotFound, "Property not found")
		}
		return utils.SendError(c, http.StatusInternalServerError, "Error fetching property")
	}

	user := c.Get("user").(models.User)

	if property.OwnerID != user.ID {
		return utils.SendError(c, http.StatusForbidden, "You are not authorized to update this property")
	}

	property.Title = dto.Title
	property.Description = dto.Description
	property.Price = dto.Price
	property.Rooms = dto.Rooms
	property.Location = dto.Location
	property.Latitude = dto.Latitude
	property.Longitude = dto.Longitude
	property.Enabled = dto.Enabled
	property.Attributes = dto.Attributes

	if err := db.DB.Save(&property).Error; err != nil {
		fmt.Println(err)
		return utils.SendError(c, http.StatusInternalServerError, "Error updating property")
	}

	return utils.SendSuccessResponse(c, http.StatusOK, "Property updated", property)
}

func GetProperties(c echo.Context) error {
	var properties interface{}

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
	result := db.DB.Preload("Owner").Where("id = ?", id).First(&property)
	if result.Error != nil {
		fmt.Println(result.Error)
		return utils.SendError(c, http.StatusInternalServerError, "Internal server error")
	}
	return utils.SendSuccessResponse(c, http.StatusOK, "Found Propery", property)
}

func UploadImage(c echo.Context) error {
	id := c.Param("id")
	propertyID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid property ID",
		})
	}

	var property models.Property
	if err := db.DB.First(&property, propertyID).Error; err != nil {
		return utils.SendError(c, http.StatusNotFound, "not found")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["images"]

	var imageURLs []string

	for _, file := range files {
		imageURL, err := utils.Uploader(file)
		if err != nil {
			return err
		}

		imageURLs = append(imageURLs, imageURL)
	}
	fmt.Println(imageURLs)
	// Update the property images
	property.Images = imageURLs

	if result := db.DB.Save(&property); result.Error != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Can not update property")
	}

	return utils.SendSuccessResponse(c, http.StatusOK, "Uploaded", property)
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

	// only get enabled properties

	conditions = append(conditions, "enabled=true")
	conditions = append(conditions, "deleted_at IS null")
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

func DeletePropery(c echo.Context) error {

	user := c.Get("user").(models.User)
	id := c.Param("id")
	var property models.Property

	res := db.DB.Preload("Owner").First(&property, id)

	if res.Error != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Property doesnt exist")
	}

	if property.OwnerID != user.ID {
		return utils.SendError(c, http.StatusMethodNotAllowed, "You can delete only your property")
	}

	return utils.SendSuccessResponse(c, http.StatusOK, "Successfully deleted", "Successfully Deleted")

}

func GetPropertyOfOwner(c echo.Context) error {
	user := c.Get("user").(models.User)

	var properties []models.Property

	if err := db.DB.Where("owner_id = ?", user.ID).Find(&properties).Error; err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Something went wrong")
	}

	return utils.SendSuccessResponse(c, 200, "Fetched all Properties", properties)
}
