package handlers

import (
	"fmt"
	"gharpeti/cmd/db"
	"gharpeti/dto"
	"gharpeti/models"
	"gharpeti/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateApplication(c echo.Context) error {

	dto := c.Get("dto").(*dto.CreateApplicationDTO)
	user := c.Get("user").(models.User)

	// check if it is a valid property or not

	var existingProperty models.Property
	if err := db.DB.First(&existingProperty, dto.PropertyID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendError(c, http.StatusNotFound, "Property not found")
		}
		return utils.SendError(c, http.StatusInternalServerError, "Failed to check property existence")
	}

	fmt.Println(existingProperty)

	// check if the user is customer or not

	if user.Type != "customer" {
		return utils.SendError(c, http.StatusMethodNotAllowed, "Only Customer can submit an application")
	}

	// check if the user already has a application for the same property
	var existingApp models.Application
	if err := db.DB.Where("user_id = ? AND property_id = ?", user.ID, dto.PropertyID).First(&existingApp).Error; err == nil {
		return utils.SendError(c, http.StatusConflict, "You already have an application for this property")
	} else if err != gorm.ErrRecordNotFound {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to check existing applications")
	}

	newApp := models.Application{
		UserID:        user.ID,
		PropertyID:    dto.PropertyID,
		Status:        "Pending",
		OfferedPrice:  dto.OfferedPrice * 100,
		ContactNumber: dto.ContactNumber,
		MoveInDate:    dto.MoveInDate,
	}

	if err := db.DB.Create(&newApp).Error; err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to create application")
	}

	return utils.SendSuccessResponse(c, http.StatusOK, "Application submitted", "Sucessfully submitted application")
}

func GetUserApplications(c echo.Context) error {

	user := c.Get("user").(models.User)
	var userApplications []models.Application

	if err := db.DB.Preload("Property").Preload("Property.Owner").Where("user_id = ?", user.ID).Find(&userApplications).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendError(c, http.StatusNotFound, "No applications found for this user")
		}
		return utils.SendError(c, http.StatusInternalServerError, "Failed to fetch user applications")
	}

	return utils.SendSuccessResponse(c, http.StatusOK, "Success", userApplications)
}

func GetOwnerApplications(c echo.Context) error {
	owner := c.Get("user").(models.User)

	if owner.Type != "gharpeti" {
		return utils.SendError(c, http.StatusMethodNotAllowed, "Only Gharpeti can view their applications")
	}

	var apps []models.Application

	if err := db.DB.
		Preload("Property", func(db *gorm.DB) *gorm.DB {
			return db.Where("owner_id = ?", owner.ID)
		}).
		Preload("User").
		Find(&apps).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendError(c, http.StatusNotFound, "No applications found for this user")
		}
		return utils.SendError(c, http.StatusInternalServerError, "Failed to fetch user applications")
	}

	return utils.SendSuccessResponse(c, http.StatusOK, "Success Owner", apps)
}

func RespondToApplication(c echo.Context) error {
	owner := c.Get("user").(models.User)
	appId := c.Param("appId")
	resBody := c.Get("dto").(dto.RespondToApplication)
	if owner.Type != "gharpeti" {
		return utils.SendError(c, http.StatusMethodNotAllowed, "Only Gharpeti can respond to applications")
	}

	var app models.Application

	if err := db.DB.Preload("Property").Preload("Property.Owner").First(&app, appId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendError(c, http.StatusNotFound, "Application Not Found")
		} else {
			return utils.SendError(c, http.StatusInternalServerError, "Something went wrong")
		}
	}

	// check if the property of the application belong to the owner

	if app.Property.OwnerID != owner.ID {
		return utils.SendError(c, http.StatusMethodNotAllowed, "You can respond to your applications only")
	}

	if resBody.Status != "" {
		app.Status = resBody.Status
	}

	if resBody.Status != "" {
		app.Feedback = resBody.Feedback
	}
	if err := db.DB.Save(&app).Error; err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Failed to update application")
	}

	return utils.SendSuccessResponse(c, http.StatusOK, "Success", app)
}
