package models

import (
	"time"

	"gorm.io/gorm"
)

type Application struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"userID"`
	User          User           `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	PropertyID    uint           `json:"propertyID"`
	Property      Property       `json:"property" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Status        string         `json:"status" gorm:"not null;default:'Pending'; enum:'Pending', 'Approved', 'Rejected' "`
	OfferedPrice  uint           `json:"offeredPrice"`
	ContactNumber uint           `json:"contactNumber"`
	MoveInDate    time.Time      `json:"moveinDate"`
	Feedback      string         `json:"feedback"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
