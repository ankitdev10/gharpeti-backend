package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName   string     `json:"fullName" gorm:"not null;default:null"`
	Email      string     `json:"email" gorm:"not null;size:255;default:null"`
	Password   string     `json:"password" gorm:"not null;size:255;default:null"`
	Location   string     `json:"location" gorm:"size:255;default:null"`
	Phone      string     `json:"phone" gorm:"not null;size:10;default:null"`
	Type       string     `json:"type" gorm:"not null;enum:'customer','gharpeti';default:null"`
	Properties []Property `json:"properties" gorm:"foreignKey:OwnerID"`
}
