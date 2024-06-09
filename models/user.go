package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName  string     `json:"firstName" gorm:"not null;default:null"`
	LastName   string     `json:"lastName" gorm:"not null;default:null"`
	Email      string     `json:"email" gorm:"not null;size:255;default:null"`
	Password   string     `json:"password" gorm:"not null;size:255;default:null"`
	Location   string     `json:"location" gorm:"size:255;default:null"`
	Contact    string     `json:"contact" gorm:"not null;size:10;default:null"`
	Role       string     `json:"role" gorm:"not null;enum:'customer','seller';default:null"`
	Properties []Property `json:"properties" gorm:"foreignKey:OwnerID"`
}
