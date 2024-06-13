package dto

type CreateUserDTO struct {
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Location string `json:"location" validate:"required"`
	Type     string `json:"type" validate:"required"`
}

type UpdateUserDTO struct {
	FullName string `json:"fullName" validate:"omitempty"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty"`
	Phone    string `json:"phone" validate:"omitempty"`
	Location string `json:"location" validate:"omitempty"`
	Type     string `json:"type" validate:"omitempty"`
}
