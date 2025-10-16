package domain

import "time"

// Category represents a product category with hierarchical support
type Category struct {
	CategoryID       int64      `json:"category_id"`
	Name             string     `json:"name"`
	Description      *string    `json:"description,omitempty"`
	ParentCategoryID *int64     `json:"parent_category_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}
