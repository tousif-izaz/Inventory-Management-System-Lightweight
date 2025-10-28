package request

import "inventory-management/pkg/service/dto"

// InventoryAdjustmentRequest represents a request to adjust inventory
type InventoryAdjustmentRequest struct {
	ProductID   int64   `json:"product_id" validate:"required,gt=0"`
	Quantity    int64   `json:"quantity" validate:"required"` // Can be negative
	Reason      string  `json:"reason" validate:"required,oneof=adjustment damage theft return other"`
	Notes       *string `json:"notes,omitempty"`
	PerformedBy int64   `json:"performed_by" validate:"required,gt=0"`
}

// ToDTO converts request to DTO
func (r *InventoryAdjustmentRequest) ToDTO() *dto.InventoryAdjustment {
	return &dto.InventoryAdjustment{
		ProductID:   r.ProductID,
		Quantity:    r.Quantity,
		Reason:      r.Reason,
		Notes:       r.Notes,
		PerformedBy: r.PerformedBy,
	}
}
