package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"ims-intro/pkg/domain"

	"github.com/labstack/gommon/log"
)

type ITerminalRepository interface {
	// Device Management
	RegisterDevice(device domain.SquareTerminalDevice) (int64, error)
	GetDevices(merchantID string) ([]domain.SquareTerminalDevice, error)
	GetDevice(deviceID int64) (*domain.SquareTerminalDevice, error)
	GetDeviceBySquareID(squareDeviceID string) (*domain.SquareTerminalDevice, error)
	GetDefaultDevice(merchantID string) (*domain.SquareTerminalDevice, error)
	UpdateDevice(deviceID int64, deviceName *string, status *string) error
	DeleteDevice(deviceID int64) error
	SetDefaultDevice(deviceID int64, merchantID string) error

	// Checkout Management
	CreateCheckout(checkout domain.SquareTerminalCheckout) (int64, error)
	GetCheckout(checkoutID int64) (*domain.SquareTerminalCheckout, error)
	GetCheckoutBySquareID(squareCheckoutID string) (*domain.SquareTerminalCheckout, error)
	GetCheckoutBySaleID(saleID int64) (*domain.SquareTerminalCheckout, error)
	UpdateCheckoutStatus(checkoutID int64, status string, paymentID *string, errorCode *string, errorMessage *string) error

	// Refund Management
	CreateRefund(refund domain.SquareRefund) (int64, error)
	GetRefund(refundID int64) (*domain.SquareRefund, error)
	GetRefundBySaleID(saleID int64) (*domain.SquareRefund, error)
	GetRefundBySquareID(squareRefundID string) (*domain.SquareRefund, error)
	UpdateRefundStatus(refundID int64, status string) error
}

type TerminalRepository struct {
	db *sql.DB
}

func NewTerminalRepository(db *sql.DB) ITerminalRepository {
	return &TerminalRepository{db}
}

