package domain

import "time"

// Product represents a product in the inventory system
type Product struct {
	ProductID       int64      `json:"product_id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	SKU             string     `json:"sku"`
	CategoryID      int64      `json:"category_id"`
	BatchNo         *string    `json:"batch_no,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	CostPrice       float64    `json:"cost_price"`
	SellingPrice    float64    `json:"selling_price"`
	CurrentQuantity int64      `json:"current_quantity"`
	MinStockLevel   int64      `json:"min_stock_level"`
	MaxStockLevel   *int64     `json:"max_stock_level,omitempty"`
	ReorderPoint    int64      `json:"reorder_point"`
	Unit            string     `json:"unit"` // pcs, kg, liter, box, carton
	ShelfLocation   *string    `json:"shelf_location,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ProductWithDetails extends Product with computed fields for API responses
type ProductWithDetails struct {
	Product
	CategoryName *string `json:"category_name,omitempty"`
}
