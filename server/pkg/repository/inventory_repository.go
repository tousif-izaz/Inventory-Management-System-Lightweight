package repository

import (
	"database/sql"
	"fmt"

	"github.com/labstack/gommon/log"
	"ims-intro/pkg/domain"
)

// IInventoryRepository defines inventory repository interface
type IInventoryRepository interface {
	// Basic CRUD
	GetAll() ([]*domain.Inventory, error)
	GetByID(inventoryID int64) (*domain.Inventory, error)
	GetByProductAndLocation(productID, locationID int64) (*domain.Inventory, error)
	GetByProduct(productID int64) ([]*domain.Inventory, error)
	GetByLocation(locationID int64) ([]*domain.Inventory, error)

	// Inventory operations
	CreateOrUpdate(productID, locationID, quantity int64) error
	UpdateQuantity(productID, locationID int64, quantityChange int64) error
	SetQuantity(productID, locationID, newQuantity int64) error

	// Stock checking
	GetAvailableQuantity(productID, locationID int64) (int64, error)
	HasSufficientStock(productID, locationID, requiredQuantity int64) (bool, error)
}

type InventoryRepository struct {
	db *sql.DB
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *sql.DB) IInventoryRepository {
	return &InventoryRepository{db}
}

// GetAll retrieves all inventory records
func (r *InventoryRepository) GetAll() ([]*domain.Inventory, error) {
	query := `
		SELECT InventoryID, ProductID, LocationID, Quantity, ReservedQuantity,
		       (Quantity - ReservedQuantity) as AvailableQuantity, LastUpdated
		FROM Inventory
		ORDER BY InventoryID
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Errorf("error getting all inventory: %v", err)
		return nil, err
	}
	defer rows.Close()

	return scanInventoryRows(rows)
}

// GetByID retrieves inventory by ID
func (r *InventoryRepository) GetByID(inventoryID int64) (*domain.Inventory, error) {
	query := `
		SELECT InventoryID, ProductID, LocationID, Quantity, ReservedQuantity,
		       (Quantity - ReservedQuantity) as AvailableQuantity, LastUpdated
		FROM Inventory
		WHERE InventoryID = ?
	`

	row := r.db.QueryRow(query, inventoryID)
	return scanInventoryRow(row)
}

// GetByProductAndLocation retrieves inventory for a specific product at a specific location
func (r *InventoryRepository) GetByProductAndLocation(productID, locationID int64) (*domain.Inventory, error) {
	query := `
		SELECT InventoryID, ProductID, LocationID, Quantity, ReservedQuantity,
		       (Quantity - ReservedQuantity) as AvailableQuantity, LastUpdated
		FROM Inventory
		WHERE ProductID = ? AND LocationID = ?
	`

	row := r.db.QueryRow(query, productID, locationID)
	inventory, err := scanInventoryRow(row)
	if err == sql.ErrNoRows {
		return nil, nil // Not found, return nil without error
	}
	return inventory, err
}

// GetByProduct retrieves all inventory records for a product across all locations
func (r *InventoryRepository) GetByProduct(productID int64) ([]*domain.Inventory, error) {
	query := `
		SELECT InventoryID, ProductID, LocationID, Quantity, ReservedQuantity,
		       (Quantity - ReservedQuantity) as AvailableQuantity, LastUpdated
		FROM Inventory
		WHERE ProductID = ?
		ORDER BY LocationID
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		log.Errorf("error getting inventory for product %d: %v", productID, err)
		return nil, err
	}
	defer rows.Close()

	return scanInventoryRows(rows)
}

// GetByLocation retrieves all inventory records at a specific location
func (r *InventoryRepository) GetByLocation(locationID int64) ([]*domain.Inventory, error) {
	query := `
		SELECT InventoryID, ProductID, LocationID, Quantity, ReservedQuantity,
		       (Quantity - ReservedQuantity) as AvailableQuantity, LastUpdated
		FROM Inventory
		WHERE LocationID = ?
		ORDER BY ProductID
	`

	rows, err := r.db.Query(query, locationID)
	if err != nil {
		log.Errorf("error getting inventory for location %d: %v", locationID, err)
		return nil, err
	}
	defer rows.Close()

	return scanInventoryRows(rows)
}