// RegisterDevice registers a new Terminal device
func (r *TerminalRepository) RegisterDevice(device domain.SquareTerminalDevice) (int64, error) {
	insertStatement := `
		INSERT INTO SquareTerminalDevices (MerchantID, SquareDeviceID, DeviceName, DeviceCode, LocationID, Status, IsDefault)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		insertStatement,
		device.MerchantID,
		device.SquareDeviceID,
		device.DeviceName,
		device.DeviceCode,
		device.LocationID,
		device.Status,
		device.IsDefault,
	)
	if err != nil {
		log.Errorf("error while registering device: %v", err)
		return 0, err
	}

	deviceID, err := result.LastInsertId()
	if err != nil {
		log.Errorf("error getting last insert ID: %v", err)
		return 0, err
	}

	log.Info(fmt.Sprintf("Device registered successfully: %d", deviceID))
	return deviceID, nil
}

// GetDevices retrieves all devices for a merchant
func (r *TerminalRepository) GetDevices(merchantID string) ([]domain.SquareTerminalDevice, error) {
	selectStatement := `
		SELECT DeviceID, MerchantID, SquareDeviceID, DeviceName, DeviceCode, LocationID, Status, IsDefault, CreatedAt, UpdatedAt
		FROM SquareTerminalDevices
		WHERE MerchantID = ?
		ORDER BY IsDefault DESC, DeviceName ASC
	`

	rows, err := r.db.Query(selectStatement, merchantID)
	if err != nil {
		log.Errorf("error while fetching devices: %v", err)
		return nil, err
	}
	defer rows.Close()

	var devices []domain.SquareTerminalDevice
	for rows.Next() {
		var device domain.SquareTerminalDevice
		var deviceCode, locationID sql.NullString

		err := rows.Scan(
			&device.DeviceID,
			&device.MerchantID,
			&device.SquareDeviceID,
			&device.DeviceName,
			&deviceCode,
			&locationID,
			&device.Status,
			&device.IsDefault,
			&device.CreatedAt,
			&device.UpdatedAt,
		)
		if err != nil {
			log.Errorf("error while scanning device: %v", err)
			return nil, err
		}

		if deviceCode.Valid {
			device.DeviceCode = &deviceCode.String
		}
		if locationID.Valid {
			device.LocationID = &locationID.String
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// GetDevice retrieves a device by ID
func (r *TerminalRepository) GetDevice(deviceID int64) (*domain.SquareTerminalDevice, error) {
	var device domain.SquareTerminalDevice
	var deviceCode, locationID sql.NullString

	selectStatement := `
		SELECT DeviceID, MerchantID, SquareDeviceID, DeviceName, DeviceCode, LocationID, Status, IsDefault, CreatedAt, UpdatedAt
		FROM SquareTerminalDevices
		WHERE DeviceID = ?
	`

	err := r.db.QueryRow(selectStatement, deviceID).Scan(
		&device.DeviceID,
		&device.MerchantID,
		&device.SquareDeviceID,
		&device.DeviceName,
		&deviceCode,
		&locationID,
		&device.Status,
		&device.IsDefault,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("device not found")
	}

	if err != nil {
		log.Errorf("error while finding device: %v", err)
		return nil, err
	}

	if deviceCode.Valid {
		device.DeviceCode = &deviceCode.String
	}
	if locationID.Valid {
		device.LocationID = &locationID.String
	}

	return &device, nil
}

// GetDeviceBySquareID retrieves a device by Square device ID
func (r *TerminalRepository) GetDeviceBySquareID(squareDeviceID string) (*domain.SquareTerminalDevice, error) {
	var device domain.SquareTerminalDevice
	var deviceCode, locationID sql.NullString

	selectStatement := `
		SELECT DeviceID, MerchantID, SquareDeviceID, DeviceName, DeviceCode, LocationID, Status, IsDefault, CreatedAt, UpdatedAt
		FROM SquareTerminalDevices
		WHERE SquareDeviceID = ?
	`

	err := r.db.QueryRow(selectStatement, squareDeviceID).Scan(
		&device.DeviceID,
		&device.MerchantID,
		&device.SquareDeviceID,
		&device.DeviceName,
		&deviceCode,
		&locationID,
		&device.Status,
		&device.IsDefault,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("device not found")
	}

	if err != nil {
		log.Errorf("error while finding device by square ID: %v", err)
		return nil, err
	}

	if deviceCode.Valid {
		device.DeviceCode = &deviceCode.String
	}
	if locationID.Valid {
		device.LocationID = &locationID.String
	}

	return &device, nil
}

// GetDefaultDevice retrieves the default device for a merchant
func (r *TerminalRepository) GetDefaultDevice(merchantID string) (*domain.SquareTerminalDevice, error) {
	var device domain.SquareTerminalDevice
	var deviceCode, locationID sql.NullString

	selectStatement := `
		SELECT DeviceID, MerchantID, SquareDeviceID, DeviceName, DeviceCode, LocationID, Status, IsDefault, CreatedAt, UpdatedAt
		FROM SquareTerminalDevices
		WHERE MerchantID = ? AND IsDefault = 1
		LIMIT 1
	`

	err := r.db.QueryRow(selectStatement, merchantID).Scan(
		&device.DeviceID,
		&device.MerchantID,
		&device.SquareDeviceID,
		&device.DeviceName,
		&deviceCode,
		&locationID,
		&device.Status,
		&device.IsDefault,
		&device.CreatedAt,
		&device.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no default device found")
	}

	if err != nil {
		log.Errorf("error while finding default device: %v", err)
		return nil, err
	}

	if deviceCode.Valid {
		device.DeviceCode = &deviceCode.String
	}
	if locationID.Valid {
		device.LocationID = &locationID.String
	}

	return &device, nil
}

// UpdateDevice updates a device's name or status
func (r *TerminalRepository) UpdateDevice(deviceID int64, deviceName *string, status *string) error {
	if deviceName != nil {
		updateStatement := `UPDATE SquareTerminalDevices SET DeviceName = ? WHERE DeviceID = ?`
		_, err := r.db.Exec(updateStatement, *deviceName, deviceID)
		if err != nil {
			log.Errorf("error updating device name: %v", err)
			return err
		}
	}

	if status != nil {
		updateStatement := `UPDATE SquareTerminalDevices SET Status = ? WHERE DeviceID = ?`
		_, err := r.db.Exec(updateStatement, *status, deviceID)
		if err != nil {
			log.Errorf("error updating device status: %v", err)
			return err
		}
	}

	log.Info(fmt.Sprintf("Device updated successfully: %d", deviceID))
	return nil
}

// DeleteDevice deletes a device
func (r *TerminalRepository) DeleteDevice(deviceID int64) error {
	deleteStatement := `DELETE FROM SquareTerminalDevices WHERE DeviceID = ?`

	result, err := r.db.Exec(deleteStatement, deviceID)
	if err != nil {
		log.Errorf("error deleting device: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("error checking rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("device not found")
	}

	log.Info(fmt.Sprintf("Device deleted successfully: %d", deviceID))
	return nil
}

// SetDefaultDevice sets a device as the default for a merchant
func (r *TerminalRepository) SetDefaultDevice(deviceID int64, merchantID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Errorf("error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	// Unset all defaults for this merchant
	_, err = tx.Exec(`UPDATE SquareTerminalDevices SET IsDefault = 0 WHERE MerchantID = ?`, merchantID)
	if err != nil {
		log.Errorf("error unsetting defaults: %v", err)
		return err
	}

	// Set the new default
	result, err := tx.Exec(`UPDATE SquareTerminalDevices SET IsDefault = 1 WHERE DeviceID = ? AND MerchantID = ?`, deviceID, merchantID)
	if err != nil {
		log.Errorf("error setting default: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("error checking rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("device not found")
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("error committing transaction: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Device set as default: %d", deviceID))
	return nil
}

// CreateCheckout creates a new Terminal checkout
func (r *TerminalRepository) CreateCheckout(checkout domain.SquareTerminalCheckout) (int64, error) {
	insertStatement := `
		INSERT INTO SquareTerminalCheckouts (SaleID, MerchantID, DeviceID, SquareCheckoutID, AmountMoney, Currency, Status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		insertStatement,
		checkout.SaleID,
		checkout.MerchantID,
		checkout.DeviceID,
		checkout.SquareCheckoutID,
		checkout.AmountMoney,
		checkout.Currency,
		checkout.Status,
	)
	if err != nil {
		log.Errorf("error creating checkout: %v", err)
		return 0, err
	}

	checkoutID, err := result.LastInsertId()
	if err != nil {
		log.Errorf("error getting last insert ID: %v", err)
		return 0, err
	}

	log.Info(fmt.Sprintf("Checkout created successfully: %d", checkoutID))
	return checkoutID, nil
}

