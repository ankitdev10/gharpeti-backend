package dto

type CreatePropertyDTO struct {
	Title       string  `json:"title" required:"true"`
	Description string  `json:"description" required:"true"`
	Location    string  `json:"location" validate:"required,max=255"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
	Rooms       int     `json:"rooms" validate:"required"`
	Price       int     `json:"price" validate:"required"`
}
