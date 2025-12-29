package repository

import (
	"database/sql"
	"fmt"
	"ims-intro/pkg/domain"
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

// DeleteSalesInRange deletes all sales and their items within a date range
func (r *SaleRepository) DeleteSalesInRange(fromDate, toDate *time.Time) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// First, get the SaleIDs to be deleted
	selectQuery := "SELECT SaleID FROM Sales WHERE 1=1"
	args := []interface{}{}

	if fromDate != nil {
		selectQuery += " AND SaleDate >= ?"
		args = append(args, *fromDate)
	}
	if toDate != nil {
		selectQuery += " AND SaleDate <= ?"
		args = append(args, *toDate)
	}

	rows, err := tx.Query(selectQuery, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to query sales: %w", err)
	}

	var saleIDs []int64
	for rows.Next() {
		var saleID int64
		if err := rows.Scan(&saleID); err != nil {
			rows.Close()
			return 0, fmt.Errorf("failed to scan sale ID: %w", err)
		}
		saleIDs = append(saleIDs, saleID)
	}
	rows.Close()

	if len(saleIDs) == 0 {
		return 0, nil
	}

	// Delete related transaction records (audit trail) where ReferenceType='sale'
	for _, saleID := range saleIDs {
		_, err := tx.Exec("DELETE FROM Transactions WHERE ReferenceType = 'sale' AND ReferenceID = ?", saleID)
		if err != nil {
			return 0, fmt.Errorf("failed to delete transactions for sale %d: %w", saleID, err)
		}
	}

	// Now delete the sales - CASCADE will handle SalesItems
	deleteQuery := "DELETE FROM Sales WHERE 1=1"
	deleteArgs := []interface{}{}

	if fromDate != nil {
		deleteQuery += " AND SaleDate >= ?"
		deleteArgs = append(deleteArgs, *fromDate)
	}
	if toDate != nil {
		deleteQuery += " AND SaleDate <= ?"
		deleteArgs = append(deleteArgs, *toDate)
	}

	result, err := tx.Exec(deleteQuery, deleteArgs...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete sales: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return int(rowsAffected), nil
}