// GetCheckout retrieves a checkout by ID
func (r *TerminalRepository) GetCheckout(checkoutID int64) (*domain.SquareTerminalCheckout, error) {
	var checkout domain.SquareTerminalCheckout
	var paymentID, errorCode, errorMessage sql.NullString
	var completedAt sql.NullTime

	selectStatement := `
		SELECT CheckoutID, SaleID, MerchantID, DeviceID, SquareCheckoutID, AmountMoney, Currency, Status,
		       SquarePaymentID, ErrorCode, ErrorMessage, CreatedAt, UpdatedAt, CompletedAt
		FROM SquareTerminalCheckouts
		WHERE CheckoutID = ?
	`

	err := r.db.QueryRow(selectStatement, checkoutID).Scan(
		&checkout.CheckoutID,
		&checkout.SaleID,
		&checkout.MerchantID,
		&checkout.DeviceID,
		&checkout.SquareCheckoutID,
		&checkout.AmountMoney,
		&checkout.Currency,
		&checkout.Status,
		&paymentID,
		&errorCode,
		&errorMessage,
		&checkout.CreatedAt,
		&checkout.UpdatedAt,
		&completedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("checkout not found")
	}

	if err != nil {
		log.Errorf("error finding checkout: %v", err)
		return nil, err
	}

	if paymentID.Valid {
		checkout.SquarePaymentID = &paymentID.String
	}
	if errorCode.Valid {
		checkout.ErrorCode = &errorCode.String
	}
	if errorMessage.Valid {
		checkout.ErrorMessage = &errorMessage.String
	}
	if completedAt.Valid {
		checkout.CompletedAt = &completedAt.Time
	}

	return &checkout, nil
}

