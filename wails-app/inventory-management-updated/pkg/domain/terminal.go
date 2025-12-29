package domain

import "time"

// SquareTerminalDevice represents a registered Square Terminal device
type SquareTerminalDevice struct {
	DeviceID       int64     `json:"device_id"`
	MerchantID     string    `json:"merchant_id"`
	SquareDeviceID string    `json:"square_device_id"`
	DeviceName     string    `json:"device_name"`
	DeviceCode     *string   `json:"device_code,omitempty"`
	LocationID     *string   `json:"location_id,omitempty"`
	Status         string    `json:"status"` // active, inactive, paired, unpaired
	IsDefault      bool      `json:"is_default"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// SquareTerminalCheckout represents a Terminal payment checkout
type SquareTerminalCheckout struct {
	CheckoutID       int64      `json:"checkout_id"`
	SaleID           int64      `json:"sale_id"`
	MerchantID       string     `json:"merchant_id"`
	DeviceID         int64      `json:"device_id"`
	SquareCheckoutID string     `json:"square_checkout_id"`
	AmountMoney      int64      `json:"amount_money"` // in cents
	Currency         string     `json:"currency"`
	Status           string     `json:"status"` // PENDING, IN_PROGRESS, COMPLETED, CANCELED, FAILED
	SquarePaymentID  *string    `json:"square_payment_id,omitempty"`
	ErrorCode        *string    `json:"error_code,omitempty"`
	ErrorMessage     *string    `json:"error_message,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
}

// SquareRefund represents a refund transaction for Square payments
type SquareRefund struct {
	RefundID       int64      `json:"refund_id"`
	SaleID         int64      `json:"sale_id"`
	PaymentID      string     `json:"payment_id"`
	SquareRefundID string     `json:"square_refund_id"`
	AmountMoney    int64      `json:"amount_money"` // in cents
	Status         string     `json:"status"`       // PENDING, COMPLETED, FAILED, REJECTED
	Reason         *string    `json:"reason,omitempty"`
	RefundMessage  *string    `json:"refund_message,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}
