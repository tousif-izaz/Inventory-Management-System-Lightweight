package repository

import (
	"database/sql"
	"fmt"
	"inventory-management/pkg/domain"
	"time"
)

type SaleRepository struct {
	db *sql.DB
}

func NewSaleRepository(db *sql.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

// CreateSale creates a new sale record
func (r *SaleRepository) CreateSale(sale *domain.Sale) (int64, error) {
	query := `
		INSERT INTO Sales (SaleDate, ReceiptNo, CustomerID, TotalAmount, TaxAmount,
			DiscountAmount, NetAmount, PaymentStatus, PaymentMethod, Notes, SoldBy)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		sale.SaleDate,
		sale.ReceiptNo,
		sale.CustomerID,
		sale.TotalAmount,
		sale.TaxAmount,
		sale.DiscountAmount,
		sale.NetAmount,
		sale.PaymentStatus,
		sale.PaymentMethod,
		sale.Notes,
		sale.SoldBy,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to create sale: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get sale ID: %w", err)
	}

	return id, nil
}

// CreateSaleItem creates a sale item record
func (r *SaleRepository) CreateSaleItem(item *domain.SaleItem) error {
	query := `
		INSERT INTO SalesItems (SaleID, ProductID, Quantity, UnitPrice,
			TaxRate, DiscountPercent, LineTotal)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		item.SaleID,
		item.ProductID,
		item.Quantity,
		item.UnitPrice,
		item.TaxRate,
		item.DiscountPercent,
		item.LineTotal,
	)

	if err != nil {
		return fmt.Errorf("failed to create sale item: %w", err)
	}

	return nil
}

// GetSaleByID retrieves a sale by ID with all items
func (r *SaleRepository) GetSaleByID(id int64) (*domain.SaleWithItems, error) {
	// Get sale header
	saleQuery := `
		SELECT SaleID, SaleDate, ReceiptNo, CustomerID, TotalAmount, TaxAmount,
			DiscountAmount, NetAmount, PaymentStatus, PaymentMethod, Notes, SoldBy, CreatedAt
		FROM Sales
		WHERE SaleID = ?
	`

	var sale domain.Sale
	err := r.db.QueryRow(saleQuery, id).Scan(
		&sale.SaleID,
		&sale.SaleDate,
		&sale.ReceiptNo,
		&sale.CustomerID,
		&sale.TotalAmount,
		&sale.TaxAmount,
		&sale.DiscountAmount,
		&sale.NetAmount,
		&sale.PaymentStatus,
		&sale.PaymentMethod,
		&sale.Notes,
		&sale.SoldBy,
		&sale.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sale not found")
		}
		return nil, fmt.Errorf("failed to get sale: %w", err)
	}

	// Get sale items
	itemsQuery := `
		SELECT SaleItemID, SaleID, ProductID, Quantity, UnitPrice,
			TaxRate, DiscountPercent, LineTotal
		FROM SalesItems
		WHERE SaleID = ?
	`

	rows, err := r.db.Query(itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get sale items: %w", err)
	}
	defer rows.Close()

	var items []domain.SaleItem
	for rows.Next() {
		var item domain.SaleItem
		err := rows.Scan(
			&item.SaleItemID,
			&item.SaleID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.TaxRate,
			&item.DiscountPercent,
			&item.LineTotal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sale item: %w", err)
		}
		items = append(items, item)
	}

	return &domain.SaleWithItems{
		Sale:  sale,
		Items: items,
	}, nil
}

// GetAllSales retrieves all sales with optional filters
func (r *SaleRepository) GetAllSales(customerID *int64, paymentStatus *string, fromDate, toDate *time.Time) ([]domain.Sale, error) {
	query := `
		SELECT SaleID, SaleDate, ReceiptNo, CustomerID, TotalAmount, TaxAmount,
			DiscountAmount, NetAmount, PaymentStatus, PaymentMethod, Notes, SoldBy, CreatedAt
		FROM Sales
		WHERE 1=1
	`
	args := []interface{}{}

	if customerID != nil {
		query += " AND CustomerID = ?"
		args = append(args, *customerID)
	}

	if paymentStatus != nil {
		query += " AND PaymentStatus = ?"
		args = append(args, *paymentStatus)
	}

	if fromDate != nil {
		query += " AND SaleDate >= ?"
		args = append(args, *fromDate)
	}

	if toDate != nil {
		query += " AND SaleDate <= ?"
		args = append(args, *toDate)
	}

	query += " ORDER BY SaleDate DESC, SaleID DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales: %w", err)
	}
	defer rows.Close()

	var sales []domain.Sale
	for rows.Next() {
		var sale domain.Sale
		err := rows.Scan(
			&sale.SaleID,
			&sale.SaleDate,
			&sale.ReceiptNo,
			&sale.CustomerID,
			&sale.TotalAmount,
			&sale.TaxAmount,
			&sale.DiscountAmount,
			&sale.NetAmount,
			&sale.PaymentStatus,
			&sale.PaymentMethod,
			&sale.Notes,
			&sale.SoldBy,
			&sale.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sale: %w", err)
		}
		sales = append(sales, sale)
	}

	return sales, nil
}

// UpdatePaymentStatus updates the payment status of a sale
func (r *SaleRepository) UpdatePaymentStatus(id int64, status string) error {
	query := `UPDATE Sales SET PaymentStatus = ? WHERE SaleID = ?`

	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("sale not found")
	}

	return nil
}

// GetSalesByCustomerID retrieves all sales for a specific customer
func (r *SaleRepository) GetSalesByCustomerID(customerID int64) ([]domain.Sale, error) {
	query := `
		SELECT SaleID, SaleDate, ReceiptNo, CustomerID, TotalAmount, TaxAmount,
			DiscountAmount, NetAmount, PaymentStatus, PaymentMethod, Notes, SoldBy, CreatedAt
		FROM Sales
		WHERE CustomerID = ?
		ORDER BY SaleDate DESC
	`

	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer sales: %w", err)
	}
	defer rows.Close()

	var sales []domain.Sale
	for rows.Next() {
		var sale domain.Sale
		err := rows.Scan(
			&sale.SaleID,
			&sale.SaleDate,
			&sale.ReceiptNo,
			&sale.CustomerID,
			&sale.TotalAmount,
			&sale.TaxAmount,
			&sale.DiscountAmount,
			&sale.NetAmount,
			&sale.PaymentStatus,
			&sale.PaymentMethod,
			&sale.Notes,
			&sale.SoldBy,
			&sale.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sale: %w", err)
		}
		sales = append(sales, sale)
	}

	return sales, nil
}

// GetSaleItemsBySaleID retrieves all items for a sale
func (r *SaleRepository) GetSaleItemsBySaleID(saleID int64) ([]domain.SaleItem, error) {
	query := `
		SELECT SaleItemID, SaleID, ProductID, Quantity, UnitPrice,
			TaxRate, DiscountPercent, LineTotal
		FROM SalesItems
		WHERE SaleID = ?
	`

	rows, err := r.db.Query(query, saleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sale items: %w", err)
	}
	defer rows.Close()

	var items []domain.SaleItem
	for rows.Next() {
		var item domain.SaleItem
		err := rows.Scan(
			&item.SaleItemID,
			&item.SaleID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.TaxRate,
			&item.DiscountPercent,
			&item.LineTotal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sale item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}
