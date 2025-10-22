package response

import (
	"ims-intro/pkg/domain"
	"time"
)

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func NewErrorResponse(errorMessage string) *ErrorResponse {
	return &ErrorResponse{errorMessage}
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewLoginResponse(token string) *LoginResponse {
	return &LoginResponse{token}
}

type UserProfileResponse struct {
	UserID     int64      `json:"user_id"`
	Username   string     `json:"username"`
	Name       string     `json:"name"`
	Email      *string    `json:"email,omitempty"`
	Role       string     `json:"role"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	LastLogin  *time.Time `json:"last_login,omitempty"`
}

func ToUserProfileResponse(user domain.User) UserProfileResponse {
	return UserProfileResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		LastLogin: user.LastLogin,
	}
}
