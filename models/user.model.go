package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	FullName     string         `json:"fullName" gorm:"not null;default:null"`
	Email        string         `json:"email" gorm:"not null;size:255;unique;default:null"`
	Password     string         `json:"password" gorm:"not null;size:255;default:null"`
	Location     string         `json:"location" gorm:"size:255;default:null"`
	Phone        string         `json:"phone" gorm:"not null;size:10;default:null"`
	Type         string         `json:"type" gorm:"not null;enum:'customer','gharpeti';default:null"`
	Properties   []Property     `json:"properties" gorm:"foreignKey:OwnerID"`
	Applications []Application  `json:"applications" gorm:"foreignKey:UserID"`
	Token        string         `json:"token"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
