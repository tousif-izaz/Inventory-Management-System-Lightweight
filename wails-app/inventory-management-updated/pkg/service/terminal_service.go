package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service/dto"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

const (
	SquareProductionTerminalURL = "https://connect.squareup.com/v2/terminals"
	SquareSandboxTerminalURL    = "https://connect.squareupsandbox.com/v2/terminals"
	SquareProductionBaseURL     = "https://connect.squareup.com"
	SquareSandboxBaseURL        = "https://connect.squareupsandbox.com"
)

type ITerminalService interface {
	// Device Management
	RegisterDevice(merchantID string, deviceDTO dto.TerminalDeviceCreate) (*domain.SquareTerminalDevice, error)
	GetDevices(merchantID string) ([]domain.SquareTerminalDevice, error)
	GetDevice(deviceID int64) (*domain.SquareTerminalDevice, error)
	UpdateDevice(deviceID int64, deviceDTO dto.TerminalDeviceUpdate) error
	DeleteDevice(deviceID int64) error
	SetDefaultDevice(deviceID int64, merchantID string) error

	// Checkout Management
	CreateCheckout(saleID int64, deviceID *int64) (*domain.SquareTerminalCheckout, error)
	GetCheckout(checkoutID int64) (*domain.SquareTerminalCheckout, error)
	GetCheckoutBySaleID(saleID int64) (*domain.SquareTerminalCheckout, error)
	CancelCheckout(checkoutID int64) error
	PollCheckoutStatus(checkoutID int64) (*domain.SquareTerminalCheckout, error)
	UpdateCheckoutFromWebhook(squareCheckoutID string, status string, paymentID *string) error

	// Refund Management
	RefundPayment(saleID int64, amountCents *int64, reason *string, refundMessage *string) (*domain.SquareRefund, error)
	GetRefundBySaleID(saleID int64) (*domain.SquareRefund, error)
}

type TerminalService struct {
	terminalRepo repository.ITerminalRepository
	squareRepo   repository.ISquareRepository
	saleRepo     *repository.SaleRepository
}

func NewTerminalService(
	terminalRepo repository.ITerminalRepository,
	squareRepo repository.ISquareRepository,
	saleRepo *repository.SaleRepository,
) ITerminalService {
	return &TerminalService{
		terminalRepo: terminalRepo,
		squareRepo:   squareRepo,
		saleRepo:     saleRepo,
	}
}

// RegisterDevice registers a new Terminal device
func (s *TerminalService) RegisterDevice(merchantID string, deviceDTO dto.TerminalDeviceCreate) (*domain.SquareTerminalDevice, error) {
	// Check if device already exists
	existing, err := s.terminalRepo.GetDeviceBySquareID(deviceDTO.SquareDeviceID)
	if err == nil && existing != nil {
		return nil, errors.New("device already registered")
	}

	// Create device
	device := domain.SquareTerminalDevice{
		MerchantID:     merchantID,
		SquareDeviceID: deviceDTO.SquareDeviceID,
		DeviceName:     deviceDTO.DeviceName,
		DeviceCode:     deviceDTO.DeviceCode,
		LocationID:     deviceDTO.LocationID,
		Status:         "active",
		IsDefault:      false,
	}

	deviceID, err := s.terminalRepo.RegisterDevice(device)
	if err != nil {
		return nil, err
	}

	device.DeviceID = deviceID
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()

	// If this is the first device, set it as default
	devices, _ := s.terminalRepo.GetDevices(merchantID)
	if len(devices) == 1 {
		s.terminalRepo.SetDefaultDevice(deviceID, merchantID)
		device.IsDefault = true
	}

	return &device, nil
}

// GetDevices retrieves all devices for a merchant
func (s *TerminalService) GetDevices(merchantID string) ([]domain.SquareTerminalDevice, error) {
	return s.terminalRepo.GetDevices(merchantID)
}

// GetDevice retrieves a specific device
func (s *TerminalService) GetDevice(deviceID int64) (*domain.SquareTerminalDevice, error) {
	return s.terminalRepo.GetDevice(deviceID)
}

