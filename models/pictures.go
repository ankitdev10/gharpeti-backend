package models

import "gorm.io/gorm"

type Picture struct {
	gorm.Model
	UserID     int    `json:"userId" gorm:"index"`
	PropertyID int    `json:"propertyId" gorm:"not null"`
	User       User   `gorm:"foreignKey:UserID"`
	URL        string `json:"url" gorm:"not null;size:255"`
}
