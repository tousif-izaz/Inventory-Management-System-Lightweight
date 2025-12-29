package domain

import "time"

// Purchase represents a purchase order/receipt header
type Purchase struct {
	PurchaseID      int64      `json:"purchase_id"`
	PurchaseDate    time.Time  `json:"purchase_date"`
	SupplierID      int64      `json:"supplier_id"`
	InvoiceNumber   string     `json:"invoice_number"`
	ReferenceNumber *string    `json:"reference_number,omitempty"`
	TotalAmount     float64    `json:"total_amount"`
	TaxAmount       float64    `json:"tax_amount"`
	DiscountAmount  float64    `json:"discount_amount"`
	NetAmount       float64    `json:"net_amount"`
	PaymentStatus   string     `json:"payment_status"` // pending, partial, paid
	PaymentMethod   *string    `json:"payment_method,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	ReceivedBy      *int64     `json:"received_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// PurchaseItem represents a line item in a purchase
type PurchaseItem struct {
	PurchaseItemID  int64      `json:"purchase_item_id"`
	PurchaseID      int64      `json:"purchase_id"`
	ProductID       int64      `json:"product_id"`
	Quantity        int64      `json:"quantity"`
	UnitCost        float64    `json:"unit_cost"`
	TaxRate         float64    `json:"tax_rate"`
	DiscountPercent float64    `json:"discount_percent"`
	LineTotal       float64    `json:"line_total"`
	LocationID      *int64     `json:"location_id,omitempty"`
	BatchNo         *string    `json:"batch_no,omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
}

// PurchaseWithItems represents a purchase with its line items
type PurchaseWithItems struct {
	Purchase
	Items []PurchaseItem `json:"items"`
}
