package repository

import (
	"database/sql"
	"fmt"

	"github.com/labstack/gommon/log"
	"ims-intro/pkg/domain"
)

// ITransactionRepository defines transaction repository interface
type ITransactionRepository interface {
	Create(transaction *domain.Transaction) error
	GetByID(transactionID int64) (*domain.Transaction, error)
	GetAll(filters map[string]interface{}) ([]*domain.Transaction, error)
	GetByProduct(productID int64) ([]*domain.Transaction, error)
	GetByLocation(locationID int64) ([]*domain.Transaction, error)
	GetByUser(userID int64) ([]*domain.Transaction, error)
	GetByType(transactionType string) ([]*domain.Transaction, error)
}

type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) ITransactionRepository {
	return &TransactionRepository{db}
}

// Create creates a new transaction record
func (r *TransactionRepository) Create(transaction *domain.Transaction) error {
	query := `
		INSERT INTO Transactions (
			TransactionType, TransactionDate, ProductID, LocationID,
			Quantity, ReferenceType, ReferenceID, PreviousQuantity,
			NewQuantity, Reason, PerformedBy, Notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		transaction.TransactionType,
		transaction.TransactionDate,
		transaction.ProductID,
		transaction.LocationID,
		transaction.Quantity,
		transaction.ReferenceType,
		transaction.ReferenceID,
		transaction.PreviousQuantity,
		transaction.NewQuantity,
		transaction.Reason,
		transaction.PerformedBy,
		transaction.Notes,
	)

	if err != nil {
		log.Errorf("error creating transaction: %v", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	transaction.TransactionID = id
	log.Infof("Transaction created: ID=%d, Type=%s, Product=%d, Quantity=%d",
		id, transaction.TransactionType, transaction.ProductID, transaction.Quantity)
	return nil
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(transactionID int64) (*domain.Transaction, error) {
	query := `
		SELECT TransactionID, TransactionType, TransactionDate, ProductID, LocationID,
		       Quantity, ReferenceType, ReferenceID, PreviousQuantity, NewQuantity,
		       Reason, PerformedBy, Notes, CreatedAt
		FROM Transactions
		WHERE TransactionID = ?
	`

	row := r.db.QueryRow(query, transactionID)
	return scanTransactionRow(row)
}

// GetAll retrieves all transactions with optional filters
func (r *TransactionRepository) GetAll(filters map[string]interface{}) ([]*domain.Transaction, error) {
	query := `
		SELECT TransactionID, TransactionType, TransactionDate, ProductID, LocationID,
		       Quantity, ReferenceType, ReferenceID, PreviousQuantity, NewQuantity,
		       Reason, PerformedBy, Notes, CreatedAt
		FROM Transactions
		WHERE 1=1
	`

	args := make([]interface{}, 0)

	// Add filters dynamically
	if txType, ok := filters["transaction_type"].(string); ok && txType != "" {
		query += " AND TransactionType = ?"
		args = append(args, txType)
	}

	if productID, ok := filters["product_id"].(int64); ok && productID > 0 {
		query += " AND ProductID = ?"
		args = append(args, productID)
	}

	if locationID, ok := filters["location_id"].(int64); ok && locationID > 0 {
		query += " AND LocationID = ?"
		args = append(args, locationID)
	}

	query += " ORDER BY TransactionDate DESC, TransactionID DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Errorf("error getting transactions: %v", err)
		return nil, err
	}
	defer rows.Close()

	return scanTransactionRows(rows)
}

// GetByProduct retrieves all transactions for a specific product
func (r *TransactionRepository) GetByProduct(productID int64) ([]*domain.Transaction, error) {
	query := `
		SELECT TransactionID, TransactionType, TransactionDate, ProductID, LocationID,
		       Quantity, ReferenceType, ReferenceID, PreviousQuantity, NewQuantity,
		       Reason, PerformedBy, Notes, CreatedAt
		FROM Transactions
		WHERE ProductID = ?
		ORDER BY TransactionDate DESC, TransactionID DESC
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		log.Errorf("error getting transactions for product %d: %v", productID, err)
		return nil, err
	}
	defer rows.Close()

	return scanTransactionRows(rows)
}

