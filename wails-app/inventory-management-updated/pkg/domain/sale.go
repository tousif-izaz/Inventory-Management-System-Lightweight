package domain

import "time"

// Sale represents a sales transaction header
type Sale struct {
	SaleID         int64      `json:"sale_id"`
	SaleDate       time.Time  `json:"sale_date"`
	ReceiptNo      string     `json:"receipt_no"`
	CustomerID     *int64     `json:"customer_id,omitempty"`
	TotalAmount    float64    `json:"total_amount"`
	TaxAmount      float64    `json:"tax_amount"`
	DiscountAmount float64    `json:"discount_amount"`
	NetAmount      float64    `json:"net_amount"`
	PaymentStatus  string     `json:"payment_status"` // pending, partial, paid, refunded
	PaymentMethod  *string    `json:"payment_method,omitempty"` // cash, card, mobile, bank_transfer, credit
	Notes          *string    `json:"notes,omitempty"`
	SoldBy         *int64     `json:"sold_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// SaleItem represents a line item in a sale
type SaleItem struct {
	SaleItemID      int64   `json:"sale_item_id"`
	SaleID          int64   `json:"sale_id"`
	ProductID       int64   `json:"product_id"`
	Quantity        int64   `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	TaxRate         float64 `json:"tax_rate"`
	DiscountPercent float64 `json:"discount_percent"`
	LineTotal       float64 `json:"line_total"`
	LocationID      *int64  `json:"location_id,omitempty"`
}

// SaleWithItems represents a sale with its line items
type SaleWithItems struct {
	Sale
	Items []SaleItem `json:"items"`
}
