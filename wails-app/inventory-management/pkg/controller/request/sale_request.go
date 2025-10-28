package request

import (
	"inventory-management/pkg/service/dto"
	"time"
)

// CreateSaleRequest represents a request to create a sale
type CreateSaleRequest struct {
	SaleDate       *time.Time         `json:"sale_date,omitempty"`
	ReceiptNo      string             `json:"receipt_no" validate:"required"`
	CustomerID     *int64             `json:"customer_id,omitempty"`
	PaymentStatus  string             `json:"payment_status" validate:"required,oneof=pending partial paid refunded"`
	PaymentMethod  *string            `json:"payment_method,omitempty" validate:"omitempty,oneof=cash card mobile bank_transfer credit"`
	Notes          *string            `json:"notes,omitempty"`
	SoldBy         *int64             `json:"sold_by,omitempty"`
	Items          []SaleItemRequest  `json:"items" validate:"required,min=1,dive"`
}

// SaleItemRequest represents a sale line item
type SaleItemRequest struct {
	ProductID       int64   `json:"product_id" validate:"required,gt=0"`
	Quantity        int64   `json:"quantity" validate:"required,gt=0"`
	UnitPrice       float64 `json:"unit_price" validate:"required,gte=0"`
	TaxRate         float64 `json:"tax_rate" validate:"gte=0,lte=100"`
	DiscountPercent float64 `json:"discount_percent" validate:"gte=0,lte=100"`
}

// ToDTO converts request to DTO
func (r *CreateSaleRequest) ToDTO() *dto.SaleCreate {
	saleDate := time.Now()
	if r.SaleDate != nil {
		saleDate = *r.SaleDate
	}

	items := make([]dto.SaleItemCreate, len(r.Items))
	for i, item := range r.Items {
		items[i] = dto.SaleItemCreate{
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			UnitPrice:       item.UnitPrice,
			TaxRate:         item.TaxRate,
			DiscountPercent: item.DiscountPercent,
		}
	}

	return &dto.SaleCreate{
		SaleDate:      saleDate,
		ReceiptNo:     r.ReceiptNo,
		CustomerID:    r.CustomerID,
		PaymentStatus: r.PaymentStatus,
		PaymentMethod: r.PaymentMethod,
		Notes:         r.Notes,
		SoldBy:        r.SoldBy,
		Items:         items,
	}
}

// UpdatePaymentStatusRequest represents a request to update sale payment status
type UpdatePaymentStatusRequest struct {
	PaymentStatus string `json:"payment_status" validate:"required,oneof=pending partial paid refunded"`
}
