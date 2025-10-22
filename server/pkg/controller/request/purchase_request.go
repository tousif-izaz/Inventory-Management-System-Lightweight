package request

import (
	"ims-intro/pkg/service/dto"
	"time"
)

// CreatePurchaseRequest represents a request to create a purchase
type CreatePurchaseRequest struct {
	PurchaseDate    *time.Time             `json:"purchase_date,omitempty"`
	SupplierID      int64                  `json:"supplier_id" validate:"required,gt=0"`
	InvoiceNumber   string                 `json:"invoice_number" validate:"required"`
	ReferenceNumber *string                `json:"reference_number,omitempty"`
	TaxAmount       float64                `json:"tax_amount" validate:"gte=0"`
	DiscountAmount  float64                `json:"discount_amount" validate:"gte=0"`
	PaymentStatus   string                 `json:"payment_status" validate:"required,oneof=pending partial paid"`
	PaymentMethod   *string                `json:"payment_method,omitempty"`
	Notes           *string                `json:"notes,omitempty"`
	ReceivedBy      *int64                 `json:"received_by,omitempty"`
	Items           []PurchaseItemRequest  `json:"items" validate:"required,min=1,dive"`
}

// PurchaseItemRequest represents a purchase line item
type PurchaseItemRequest struct {
	ProductID       int64      `json:"product_id" validate:"required,gt=0"`
	Quantity        int64      `json:"quantity" validate:"required,gt=0"`
	UnitCost        float64    `json:"unit_cost" validate:"required,gte=0"`
	TaxRate         float64    `json:"tax_rate" validate:"gte=0,lte=100"`
	DiscountPercent float64    `json:"discount_percent" validate:"gte=0,lte=100"`
	BatchNo         *string    `json:"batch_no,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
}

// ToDTO converts request to DTO
func (r *CreatePurchaseRequest) ToDTO() *dto.PurchaseCreate {
	purchaseDate := time.Now()
	if r.PurchaseDate != nil {
		purchaseDate = *r.PurchaseDate
	}

	items := make([]dto.PurchaseItemCreate, len(r.Items))
	for i, item := range r.Items {
		items[i] = dto.PurchaseItemCreate{
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			UnitCost:        item.UnitCost,
			TaxRate:         item.TaxRate,
			DiscountPercent: item.DiscountPercent,
			BatchNo:         item.BatchNo,
			ExpiryDate:      item.ExpiryDate,
		}
	}

	return &dto.PurchaseCreate{
		PurchaseDate:    purchaseDate,
		SupplierID:      r.SupplierID,
		InvoiceNumber:   r.InvoiceNumber,
		ReferenceNumber: r.ReferenceNumber,
		TaxAmount:       r.TaxAmount,
		DiscountAmount:  r.DiscountAmount,
		PaymentStatus:   r.PaymentStatus,
		PaymentMethod:   r.PaymentMethod,
		Notes:           r.Notes,
		ReceivedBy:      r.ReceivedBy,
		Items:           items,
	}
}

// UpdatePurchasePaymentStatusRequest represents a request to update purchase payment status
type UpdatePurchasePaymentStatusRequest struct {
	PaymentStatus string `json:"payment_status" validate:"required,oneof=pending partial paid"`
}
