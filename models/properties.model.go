package models

import (
	"time"

	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Property struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null;default:null"`
	Description string         `json:"description" gorm:"not null;default:null"`
	Enabled     bool           `json:"enabled" gorm:"not null;default:true"`
	Location    string         `json:"location" gorm:"not null;size:255;default:null"`
	Images      pq.StringArray `json:"images" gorm:"type:text"`
	Latitude    float64        `json:"latitude" gorm:"not null;default:null"`
	Longitude   float64        `json:"longitude" gorm:"not null;default:null"`
	OwnerID     uint           `json:"ownerId" gorm:"not null"`
	Rooms       int            `json:"rooms" gorm:"not null;default:null"`
	Price       int            `json:"price" gorm:"not null;default:null"`
	Attributes  pgtype.JSONB   `json:"attributes" gorm:"type:jsonb"`
	Owner       User           `json:"owner" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
