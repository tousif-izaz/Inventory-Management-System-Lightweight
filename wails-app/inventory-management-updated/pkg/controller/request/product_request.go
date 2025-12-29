package request

import (
	"ims-intro/pkg/service/dto"
	"time"
)

// AddProductRequest represents a request to create a new product
type AddProductRequest struct {
	Name            string     `json:"name" validate:"required"`
	Description     *string    `json:"description,omitempty"`
	SKU             string     `json:"sku" validate:"required"`
	CategoryID      int64      `json:"category_id" validate:"required,gt=0"`
	BatchNo         *string    `json:"batch_no,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	CostPrice       float64    `json:"cost_price" validate:"required,gte=0"`
	SellingPrice    float64    `json:"selling_price" validate:"required,gte=0"`
	CurrentQuantity int64      `json:"current_quantity" validate:"gte=0"`
	MinStockLevel   int64      `json:"min_stock_level" validate:"gte=0"`
	MaxStockLevel   *int64     `json:"max_stock_level,omitempty" validate:"omitempty,gte=0"`
	ReorderPoint    int64      `json:"reorder_point" validate:"gte=0"`
	Unit            string     `json:"unit" validate:"required,oneof=pcs kg liter box carton"`
	ShelfLocation   *string    `json:"shelf_location,omitempty"`
}

// ToModel converts request to DTO
func (request *AddProductRequest) ToModel() *dto.ProductCreate {
	return &dto.ProductCreate{
		Name:            request.Name,
		Description:     request.Description,
		SKU:             request.SKU,
		CategoryID:      request.CategoryID,
		BatchNo:         request.BatchNo,
		ExpiryDate:      request.ExpiryDate,
		CostPrice:       request.CostPrice,
		SellingPrice:    request.SellingPrice,
		CurrentQuantity: request.CurrentQuantity,
		MinStockLevel:   request.MinStockLevel,
		MaxStockLevel:   request.MaxStockLevel,
		ReorderPoint:    request.ReorderPoint,
		Unit:            request.Unit,
		ShelfLocation:   request.ShelfLocation,
	}
}
