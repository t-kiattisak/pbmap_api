package dto

type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	DisplayName string `json:"display_name" validate:"required,min=3,max=50"`
	Role        string `json:"role" validate:"required,oneof=citizen officer admin"`
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name" validate:"omitempty,min=3,max=50"`
	Role        string `json:"role" validate:"omitempty,oneof=citizen officer admin"`
}