// CreateOrUpdate creates or updates inventory record (upsert)
func (r *InventoryRepository) CreateOrUpdate(productID, locationID, quantity int64) error {
	query := `
		INSERT INTO Inventory (ProductID, LocationID, Quantity, ReservedQuantity)
		VALUES (?, ?, ?, 0)
		ON CONFLICT(ProductID, LocationID)
		DO UPDATE SET Quantity = Quantity + excluded.Quantity
	`

	_, err := r.db.Exec(query, productID, locationID, quantity)
	if err != nil {
		log.Errorf("error creating/updating inventory: %v", err)
		return err
	}

	log.Infof("Inventory created/updated: Product=%d, Location=%d, Quantity=%d", productID, locationID, quantity)
	return nil
}

// UpdateQuantity updates inventory by adding/subtracting from current quantity
func (r *InventoryRepository) UpdateQuantity(productID, locationID int64, quantityChange int64) error {
	// First check if record exists
	existing, err := r.GetByProductAndLocation(productID, locationID)
	if err != nil {
		return err
	}

	if existing == nil {
		// Create new record if adding stock
		if quantityChange > 0 {
			return r.CreateOrUpdate(productID, locationID, quantityChange)
		}
		return fmt.Errorf("cannot reduce inventory: product %d not found at location %d", productID, locationID)
	}

	newQuantity := existing.Quantity + quantityChange

	// Validate non-negative quantity
	if newQuantity < 0 {
		return fmt.Errorf("insufficient inventory: product %d at location %d has %d units, cannot reduce by %d",
			productID, locationID, existing.Quantity, -quantityChange)
	}

	query := `
		UPDATE Inventory
		SET Quantity = Quantity + ?
		WHERE ProductID = ? AND LocationID = ?
	`

	_, err = r.db.Exec(query, quantityChange, productID, locationID)
	if err != nil {
		log.Errorf("error updating inventory quantity: %v", err)
		return err
	}

	log.Infof("Inventory updated: Product=%d, Location=%d, Change=%d, NewQuantity=%d",
		productID, locationID, quantityChange, newQuantity)
	return nil
}

// SetQuantity sets inventory to a specific absolute value
func (r *InventoryRepository) SetQuantity(productID, locationID, newQuantity int64) error {
	if newQuantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	// First check if record exists
	existing, err := r.GetByProductAndLocation(productID, locationID)
	if err != nil {
		return err
	}

	if existing == nil {
		// Create new record
		return r.CreateOrUpdate(productID, locationID, newQuantity)
	}

	query := `
		UPDATE Inventory
		SET Quantity = ?
		WHERE ProductID = ? AND LocationID = ?
	`

	_, err = r.db.Exec(query, newQuantity, productID, locationID)
	if err != nil {
		log.Errorf("error setting inventory quantity: %v", err)
		return err
	}

	log.Infof("Inventory set: Product=%d, Location=%d, NewQuantity=%d", productID, locationID, newQuantity)
	return nil
}

// GetAvailableQuantity returns the available (non-reserved) quantity
func (r *InventoryRepository) GetAvailableQuantity(productID, locationID int64) (int64, error) {
	inventory, err := r.GetByProductAndLocation(productID, locationID)
	if err != nil {
		return 0, err
	}
	if inventory == nil {
		return 0, nil
	}
	return inventory.AvailableQuantity, nil
}

// HasSufficientStock checks if there's enough available stock
func (r *InventoryRepository) HasSufficientStock(productID, locationID, requiredQuantity int64) (bool, error) {
	available, err := r.GetAvailableQuantity(productID, locationID)
	if err != nil {
		return false, err
	}
	return available >= requiredQuantity, nil
}

// Helper functions

func scanInventoryRow(row *sql.Row) (*domain.Inventory, error) {
	inventory := &domain.Inventory{}
	err := row.Scan(
		&inventory.InventoryID,
		&inventory.ProductID,
		&inventory.LocationID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.AvailableQuantity,
		&inventory.LastUpdated,
	)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

func scanInventoryRows(rows *sql.Rows) ([]*domain.Inventory, error) {
	inventories := make([]*domain.Inventory, 0)

	for rows.Next() {
		inventory := &domain.Inventory{}
		err := rows.Scan(
			&inventory.InventoryID,
			&inventory.ProductID,
			&inventory.LocationID,
			&inventory.Quantity,
			&inventory.ReservedQuantity,
			&inventory.AvailableQuantity,
			&inventory.LastUpdated,
		)
		if err != nil {
			log.Errorf("error scanning inventory row: %v", err)
			continue
		}
		inventories = append(inventories, inventory)
	}

	return inventories, nil
}
