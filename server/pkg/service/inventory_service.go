package service

import (
	"database/sql"
	"fmt"
	"time"

	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service/dto"
)

// IInventoryService defines inventory service interface
type IInventoryService interface {
	// Query operations
	GetAll() ([]*domain.Inventory, error)
	GetByProduct(productID int64) ([]*domain.Inventory, error)
	GetByLocation(locationID int64) ([]*domain.Inventory, error)
	GetByProductAndLocation(productID, locationID int64) (*domain.Inventory, error)

	// Inventory operations
	AdjustInventory(adjustment *dto.InventoryAdjustment) error
	TransferInventory(transfer *dto.InventoryTransfer) error

	// Stock validation
	ValidateSufficientStock(productID, locationID, requiredQuantity int64) error
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

// GetByProduct retrieves inventory for a specific product across all locations
func (s *InventoryService) GetByProduct(productID int64) ([]*domain.Inventory, error) {
	if productID <= 0 {
		return nil, fmt.Errorf("invalid product ID")
	}
	return s.inventoryRepo.GetByProduct(productID)
}

// GetByLocation retrieves inventory at a specific location
func (s *InventoryService) GetByLocation(locationID int64) ([]*domain.Inventory, error) {
	if locationID <= 0 {
		return nil, fmt.Errorf("invalid location ID")
	}
	return s.inventoryRepo.GetByLocation(locationID)
}

// GetByProductAndLocation retrieves inventory for a specific product at a specific location
func (s *InventoryService) GetByProductAndLocation(productID, locationID int64) (*domain.Inventory, error) {
	if productID <= 0 || locationID <= 0 {
		return nil, fmt.Errorf("invalid product or location ID")
	}
	return s.inventoryRepo.GetByProductAndLocation(productID, locationID)
}

// AdjustInventory performs manual inventory adjustment
func (s *InventoryService) AdjustInventory(adjustment *dto.InventoryAdjustment) error {
	// Validate input
	if err := s.validateAdjustment(adjustment); err != nil {
		return err
	}

	// Get current inventory
	current, err := s.inventoryRepo.GetByProductAndLocation(adjustment.ProductID, adjustment.LocationID)
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
		err = s.inventoryRepo.CreateOrUpdate(adjustment.ProductID, adjustment.LocationID, adjustment.Quantity)
	} else {
		// Update existing record
		err = s.inventoryRepo.UpdateQuantity(adjustment.ProductID, adjustment.LocationID, adjustment.Quantity)
	}

	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	// Create transaction record for audit trail
	transactionRecord := &domain.Transaction{
		TransactionType:  adjustment.Reason, // e.g., "adjustment", "damage", "theft"
		TransactionDate:  time.Now(),
		ProductID:        adjustment.ProductID,
		LocationID:       &adjustment.LocationID,
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

// TransferInventory transfers inventory between locations
func (s *InventoryService) TransferInventory(transfer *dto.InventoryTransfer) error {
	// Validate input
	if err := s.validateTransfer(transfer); err != nil {
		return err
	}

	// Validate sufficient stock at source
	hasSufficient, err := s.inventoryRepo.HasSufficientStock(
		transfer.ProductID,
		transfer.FromLocationID,
		transfer.Quantity,
	)
	if err != nil {
		return fmt.Errorf("failed to check stock: %w", err)
	}

	if !hasSufficient {
		sourceInv, _ := s.inventoryRepo.GetByProductAndLocation(transfer.ProductID, transfer.FromLocationID)
		available := int64(0)
		if sourceInv != nil {
			available = sourceInv.AvailableQuantity
		}
		return fmt.Errorf("insufficient inventory at source location: available=%d, required=%d",
			available, transfer.Quantity)
	}

	// Get current quantities for audit trail
	sourceInv, err := s.inventoryRepo.GetByProductAndLocation(transfer.ProductID, transfer.FromLocationID)
	if err != nil {
		return fmt.Errorf("failed to get source inventory: %w", err)
	}

	destInv, _ := s.inventoryRepo.GetByProductAndLocation(transfer.ProductID, transfer.ToLocationID)
	destPrevQty := int64(0)
	if destInv != nil {
		destPrevQty = destInv.Quantity
	}

	sourcePrevQty := sourceInv.Quantity
	sourceNewQty := sourcePrevQty - transfer.Quantity
	destNewQty := destPrevQty + transfer.Quantity

	// Start database transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Reduce quantity at source location
	err = s.inventoryRepo.UpdateQuantity(transfer.ProductID, transfer.FromLocationID, -transfer.Quantity)
	if err != nil {
		return fmt.Errorf("failed to reduce source inventory: %w", err)
	}

	// Increase quantity at destination location
	err = s.inventoryRepo.UpdateQuantity(transfer.ProductID, transfer.ToLocationID, transfer.Quantity)
	if err != nil {
		return fmt.Errorf("failed to increase destination inventory: %w", err)
	}

	// Create transaction record for source (outbound)
	sourceTransaction := &domain.Transaction{
		TransactionType:  "transfer",
		TransactionDate:  time.Now(),
		ProductID:        transfer.ProductID,
		LocationID:       &transfer.FromLocationID,
		Quantity:         -transfer.Quantity, // Negative for outbound
		ReferenceType:    transfer.Reason,
		PreviousQuantity: &sourcePrevQty,
		NewQuantity:      &sourceNewQty,
		Reason:           transfer.Reason,
		PerformedBy:      &transfer.PerformedBy,
		Notes:            stringPtr(fmt.Sprintf("Transfer to location %d", transfer.ToLocationID)),
	}

	err = s.transactionRepo.Create(sourceTransaction)
	if err != nil {
		return fmt.Errorf("failed to create source transaction record: %w", err)
	}

	// Create transaction record for destination (inbound)
	destTransaction := &domain.Transaction{
		TransactionType:  "transfer",
		TransactionDate:  time.Now(),
		ProductID:        transfer.ProductID,
		LocationID:       &transfer.ToLocationID,
		Quantity:         transfer.Quantity, // Positive for inbound
		ReferenceType:    transfer.Reason,
		PreviousQuantity: &destPrevQty,
		NewQuantity:      &destNewQty,
		Reason:           transfer.Reason,
		PerformedBy:      &transfer.PerformedBy,
		Notes:            stringPtr(fmt.Sprintf("Transfer from location %d", transfer.FromLocationID)),
	}

	err = s.transactionRepo.Create(destTransaction)
	if err != nil {
		return fmt.Errorf("failed to create destination transaction record: %w", err)
	}

	// Commit database transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ValidateSufficientStock validates if there's sufficient stock for an operation
func (s *InventoryService) ValidateSufficientStock(productID, locationID, requiredQuantity int64) error {
	if requiredQuantity <= 0 {
		return fmt.Errorf("required quantity must be positive")
	}

	hasSufficient, err := s.inventoryRepo.HasSufficientStock(productID, locationID, requiredQuantity)
	if err != nil {
		return fmt.Errorf("failed to check stock: %w", err)
	}

	if !hasSufficient {
		inventory, _ := s.inventoryRepo.GetByProductAndLocation(productID, locationID)
		available := int64(0)
		if inventory != nil {
			available = inventory.AvailableQuantity
		}
		return fmt.Errorf("insufficient inventory: product_id=%d, location_id=%d, available=%d, required=%d",
			productID, locationID, available, requiredQuantity)
	}

	return nil
}

// Validation helpers

func (s *InventoryService) validateAdjustment(adjustment *dto.InventoryAdjustment) error {
	if adjustment.ProductID <= 0 {
		return fmt.Errorf("invalid product ID")
	}
	if adjustment.LocationID <= 0 {
		return fmt.Errorf("invalid location ID")
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

func (s *InventoryService) validateTransfer(transfer *dto.InventoryTransfer) error {
	if transfer.ProductID <= 0 {
		return fmt.Errorf("invalid product ID")
	}
	if transfer.FromLocationID <= 0 {
		return fmt.Errorf("invalid source location ID")
	}
	if transfer.ToLocationID <= 0 {
		return fmt.Errorf("invalid destination location ID")
	}
	if transfer.FromLocationID == transfer.ToLocationID {
		return fmt.Errorf("source and destination locations cannot be the same")
	}
	if transfer.Quantity <= 0 {
		return fmt.Errorf("transfer quantity must be positive")
	}
	if transfer.PerformedBy <= 0 {
		return fmt.Errorf("performed_by user ID is required")
	}

	return nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
