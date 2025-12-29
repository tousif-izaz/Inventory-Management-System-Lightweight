package domain

import "time"

// Customer represents a customer for sales tracking
type Customer struct {
	CustomerID    int64      `json:"customer_id"`
	Name          string     `json:"name"`
	Email         *string    `json:"email,omitempty"`
	Phone         *string    `json:"phone,omitempty"`
	Address       *string    `json:"address,omitempty"`
	City          *string    `json:"city,omitempty"`
	Country       *string    `json:"country,omitempty"`
	LoyaltyPoints int64      `json:"loyalty_points"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
