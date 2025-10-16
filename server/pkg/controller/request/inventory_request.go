package request

import "ims-intro/pkg/service/dto"

// InventoryAdjustmentRequest represents a request to adjust inventory
type InventoryAdjustmentRequest struct {
	ProductID   int64   `json:"product_id" validate:"required,gt=0"`
	LocationID  int64   `json:"location_id" validate:"required,gt=0"`
	Quantity    int64   `json:"quantity" validate:"required"` // Can be negative
	Reason      string  `json:"reason" validate:"required,oneof=adjustment damage theft return other"`
	Notes       *string `json:"notes,omitempty"`
	PerformedBy int64   `json:"performed_by" validate:"required,gt=0"`
}

// ToDTO converts request to DTO
func (r *InventoryAdjustmentRequest) ToDTO() *dto.InventoryAdjustment {
	return &dto.InventoryAdjustment{
		ProductID:   r.ProductID,
		LocationID:  r.LocationID,
		Quantity:    r.Quantity,
		Reason:      r.Reason,
		Notes:       r.Notes,
		PerformedBy: r.PerformedBy,
	}
}

// InventoryTransferRequest represents a request to transfer inventory between locations
type InventoryTransferRequest struct {
	ProductID      int64   `json:"product_id" validate:"required,gt=0"`
	FromLocationID int64   `json:"from_location_id" validate:"required,gt=0"`
	ToLocationID   int64   `json:"to_location_id" validate:"required,gt=0,nefield=FromLocationID"`
	Quantity       int64   `json:"quantity" validate:"required,gt=0"`
	Reason         *string `json:"reason,omitempty"`
	PerformedBy    int64   `json:"performed_by" validate:"required,gt=0"`
}

// ToDTO converts request to DTO
func (r *InventoryTransferRequest) ToDTO() *dto.InventoryTransfer {
	return &dto.InventoryTransfer{
		ProductID:      r.ProductID,
		FromLocationID: r.FromLocationID,
		ToLocationID:   r.ToLocationID,
		Quantity:       r.Quantity,
		Reason:         r.Reason,
		PerformedBy:    r.PerformedBy,
	}
}
