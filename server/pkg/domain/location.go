package domain

import "time"

// Location represents a physical storage location
type Location struct {
	LocationID int64      `json:"location_id"`
	Name       string     `json:"name"`
	Type       *string    `json:"type,omitempty"` // warehouse, store, shelf, zone
	Address    *string    `json:"address,omitempty"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
}
