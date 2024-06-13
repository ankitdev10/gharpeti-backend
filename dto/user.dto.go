package dto

type CreateUserDTO struct {
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Location string `json:"location" validate:"required"`
	Type     string `json:"type" validate:"required"`
}
