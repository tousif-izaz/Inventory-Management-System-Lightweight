package response

import (
	"ims-intro/pkg/domain"
	"time"
)

// ProductResponse represents a product in API responses with additional computed fields
type ProductResponse struct {
	ProductID       int64      `json:"product_id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	SKU             string     `json:"sku"`
	CategoryID      int64      `json:"category_id"`
	CategoryName    *string    `json:"category_name,omitempty"`
	BatchNo         *string    `json:"batch_no,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	CostPrice       float64    `json:"cost_price"`
	SellingPrice    float64    `json:"selling_price"`
	CurrentQuantity int64      `json:"current_quantity"`
	MinStockLevel   int64      `json:"min_stock_level"`
	MaxStockLevel   *int64     `json:"max_stock_level,omitempty"`
	ReorderPoint    int64      `json:"reorder_point"`
	Unit            string     `json:"unit"`
	ShelfLocation   *string    `json:"shelf_location,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func toProductResponse(product *domain.Product) *ProductResponse {
	return &ProductResponse{
		ProductID:       product.ProductID,
		Name:            product.Name,
		Description:     product.Description,
		SKU:             product.SKU,
		CategoryID:      product.CategoryID,
		CategoryName:    nil, // Not available in basic Product
		BatchNo:         product.BatchNo,
		ExpiryDate:      product.ExpiryDate,
		CostPrice:       product.CostPrice,
		SellingPrice:    product.SellingPrice,
		CurrentQuantity: product.CurrentQuantity,
		MinStockLevel:   product.MinStockLevel,
		MaxStockLevel:   product.MaxStockLevel,
		ReorderPoint:    product.ReorderPoint,
		Unit:            product.Unit,
		ShelfLocation:   product.ShelfLocation,
		IsActive:        product.IsActive,
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
	}
}

// ToProductResponse converts a single Product to ProductResponse (public version)
func ToProductResponse(product *domain.Product) *ProductResponse {
	return toProductResponse(product)
}

func toProductWithDetailsResponse(product *domain.ProductWithDetails) *ProductResponse {
	return &ProductResponse{
		ProductID:       product.ProductID,
		Name:            product.Name,
		Description:     product.Description,
		SKU:             product.SKU,
		CategoryID:      product.CategoryID,
		CategoryName:    product.CategoryName,
		BatchNo:         product.BatchNo,
		ExpiryDate:      product.ExpiryDate,
		CostPrice:       product.CostPrice,
		SellingPrice:    product.SellingPrice,
		CurrentQuantity: product.CurrentQuantity,
		MinStockLevel:   product.MinStockLevel,
		MaxStockLevel:   product.MaxStockLevel,
		ReorderPoint:    product.ReorderPoint,
		Unit:            product.Unit,
		ShelfLocation:   product.ShelfLocation,
		IsActive:        product.IsActive,
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
	}
}

func ToProductResponseList(products []*domain.Product) []*ProductResponse {
	responses := make([]*ProductResponse, 0)
	for _, product := range products {
		responses = append(responses, toProductResponse(product))
	}
	return responses
}

func ToProductWithDetailsResponseList(products []*domain.ProductWithDetails) []*ProductResponse {
	responses := make([]*ProductResponse, 0)
	for _, product := range products {
		responses = append(responses, toProductWithDetailsResponse(product))
	}
	return responses
}