// GetByLocation retrieves all transactions for a specific location
func (r *TransactionRepository) GetByLocation(locationID int64) ([]*domain.Transaction, error) {
	query := `
		SELECT TransactionID, TransactionType, TransactionDate, ProductID, LocationID,
		       Quantity, ReferenceType, ReferenceID, PreviousQuantity, NewQuantity,
		       Reason, PerformedBy, Notes, CreatedAt
		FROM Transactions
		WHERE LocationID = ?
		ORDER BY TransactionDate DESC, TransactionID DESC
	`

	rows, err := r.db.Query(query, locationID)
	if err != nil {
		log.Errorf("error getting transactions for location %d: %v", locationID, err)
		return nil, err
	}
	defer rows.Close()

	return scanTransactionRows(rows)
}

// GetByUser retrieves all transactions performed by a specific user
func (r *TransactionRepository) GetByUser(userID int64) ([]*domain.Transaction, error) {
	query := `
		SELECT TransactionID, TransactionType, TransactionDate, ProductID, LocationID,
		       Quantity, ReferenceType, ReferenceID, PreviousQuantity, NewQuantity,
		       Reason, PerformedBy, Notes, CreatedAt
		FROM Transactions
		WHERE PerformedBy = ?
		ORDER BY TransactionDate DESC, TransactionID DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		log.Errorf("error getting transactions for user %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	return scanTransactionRows(rows)
}

// GetByType retrieves all transactions of a specific type
func (r *TransactionRepository) GetByType(transactionType string) ([]*domain.Transaction, error) {
	query := `
		SELECT TransactionID, TransactionType, TransactionDate, ProductID, LocationID,
		       Quantity, ReferenceType, ReferenceID, PreviousQuantity, NewQuantity,
		       Reason, PerformedBy, Notes, CreatedAt
		FROM Transactions
		WHERE TransactionType = ?
		ORDER BY TransactionDate DESC, TransactionID DESC
	`

	rows, err := r.db.Query(query, transactionType)
	if err != nil {
		log.Errorf("error getting transactions of type %s: %v", transactionType, err)
		return nil, err
	}
	defer rows.Close()

	return scanTransactionRows(rows)
}

// Helper functions

func scanTransactionRow(row *sql.Row) (*domain.Transaction, error) {
	transaction := &domain.Transaction{}
	err := row.Scan(
		&transaction.TransactionID,
		&transaction.TransactionType,
		&transaction.TransactionDate,
		&transaction.ProductID,
		&transaction.LocationID,
		&transaction.Quantity,
		&transaction.ReferenceType,
		&transaction.ReferenceID,
		&transaction.PreviousQuantity,
		&transaction.NewQuantity,
		&transaction.Reason,
		&transaction.PerformedBy,
		&transaction.Notes,
		&transaction.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func scanTransactionRows(rows *sql.Rows) ([]*domain.Transaction, error) {
	transactions := make([]*domain.Transaction, 0)

	for rows.Next() {
		transaction := &domain.Transaction{}
		err := rows.Scan(
			&transaction.TransactionID,
			&transaction.TransactionType,
			&transaction.TransactionDate,
			&transaction.ProductID,
			&transaction.LocationID,
			&transaction.Quantity,
			&transaction.ReferenceType,
			&transaction.ReferenceID,
			&transaction.PreviousQuantity,
			&transaction.NewQuantity,
			&transaction.Reason,
			&transaction.PerformedBy,
			&transaction.Notes,
			&transaction.CreatedAt,
		)
		if err != nil {
			log.Errorf("error scanning transaction row: %v", err)
			continue
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}
