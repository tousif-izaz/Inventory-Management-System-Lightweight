package domain

import "time"

// User represents a system user
type User struct {
	UserID   int64      `json:"user_id"`
	Username string     `json:"username"`
	Name     string     `json:"name"`
	Email    *string    `json:"email,omitempty"`
	Password string     `json:"-"` // Never expose in JSON
	Role     string     `json:"role"` // admin, manager, staff, viewer
	IsActive bool       `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}
