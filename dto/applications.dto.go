package dto

import "time"

type CreateApplicationDTO struct {
	PropertyID    uint      `json:"propertyID" validate:"required"`
	Status        string    `json:"status" validate:"required,oneof=Pending Approved Rejected"`
	OfferedPrice  uint      `json:"offeredPrice" validate:"required,min=0"`
	ContactNumber uint      `json:"contactNumber" validate:"required"`
	MoveInDate    time.Time `json:"moveinDate" validate:"required"`
}

type RespondToApplication struct {
	Status   string `json:"status" validate:"required, oneof=Pending Approved Rejected"`
	Feedback string `json:"feedback" validate:"required"`
}
