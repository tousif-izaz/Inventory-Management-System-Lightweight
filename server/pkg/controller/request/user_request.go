package request

import "ims-intro/pkg/service/dto"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Username string  `json:"username" validate:"required"`
	Name     string  `json:"name" validate:"required"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password string  `json:"password" validate:"required,min=8"`
	Role     string  `json:"role" validate:"required,oneof=admin manager staff"`
}

func (request *SignUpRequest) ToDtoModel() dto.UserCreate {
	return dto.UserCreate{
		Username: request.Username,
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
	}
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