// UpdateDevice updates a device's properties
func (s *TerminalService) UpdateDevice(deviceID int64, deviceDTO dto.TerminalDeviceUpdate) error {
	return s.terminalRepo.UpdateDevice(deviceID, deviceDTO.DeviceName, deviceDTO.Status)
}

// DeleteDevice deletes a device
func (s *TerminalService) DeleteDevice(deviceID int64) error {
	return s.terminalRepo.DeleteDevice(deviceID)
}

// SetDefaultDevice sets a device as the default
func (s *TerminalService) SetDefaultDevice(deviceID int64, merchantID string) error {
	return s.terminalRepo.SetDefaultDevice(deviceID, merchantID)
}

// CreateCheckout creates a new Terminal checkout
func (s *TerminalService) CreateCheckout(saleID int64, deviceID *int64) (*domain.SquareTerminalCheckout, error) {
	// 1. Get sale details
	sale, err := s.saleRepo.GetSaleByID(saleID)
	if err != nil {
		log.Errorf("error getting sale: %v", err)
		return nil, err
	}

	// 2. Determine which device to use
	var device *domain.SquareTerminalDevice
	if deviceID != nil {
		device, err = s.terminalRepo.GetDevice(*deviceID)
		if err != nil {
			log.Errorf("error getting specified device: %v", err)
			return nil, err
		}
	} else {
		// Use default device for this merchant
		// First, we need to get merchant ID from OAuth token
		// For now, we'll assume we can get it from the first device or configuration
		// In production, you'd get this from the session/auth context
		log.Warn("No device specified, need to implement merchant ID resolution")
		return nil, errors.New("device ID required or merchant context not available")
	}

	// 3. Get merchant OAuth token
	token, err := s.squareRepo.GetOAuthTokenByMerchant(device.MerchantID)
	if err != nil {
		log.Errorf("error getting merchant token: %v", err)
		return nil, err
	}

	// 4. Get Square configuration
	config, err := s.squareRepo.GetConfig()
	if err != nil {
		log.Errorf("error getting Square config: %v", err)
		return nil, err
	}

	// 5. Calculate amount in cents
	amountCents := int64(sale.NetAmount * 100)

	// 6. Create Square Terminal checkout request
	squareRequest := dto.SquareTerminalCheckoutRequest{
		IdempotencyKey: uuid.New().String(),
		Checkout: dto.SquareTerminalCheckoutDetails{
			AmountMoney: dto.SquareMoney{
				Amount:   amountCents,
				Currency: "USD",
			},
			DeviceOptions: dto.SquareDeviceOptions{
				DeviceID: device.SquareDeviceID,
			},
			ReferenceID: fmt.Sprintf("SALE-%d", saleID),
			Note:        fmt.Sprintf("Sale #%s", sale.ReceiptNo),
		},
		DeadlineDuration: "PT5M", // 5 minutes timeout
	}

	// 7. Send request to Square Terminal API
	baseURL := SquareProductionTerminalURL
	if config.Environment == "sandbox" {
		baseURL = SquareSandboxTerminalURL
	}

	squareResponse, err := s.callSquareTerminalCheckoutAPI(baseURL, token.AccessToken, squareRequest)
	if err != nil {
		log.Errorf("error calling Square Terminal API: %v", err)
		return nil, err
	}

	if squareResponse.Checkout == nil {
		log.Errorf("Square API returned no checkout object")
		return nil, errors.New("failed to create checkout: no checkout object returned")
	}

	// 8. Save checkout record to database
	checkout := domain.SquareTerminalCheckout{
		SaleID:           saleID,
		MerchantID:       device.MerchantID,
		DeviceID:         device.DeviceID,
		SquareCheckoutID: squareResponse.Checkout.ID,
		AmountMoney:      amountCents,
		Currency:         "USD",
		Status:           squareResponse.Checkout.Status,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	checkoutID, err := s.terminalRepo.CreateCheckout(checkout)
	if err != nil {
		return nil, err
	}

	checkout.CheckoutID = checkoutID

	log.Info(fmt.Sprintf("Terminal checkout created: %d (Square: %s)", checkoutID, squareResponse.Checkout.ID))
	return &checkout, nil
}

// callSquareTerminalCheckoutAPI makes the HTTP request to Square Terminal API
func (s *TerminalService) callSquareTerminalCheckoutAPI(baseURL string, accessToken string, request dto.SquareTerminalCheckoutRequest) (*dto.SquareTerminalCheckoutAPIResponse, error) {
	url := fmt.Sprintf("%s/checkouts", baseURL)

	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Errorf("error marshaling checkout request: %v", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Errorf("error creating request: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Square-Version", "2024-11-13")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error making request to Square: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var squareResponse dto.SquareTerminalCheckoutAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&squareResponse); err != nil {
		log.Errorf("error decoding Square response: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Square API error (status %d): %v", resp.StatusCode, squareResponse.Errors)
		if len(squareResponse.Errors) > 0 {
			return nil, fmt.Errorf("Square API error: %s - %s", squareResponse.Errors[0].Code, squareResponse.Errors[0].Detail)
		}
		return nil, fmt.Errorf("Square API returned status %d", resp.StatusCode)
	}

	return &squareResponse, nil
}

// GetCheckout retrieves a checkout by ID
func (s *TerminalService) GetCheckout(checkoutID int64) (*domain.SquareTerminalCheckout, error) {
	return s.terminalRepo.GetCheckout(checkoutID)
}

// GetCheckoutBySaleID retrieves a checkout for a sale
func (s *TerminalService) GetCheckoutBySaleID(saleID int64) (*domain.SquareTerminalCheckout, error) {
	return s.terminalRepo.GetCheckoutBySaleID(saleID)
}

// CancelCheckout cancels a Terminal checkout
func (s *TerminalService) CancelCheckout(checkoutID int64) error {
	// 1. Get checkout
	checkout, err := s.terminalRepo.GetCheckout(checkoutID)
	if err != nil {
		return err
	}

	// 2. Get device to get merchant ID
	device, err := s.terminalRepo.GetDevice(checkout.DeviceID)
	if err != nil {
		return err
	}

	// 3. Get merchant token
	token, err := s.squareRepo.GetOAuthTokenByMerchant(device.MerchantID)
	if err != nil {
		return err
	}

	// 4. Get config for environment
	config, err := s.squareRepo.GetConfig()
	if err != nil {
		return err
	}

	baseURL := SquareProductionTerminalURL
	if config.Environment == "sandbox" {
		baseURL = SquareSandboxTerminalURL
	}

	// 5. Call Square API to cancel
	url := fmt.Sprintf("%s/checkouts/%s/cancel", baseURL, checkout.SquareCheckoutID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Errorf("error creating cancel request: %v", err)
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Square-Version", "2024-11-13")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error canceling checkout: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse dto.SquareTerminalCheckoutCancelResponse
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		log.Errorf("Square cancel API error (status %d): %v", resp.StatusCode, errorResponse.Errors)
		return fmt.Errorf("failed to cancel checkout")
	}

	// 6. Update checkout status in database
	return s.terminalRepo.UpdateCheckoutStatus(checkoutID, "CANCELED", nil, nil, nil)
}

// PollCheckoutStatus polls Square API for the latest checkout status
func (s *TerminalService) PollCheckoutStatus(checkoutID int64) (*domain.SquareTerminalCheckout, error) {
	// 1. Get checkout from DB
	checkout, err := s.terminalRepo.GetCheckout(checkoutID)
	if err != nil {
		return nil, err
	}

	// 2. Get device to get merchant ID
	device, err := s.terminalRepo.GetDevice(checkout.DeviceID)
	if err != nil {
		return nil, err
	}

	// 3. Get merchant token
	token, err := s.squareRepo.GetOAuthTokenByMerchant(device.MerchantID)
	if err != nil {
		return nil, err
	}

	// 4. Get config
	config, err := s.squareRepo.GetConfig()
	if err != nil {
		return nil, err
	}

	baseURL := SquareProductionTerminalURL
	if config.Environment == "sandbox" {
		baseURL = SquareSandboxTerminalURL
	}

	// 5. Call Square API to get status
	url := fmt.Sprintf("%s/checkouts/%s", baseURL, checkout.SquareCheckoutID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("error creating get request: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Square-Version", "2024-11-13")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error getting checkout status: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var squareResponse dto.SquareTerminalCheckoutGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&squareResponse); err != nil {
		log.Errorf("error decoding response: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK || squareResponse.Checkout == nil {
		log.Errorf("Square get API error (status %d): %v", resp.StatusCode, squareResponse.Errors)
		return checkout, nil // Return current checkout state
	}

	// 6. Update checkout in database if status changed
	if squareResponse.Checkout.Status != checkout.Status {
		var paymentID *string
		if len(squareResponse.Checkout.PaymentIDs) > 0 {
			paymentID = &squareResponse.Checkout.PaymentIDs[0]
		}

		err = s.terminalRepo.UpdateCheckoutStatus(
			checkoutID,
			squareResponse.Checkout.Status,
			paymentID,
			nil,
			nil,
		)
		if err != nil {
			log.Errorf("error updating checkout status: %v", err)
		}

		checkout.Status = squareResponse.Checkout.Status
		checkout.SquarePaymentID = paymentID
		checkout.UpdatedAt = time.Now()

		if squareResponse.Checkout.Status == "COMPLETED" || squareResponse.Checkout.Status == "FAILED" || squareResponse.Checkout.Status == "CANCELED" {
			now := time.Now()
			checkout.CompletedAt = &now
		}
	}

	return checkout, nil
}

// UpdateCheckoutFromWebhook updates checkout status from webhook event
func (s *TerminalService) UpdateCheckoutFromWebhook(squareCheckoutID string, status string, paymentID *string) error {
	// Get checkout by Square ID
	checkout, err := s.terminalRepo.GetCheckoutBySquareID(squareCheckoutID)
	if err != nil {
		log.Errorf("error getting checkout for webhook: %v", err)
		return err
	}

	// Update status
	return s.terminalRepo.UpdateCheckoutStatus(checkout.CheckoutID, status, paymentID, nil, nil)
}

// ===== Refund Management =====

// RefundPayment processes a refund through Square Payments API
func (s *TerminalService) RefundPayment(saleID int64, amountCents *int64, reason *string, refundMessage *string) (*domain.SquareRefund, error) {
	// 1. Get sale to find payment details
	sale, err := s.saleRepo.GetSaleByID(saleID)
	if err != nil {
		return nil, fmt.Errorf("sale not found: %w", err)
	}

	// 2. Check if this was a Square payment (Terminal or card)
	if sale.PaymentMethod == nil {
		return nil, errors.New("sale has no payment method")
	}

	// Only refund card payments (square_terminal, credit_card, debit_card)
	paymentMethod := *sale.PaymentMethod
	if paymentMethod != "square_terminal" && paymentMethod != "credit_card" && paymentMethod != "debit_card" {
		return nil, fmt.Errorf("refunds only supported for card payments, got: %s", paymentMethod)
	}

	// 3. Get checkout/payment ID
	var paymentID string
	var merchantID string

	if paymentMethod == "square_terminal" {
		// Get from Terminal checkout
		checkout, err := s.terminalRepo.GetCheckoutBySaleID(saleID)
		if err != nil || checkout == nil {
			return nil, errors.New("no Square Terminal checkout found for this sale")
		}
		if checkout.SquarePaymentID == nil {
			return nil, errors.New("Terminal checkout has no payment ID - payment may not be completed")
		}
		paymentID = *checkout.SquarePaymentID
		merchantID = checkout.MerchantID
	} else {
		// For non-terminal card payments, we need payment_id stored somewhere
		// TODO: Add square_payment_id column to Sales table for non-terminal card payments
		return nil, errors.New("non-terminal card payment refunds require payment ID tracking - not yet implemented")
	}

	// 4. Determine refund amount (full or partial)
	refundAmount := sale.NetAmount
	if amountCents != nil {
		refundAmount = float64(*amountCents) / 100.0
	}
	refundAmountCents := int64(refundAmount * 100)

	// Validate refund amount
	saleAmountCents := int64(sale.NetAmount * 100)
	if refundAmountCents > saleAmountCents {
		return nil, fmt.Errorf("refund amount ($%.2f) cannot exceed sale amount ($%.2f)",
			float64(refundAmountCents)/100, sale.NetAmount)
	}

	// 5. Get Square OAuth token
	token, err := s.squareRepo.GetOAuthTokenByMerchant(merchantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Square token: %w", err)
	}

	// 6. Determine Square API base URL
	baseURL := SquareProductionBaseURL
	if s.getSquareEnvironment() == "sandbox" {
		baseURL = SquareSandboxBaseURL
	}

	// 7. Generate idempotency key
	idempotencyKey := uuid.New().String()

	// 8. Prepare refund request
	refundRequest := dto.RefundRequest{
		IdempotencyKey: idempotencyKey,
		AmountMoney: dto.SquareMoneyObject{
			Amount:   refundAmountCents,
			Currency: "USD",
		},
		PaymentID: paymentID,
		Reason:    "Customer refund",
	}
	if reason != nil && *reason != "" {
		refundRequest.Reason = *reason
	}

	requestBody, err := json.Marshal(refundRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refund request: %w", err)
	}

	// 9. Call Square Refunds API
	log.Infof("Creating Square refund: Sale=%d, Payment=%s, Amount=$%.2f",
		saleID, paymentID, float64(refundAmountCents)/100)

	req, err := http.NewRequest("POST", baseURL+"/v2/refunds", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create refund request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Square-Version", "2024-10-17")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Square API: %w", err)
	}
	defer resp.Body.Close()

	// 10. Parse response
	var refundResponse dto.RefundResponse
	if err := json.NewDecoder(resp.Body).Decode(&refundResponse); err != nil {
		return nil, fmt.Errorf("failed to parse Square response: %w", err)
	}

	// Check for errors
	if len(refundResponse.Errors) > 0 {
		errMsg := fmt.Sprintf("%s - %s", refundResponse.Errors[0].Code, refundResponse.Errors[0].Detail)
		log.Errorf("Square API error creating refund: %s", errMsg)
		return nil, fmt.Errorf("Square API error: %s", errMsg)
	}

	if refundResponse.Refund == nil {
		log.Errorf("Square API returned no refund object")
		return nil, errors.New("Square API returned no refund object")
	}

	// 11. Store refund in database
	refund := domain.SquareRefund{
		SaleID:         saleID,
		PaymentID:      paymentID,
		SquareRefundID: refundResponse.Refund.ID,
		AmountMoney:    refundAmountCents,
		Status:         refundResponse.Refund.Status,
		Reason:         reason,
		RefundMessage:  refundMessage,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	refundID, err := s.terminalRepo.CreateRefund(refund)
	if err != nil {
		return nil, fmt.Errorf("failed to store refund: %w", err)
	}

	refund.RefundID = refundID
	log.Infof("Refund created successfully: ID=%d, Square ID=%s, Amount=$%.2f, Status=%s",
		refundID, refund.SquareRefundID, float64(refundAmountCents)/100, refund.Status)

	return &refund, nil
}

// GetRefundBySaleID retrieves refund for a sale
func (s *TerminalService) GetRefundBySaleID(saleID int64) (*domain.SquareRefund, error) {
	return s.terminalRepo.GetRefundBySaleID(saleID)
}

// Helper to get Square environment
func (s *TerminalService) getSquareEnvironment() string {
	// This should come from configuration
	// For now, check if we have sandbox in the base URL constant
	return "production" // Default to production, can be made configurable
}
