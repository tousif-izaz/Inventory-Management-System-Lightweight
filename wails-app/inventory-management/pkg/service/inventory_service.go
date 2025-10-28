package service

import (
	"database/sql"
	"fmt"
	"time"

	"inventory-management/pkg/domain"
	"inventory-management/pkg/repository"
	"inventory-management/pkg/service/dto"
)

// IInventoryService defines inventory service interface
type IInventoryService interface {
	// Query operations
	GetAll() ([]*domain.Inventory, error)
	GetByProduct(productID int64) (*domain.Inventory, error)

	// Inventory operations
	AdjustInventory(adjustment *dto.InventoryAdjustment) error

	// Stock validation
	ValidateSufficientStock(productID, requiredQuantity int64) error
}

type InventoryService struct {
	inventoryRepo   repository.IInventoryRepository
	transactionRepo repository.ITransactionRepository
	db              *sql.DB
}

// NewInventoryService creates a new inventory service
func NewInventoryService(
	inventoryRepo repository.IInventoryRepository,
	transactionRepo repository.ITransactionRepository,
	db *sql.DB,
) IInventoryService {
	return &InventoryService{
		inventoryRepo:   inventoryRepo,
		transactionRepo: transactionRepo,
		db:              db,
	}
}

// GetAll retrieves all inventory records
func (s *InventoryService) GetAll() ([]*domain.Inventory, error) {
	return s.inventoryRepo.GetAll()
}

// GetByProduct retrieves inventory for a specific product
func (s *InventoryService) GetByProduct(productID int64) (*domain.Inventory, error) {
	if productID <= 0 {
		return nil, fmt.Errorf("invalid product ID")
	}
	return s.inventoryRepo.GetByProduct(productID)
}

// AdjustInventory performs manual inventory adjustment
func (s *InventoryService) AdjustInventory(adjustment *dto.InventoryAdjustment) error {
	// Validate input
	if err := s.validateAdjustment(adjustment); err != nil {
		return err
	}

	// Get current inventory
	current, err := s.inventoryRepo.GetByProduct(adjustment.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get current inventory: %w", err)
	}

	var previousQty int64
	if current != nil {
		previousQty = current.Quantity
	} else {
		previousQty = 0
	}

	newQty := previousQty + adjustment.Quantity

	// Validate non-negative result
	if newQty < 0 {
		return fmt.Errorf("adjustment would result in negative inventory: current=%d, adjustment=%d",
			previousQty, adjustment.Quantity)
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update inventory
	if current == nil {
		// Create new inventory record
		err = s.inventoryRepo.CreateOrUpdate(adjustment.ProductID, adjustment.Quantity)
	} else {
		// Update existing record
		err = s.inventoryRepo.UpdateQuantity(adjustment.ProductID, adjustment.Quantity)
	}

	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	// Create transaction record for audit trail
	transactionRecord := &domain.Transaction{
		TransactionType:  adjustment.Reason, // e.g., "adjustment", "damage", "theft"
		TransactionDate:  time.Now(),
		ProductID:        adjustment.ProductID,
		Quantity:         adjustment.Quantity,
		PreviousQuantity: &previousQty,
		NewQuantity:      &newQty,
		Reason:           &adjustment.Reason,
		PerformedBy:      &adjustment.PerformedBy,
		Notes:            adjustment.Notes,
	}

	err = s.transactionRepo.Create(transactionRecord)
	if err != nil {
		return fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ValidateSufficientStock validates if there's sufficient stock for an operation
func (s *InventoryService) ValidateSufficientStock(productID, requiredQuantity int64) error {
	if requiredQuantity <= 0 {
		return fmt.Errorf("required quantity must be positive")
	}

	hasSufficient, err := s.inventoryRepo.HasSufficientStock(productID, requiredQuantity)
	if err != nil {
		return fmt.Errorf("failed to check stock: %w", err)
	}

	if !hasSufficient {
		inventory, _ := s.inventoryRepo.GetByProduct(productID)
		available := int64(0)
		if inventory != nil {
			available = inventory.AvailableQuantity
		}
		return fmt.Errorf("insufficient inventory: product_id=%d, available=%d, required=%d",
			productID, available, requiredQuantity)
	}

	return nil
}

// Validation helpers

func (s *InventoryService) validateAdjustment(adjustment *dto.InventoryAdjustment) error {
	if adjustment.ProductID <= 0 {
		return fmt.Errorf("invalid product ID")
	}
	if adjustment.Quantity == 0 {
		return fmt.Errorf("adjustment quantity cannot be zero")
	}
	if adjustment.Reason == "" {
		return fmt.Errorf("reason is required for inventory adjustment")
	}
	if adjustment.PerformedBy <= 0 {
		return fmt.Errorf("performed_by user ID is required")
	}

	// Validate reason is one of the allowed values
	validReasons := map[string]bool{
		"adjustment": true,
		"damage":     true,
		"theft":      true,
		"return":     true,
		"other":      true,
	}

	if !validReasons[adjustment.Reason] {
		return fmt.Errorf("invalid reason: must be one of adjustment, damage, theft, return, other")
	}

	return nil
}
