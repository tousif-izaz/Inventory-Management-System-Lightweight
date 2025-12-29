package domain

import "time"

// Transaction represents an audit trail entry for inventory movements
type Transaction struct {
	TransactionID   int64      `json:"transaction_id"`
	TransactionType string     `json:"transaction_type"` // purchase, sale, adjustment, transfer, return, damage, theft
	TransactionDate time.Time  `json:"transaction_date"`
	ProductID       int64      `json:"product_id"`
	LocationID      *int64     `json:"location_id,omitempty"`
	Quantity        int64      `json:"quantity"` // Can be negative for outbound
	ReferenceType   *string    `json:"reference_type,omitempty"` // e.g., "sale", "purchase"
	ReferenceID     *int64     `json:"reference_id,omitempty"` // ID of related transaction
	PreviousQuantity *int64    `json:"previous_quantity,omitempty"`
	NewQuantity     *int64     `json:"new_quantity,omitempty"`
	Reason          *string    `json:"reason,omitempty"`
	PerformedBy     *int64     `json:"performed_by,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}
