package dto

// Terminal Device DTOs

// TerminalDeviceCreate represents the request to register a new Terminal device
type TerminalDeviceCreate struct {
	SquareDeviceID string  `json:"square_device_id" validate:"required"`
	DeviceName     string  `json:"device_name" validate:"required"`
	DeviceCode     *string `json:"device_code,omitempty"`
	LocationID     *string `json:"location_id,omitempty"`
}

// TerminalDeviceUpdate represents the request to update a Terminal device
type TerminalDeviceUpdate struct {
	DeviceName *string `json:"device_name,omitempty"`
	Status     *string `json:"status,omitempty"`
}

// Terminal Checkout DTOs

// TerminalCheckoutCreate represents the request to create a Terminal checkout
type TerminalCheckoutCreate struct {
	SaleID   int64  `json:"sale_id" validate:"required"`
	DeviceID *int64 `json:"device_id,omitempty"` // Use default if not specified
}

// TerminalCheckoutResponse represents the response after creating a checkout
type TerminalCheckoutResponse struct {
	CheckoutID       int64  `json:"checkout_id"`
	SaleID           int64  `json:"sale_id"`
	SquareCheckoutID string `json:"square_checkout_id"`
	Status           string `json:"status"`
	AmountMoney      int64  `json:"amount_money"`
	Currency         string `json:"currency"`
	DeviceID         int64  `json:"device_id"`
	DeviceName       string `json:"device_name"`
}

// Square API Request/Response DTOs

// SquareTerminalCheckoutRequest represents the request to Square Terminal API
type SquareTerminalCheckoutRequest struct {
	IdempotencyKey   string                        `json:"idempotency_key"`
	Checkout         SquareTerminalCheckoutDetails `json:"checkout"`
	DeadlineDuration string                        `json:"deadline_duration,omitempty"`
}

// SquareTerminalCheckoutDetails contains the checkout details
type SquareTerminalCheckoutDetails struct {
	AmountMoney   SquareMoney         `json:"amount_money"`
	DeviceOptions SquareDeviceOptions `json:"device_options"`
	ReferenceID   string              `json:"reference_id,omitempty"`
	Note          string              `json:"note,omitempty"`
}

// SquareMoney represents money in Square API
type SquareMoney struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// SquareDeviceOptions specifies which device to use
type SquareDeviceOptions struct {
	DeviceID string `json:"device_id"`
}

// SquareTerminalCheckoutAPIResponse represents the response from Square Terminal API
type SquareTerminalCheckoutAPIResponse struct {
	Checkout *SquareCheckoutObject `json:"checkout,omitempty"`
	Errors   []SquareError         `json:"errors,omitempty"`
}

// SquareCheckoutObject represents a checkout object from Square
type SquareCheckoutObject struct {
	ID         string   `json:"id"`
	Status     string   `json:"status"`
	PaymentIDs []string `json:"payment_ids,omitempty"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

// SquareError represents an error from Square API
type SquareError struct {
	Category string `json:"category"`
	Code     string `json:"code"`
	Detail   string `json:"detail,omitempty"`
	Field    string `json:"field,omitempty"`
}

// SquareTerminalCheckoutGetResponse represents the response when getting a checkout
type SquareTerminalCheckoutGetResponse struct {
	Checkout *SquareCheckoutObject `json:"checkout,omitempty"`
	Errors   []SquareError         `json:"errors,omitempty"`
}

// SquareTerminalCheckoutCancelResponse represents the response when canceling a checkout
type SquareTerminalCheckoutCancelResponse struct {
	Checkout *SquareCheckoutObject `json:"checkout,omitempty"`
	Errors   []SquareError         `json:"errors,omitempty"`
}

// Webhook DTOs

// SquareWebhookPayload represents the webhook payload from Square
type SquareWebhookPayload struct {
	MerchantID string                `json:"merchant_id"`
	Type       string                `json:"type"`
	EventID    string                `json:"event_id"`
	CreatedAt  string                `json:"created_at"`
	Data       SquareWebhookDataWrap `json:"data"`
}

// SquareWebhookDataWrap wraps the webhook data
type SquareWebhookDataWrap struct {
	Type   string                 `json:"type"`
	ID     string                 `json:"id"`
	Object SquareWebhookDataEvent `json:"object"`
}

// SquareWebhookDataEvent contains the actual event data
type SquareWebhookDataEvent struct {
	Checkout *SquareCheckoutObject `json:"checkout,omitempty"`
	Payment  *SquarePaymentObject  `json:"payment,omitempty"`
}

// SquarePaymentObject represents a payment object from webhook
type SquarePaymentObject struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	TotalMoney   int64  `json:"total_money"`
	ReferenceID  string `json:"reference_id,omitempty"`
	ReceiptURL   string `json:"receipt_url,omitempty"`
	ReceiptNumber string `json:"receipt_number,omitempty"`
}

// Refund DTOs

// SquareMoneyObject represents money in Square API
type SquareMoneyObject struct {
	Amount   int64  `json:"amount"`   // Amount in cents
	Currency string `json:"currency"` // Currency code (e.g., "USD")
}

// RefundCreate represents a request to create a refund
type RefundCreate struct {
	SaleID        int64   `json:"sale_id"`
	AmountCents   *int64  `json:"amount_cents,omitempty"` // Optional for partial refunds
	Reason        *string `json:"reason,omitempty"`
	RefundMessage *string `json:"refund_message,omitempty"` // Custom message for refund
}

// RefundResponse represents the response from Square Refunds API
type RefundResponse struct {
	Refund *SquareRefundObject `json:"refund,omitempty"`
	Errors []SquareError       `json:"errors,omitempty"`
}

// SquareRefundObject represents a refund object from Square API
type SquareRefundObject struct {
	ID          string                   `json:"id"`
	Status      string                   `json:"status"`
	AmountMoney SquareMoneyObject        `json:"amount_money"`
	PaymentID   string                   `json:"payment_id"`
	Reason      string                   `json:"reason,omitempty"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
}

// RefundRequest represents Square Refunds API request
type RefundRequest struct {
	IdempotencyKey string            `json:"idempotency_key"`
	AmountMoney    SquareMoneyObject `json:"amount_money"`
	PaymentID      string            `json:"payment_id"`
	Reason         string            `json:"reason,omitempty"`
}
