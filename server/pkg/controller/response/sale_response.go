package response

import "time"

// SaleResponse represents a sale in API responses
type SaleResponse struct {
	SaleID         int64      `json:"sale_id"`
	SaleDate       time.Time  `json:"sale_date"`
	ReceiptNo      string     `json:"receipt_no"`
	CustomerID     *int64     `json:"customer_id,omitempty"`
	TotalAmount    float64    `json:"total_amount"`
	TaxAmount      float64    `json:"tax_amount"`
	DiscountAmount float64    `json:"discount_amount"`
	NetAmount      float64    `json:"net_amount"`
	PaymentStatus  string     `json:"payment_status"`
	PaymentMethod  *string    `json:"payment_method,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	SoldBy         *int64     `json:"sold_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// SaleItemResponse represents a sale line item in API responses
type SaleItemResponse struct {
	SaleItemID      int64   `json:"sale_item_id"`
	SaleID          int64   `json:"sale_id"`
	ProductID       int64   `json:"product_id"`
	ProductName     *string `json:"product_name,omitempty"`
	ProductSKU      *string `json:"product_sku,omitempty"`
	Quantity        int64   `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	TaxRate         float64 `json:"tax_rate"`
	DiscountPercent float64 `json:"discount_percent"`
	LineTotal       float64 `json:"line_total"`
}

// SaleWithItemsResponse represents a sale with its line items
type SaleWithItemsResponse struct {
	SaleResponse
	Items []SaleItemResponse `json:"items"`
}
