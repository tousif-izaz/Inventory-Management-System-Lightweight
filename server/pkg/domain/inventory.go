package domain

import "time"

// Inventory represents stock levels at a specific location
type Inventory struct {
	InventoryID       int64     `json:"inventory_id"`
	ProductID         int64     `json:"product_id"`
	LocationID        int64     `json:"location_id"`
	Quantity          int64     `json:"quantity"`
	ReservedQuantity  int64     `json:"reserved_quantity"`
	AvailableQuantity int64     `json:"available_quantity"` // Virtual/computed field
	LastUpdated       time.Time `json:"last_updated"`
}