// GetCheckoutBySquareID retrieves a checkout by Square checkout ID
func (r *TerminalRepository) GetCheckoutBySquareID(squareCheckoutID string) (*domain.SquareTerminalCheckout, error) {
	var checkout domain.SquareTerminalCheckout
	var paymentID, errorCode, errorMessage sql.NullString
	var completedAt sql.NullTime

	selectStatement := `
		SELECT CheckoutID, SaleID, MerchantID, DeviceID, SquareCheckoutID, AmountMoney, Currency, Status,
		       SquarePaymentID, ErrorCode, ErrorMessage, CreatedAt, UpdatedAt, CompletedAt
		FROM SquareTerminalCheckouts
		WHERE SquareCheckoutID = ?
	`

	err := r.db.QueryRow(selectStatement, squareCheckoutID).Scan(
		&checkout.CheckoutID,
		&checkout.SaleID,
		&checkout.MerchantID,
		&checkout.DeviceID,
		&checkout.SquareCheckoutID,
		&checkout.AmountMoney,
		&checkout.Currency,
		&checkout.Status,
		&paymentID,
		&errorCode,
		&errorMessage,
		&checkout.CreatedAt,
		&checkout.UpdatedAt,
		&completedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("checkout not found")
	}

	if err != nil {
		log.Errorf("error finding checkout by square ID: %v", err)
		return nil, err
	}

	if paymentID.Valid {
		checkout.SquarePaymentID = &paymentID.String
	}
	if errorCode.Valid {
		checkout.ErrorCode = &errorCode.String
	}
	if errorMessage.Valid {
		checkout.ErrorMessage = &errorMessage.String
	}
	if completedAt.Valid {
		checkout.CompletedAt = &completedAt.Time
	}

	return &checkout, nil
}

// GetCheckoutBySaleID retrieves a checkout by sale ID
func (r *TerminalRepository) GetCheckoutBySaleID(saleID int64) (*domain.SquareTerminalCheckout, error) {
	var checkout domain.SquareTerminalCheckout
	var paymentID, errorCode, errorMessage sql.NullString
	var completedAt sql.NullTime

	selectStatement := `
		SELECT CheckoutID, SaleID, MerchantID, DeviceID, SquareCheckoutID, AmountMoney, Currency, Status,
		       SquarePaymentID, ErrorCode, ErrorMessage, CreatedAt, UpdatedAt, CompletedAt
		FROM SquareTerminalCheckouts
		WHERE SaleID = ?
	`

	err := r.db.QueryRow(selectStatement, saleID).Scan(
		&checkout.CheckoutID,
		&checkout.SaleID,
		&checkout.MerchantID,
		&checkout.DeviceID,
		&checkout.SquareCheckoutID,
		&checkout.AmountMoney,
		&checkout.Currency,
		&checkout.Status,
		&paymentID,
		&errorCode,
		&errorMessage,
		&checkout.CreatedAt,
		&checkout.UpdatedAt,
		&completedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("checkout not found for this sale")
	}

	if err != nil {
		log.Errorf("error finding checkout by sale ID: %v", err)
		return nil, err
	}

	if paymentID.Valid {
		checkout.SquarePaymentID = &paymentID.String
	}
	if errorCode.Valid {
		checkout.ErrorCode = &errorCode.String
	}
	if errorMessage.Valid {
		checkout.ErrorMessage = &errorMessage.String
	}
	if completedAt.Valid {
		checkout.CompletedAt = &completedAt.Time
	}

	return &checkout, nil
}

// UpdateCheckoutStatus updates the status of a checkout
func (r *TerminalRepository) UpdateCheckoutStatus(checkoutID int64, status string, paymentID *string, errorCode *string, errorMessage *string) error {
	updateStatement := `
		UPDATE SquareTerminalCheckouts
		SET Status = ?, SquarePaymentID = ?, ErrorCode = ?, ErrorMessage = ?,
		    CompletedAt = CASE WHEN ? IN ('COMPLETED', 'FAILED', 'CANCELED') THEN CURRENT_TIMESTAMP ELSE CompletedAt END
		WHERE CheckoutID = ?
	`

	_, err := r.db.Exec(updateStatement, status, paymentID, errorCode, errorMessage, status, checkoutID)
	if err != nil {
		log.Errorf("error updating checkout status: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Checkout status updated: %d -> %s", checkoutID, status))
	return nil
}

// ===== Refund Management =====

