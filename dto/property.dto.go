package dto

import (
	"github.com/jackc/pgtype"
)

type CreatePropertyDTO struct {
	Title       string       `json:"title" required:"true"`
	Description string       `json:"description" required:"true"`
	Location    string       `json:"location" validate:"required,max=255"`
	Latitude    float64      `json:"latitude" validate:"required"`
	Longitude   float64      `json:"longitude" validate:"required"`
	Enabled     bool         `json:"enabled"`
	Attributes  pgtype.JSONB `json:"attributes"`
	Images      []string     `json:"images"`
	Rooms       int          `json:"rooms" validate:"required"`
	Price       int          `json:"price" validate:"required"`
}

type UpdatePropertyDTO struct {
	ID          int          `json:"id" required:"true"`
	Title       string       `json:"title" required:"true"`
	Description string       `json:"description" required:"true"`
	Location    string       `json:"location" validate:"required,max=255"`
	Attributes  pgtype.JSONB `json:"attributes"`
	Latitude    float64      `json:"latitude" validate:"required"`
	Longitude   float64      `json:"longitude" validate:"required"`
	Enabled     bool         `json:"enabled"`
	Images      []string     `json:"images"`
	Rooms       int          `json:"rooms" validate:"required"`
	Price       int          `json:"price" validate:"required"`
}