// CreateRefund creates a new refund record
func (r *TerminalRepository) CreateRefund(refund domain.SquareRefund) (int64, error) {
	insertStatement := `
		INSERT INTO SquareRefunds (sale_id, payment_id, square_refund_id, amount_money, status, reason, refund_message)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		insertStatement,
		refund.SaleID,
		refund.PaymentID,
		refund.SquareRefundID,
		refund.AmountMoney,
		refund.Status,
		refund.Reason,
		refund.RefundMessage,
	)

	if err != nil {
		log.Errorf("error creating refund: %v", err)
		return 0, err
	}

	refundID, err := result.LastInsertId()
	if err != nil {
		log.Errorf("error getting refund ID: %v", err)
		return 0, err
	}

	log.Info(fmt.Sprintf("Refund created: ID=%d, Sale=%d, Amount=%d", refundID, refund.SaleID, refund.AmountMoney))
	return refundID, nil
}

// GetRefund retrieves a refund by ID
func (r *TerminalRepository) GetRefund(refundID int64) (*domain.SquareRefund, error) {
	selectStatement := `
		SELECT id, sale_id, payment_id, square_refund_id, amount_money, status, reason, refund_message, created_at, updated_at, completed_at
		FROM SquareRefunds
		WHERE id = ?
	`

	refund := &domain.SquareRefund{}
	err := r.db.QueryRow(selectStatement, refundID).Scan(
		&refund.RefundID,
		&refund.SaleID,
		&refund.PaymentID,
		&refund.SquareRefundID,
		&refund.AmountMoney,
		&refund.Status,
		&refund.Reason,
		&refund.RefundMessage,
		&refund.CreatedAt,
		&refund.UpdatedAt,
		&refund.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("refund not found")
	}
	if err != nil {
		log.Errorf("error getting refund: %v", err)
		return nil, err
	}

	return refund, nil
}

// GetRefundBySaleID retrieves a refund by sale ID
func (r *TerminalRepository) GetRefundBySaleID(saleID int64) (*domain.SquareRefund, error) {
	selectStatement := `
		SELECT id, sale_id, payment_id, square_refund_id, amount_money, status, reason, refund_message, created_at, updated_at, completed_at
		FROM SquareRefunds
		WHERE sale_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	refund := &domain.SquareRefund{}
	err := r.db.QueryRow(selectStatement, saleID).Scan(
		&refund.RefundID,
		&refund.SaleID,
		&refund.PaymentID,
		&refund.SquareRefundID,
		&refund.AmountMoney,
		&refund.Status,
		&refund.Reason,
		&refund.RefundMessage,
		&refund.CreatedAt,
		&refund.UpdatedAt,
		&refund.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No refund for this sale
	}
	if err != nil {
		log.Errorf("error getting refund by sale ID: %v", err)
		return nil, err
	}

	return refund, nil
}

// GetRefundBySquareID retrieves a refund by Square refund ID
func (r *TerminalRepository) GetRefundBySquareID(squareRefundID string) (*domain.SquareRefund, error) {
	selectStatement := `
		SELECT id, sale_id, payment_id, square_refund_id, amount_money, status, reason, refund_message, created_at, updated_at, completed_at
		FROM SquareRefunds
		WHERE square_refund_id = ?
	`

	refund := &domain.SquareRefund{}
	err := r.db.QueryRow(selectStatement, squareRefundID).Scan(
		&refund.RefundID,
		&refund.SaleID,
		&refund.PaymentID,
		&refund.SquareRefundID,
		&refund.AmountMoney,
		&refund.Status,
		&refund.Reason,
		&refund.RefundMessage,
		&refund.CreatedAt,
		&refund.UpdatedAt,
		&refund.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("refund not found")
	}
	if err != nil {
		log.Errorf("error getting refund by Square ID: %v", err)
		return nil, err
	}

	return refund, nil
}

// UpdateRefundStatus updates the refund status
func (r *TerminalRepository) UpdateRefundStatus(refundID int64, status string) error {
	updateStatement := `
		UPDATE SquareRefunds
		SET status = ?, updated_at = CURRENT_TIMESTAMP,
		    completed_at = CASE WHEN ? IN ('COMPLETED', 'FAILED', 'REJECTED') THEN CURRENT_TIMESTAMP ELSE completed_at END
		WHERE id = ?
	`

	_, err := r.db.Exec(updateStatement, status, status, refundID)
	if err != nil {
		log.Errorf("error updating refund status: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Refund status updated: %d -> %s", refundID, status))
	return nil
}
