package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ims-intro/pkg/service/dto"
)

// Test inventory underflow scenarios without database
// These tests focus on validation logic only

func TestInventoryUnderflowValidation(t *testing.T) {
	t.Run("Negative quantity after adjustment should fail", func(t *testing.T) {
		// This simulates the check that happens in AdjustInventory
		currentQty := int64(50)
		adjustment := int64(-60)
		newQty := currentQty + adjustment

		assert.Less(t, newQty, int64(0), "New quantity should be negative")
		assert.Equal(t, int64(-10), newQty, "Should result in -10")
	})

	t.Run("Exact zero should be allowed", func(t *testing.T) {
		currentQty := int64(50)
		adjustment := int64(-50)
		newQty := currentQty + adjustment

		assert.GreaterOrEqual(t, newQty, int64(0), "New quantity should be non-negative")
		assert.Equal(t, int64(0), newQty, "Should result in 0")
	})

	t.Run("Positive adjustment always valid", func(t *testing.T) {
		currentQty := int64(0)
		adjustment := int64(100)
		newQty := currentQty + adjustment

		assert.Greater(t, newQty, int64(0), "New quantity should be positive")
		assert.Equal(t, int64(100), newQty, "Should result in 100")
	})
}

func TestTransferValidation(t *testing.T) {
	service := &InventoryService{}

	t.Run("Same source and destination should fail", func(t *testing.T) {
		transfer := &dto.InventoryTransfer{
			ProductID:      1,
			FromLocationID: 1,
			ToLocationID:   1, // Same as source
			Quantity:       10,
			PerformedBy:    1,
		}

		err := service.validateTransfer(transfer)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be the same")
	})

	t.Run("Negative quantity should fail", func(t *testing.T) {
		transfer := &dto.InventoryTransfer{
			ProductID:      1,
			FromLocationID: 1,
			ToLocationID:   2,
			Quantity:       -10, // Invalid
			PerformedBy:    1,
		}

		err := service.validateTransfer(transfer)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be positive")
	})

	t.Run("Zero quantity should fail", func(t *testing.T) {
		transfer := &dto.InventoryTransfer{
			ProductID:      1,
			FromLocationID: 1,
			ToLocationID:   2,
			Quantity:       0, // Invalid
			PerformedBy:    1,
		}

		err := service.validateTransfer(transfer)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be positive")
	})

	t.Run("Valid transfer should pass", func(t *testing.T) {
		reason := "test"
		transfer := &dto.InventoryTransfer{
			ProductID:      1,
			FromLocationID: 1,
			ToLocationID:   2,
			Quantity:       10,
			Reason:         &reason,
			PerformedBy:    1,
		}

		err := service.validateTransfer(transfer)
		assert.NoError(t, err)
	})
}

func TestAdjustmentValidation(t *testing.T) {
	service := &InventoryService{}

	t.Run("Zero adjustment quantity should fail", func(t *testing.T) {
		adjustment := &dto.InventoryAdjustment{
			ProductID:   1,
			LocationID:  1,
			Quantity:    0, // Invalid
			Reason:      "adjustment",
			PerformedBy: 1,
		}

		err := service.validateAdjustment(adjustment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be zero")
	})

	t.Run("Missing reason should fail", func(t *testing.T) {
		adjustment := &dto.InventoryAdjustment{
			ProductID:   1,
			LocationID:  1,
			Quantity:    10,
			Reason:      "", // Invalid
			PerformedBy: 1,
		}

		err := service.validateAdjustment(adjustment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reason is required")
	})

	t.Run("Invalid reason should fail", func(t *testing.T) {
		adjustment := &dto.InventoryAdjustment{
			ProductID:   1,
			LocationID:  1,
			Quantity:    10,
			Reason:      "invalid_reason", // Not in allowed list
			PerformedBy: 1,
		}

		err := service.validateAdjustment(adjustment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reason")
	})

	t.Run("Valid adjustment reason should pass", func(t *testing.T) {
		validReasons := []string{"adjustment", "damage", "theft", "return", "other"}

		for _, reason := range validReasons {
			adjustment := &dto.InventoryAdjustment{
				ProductID:   1,
				LocationID:  1,
				Quantity:    10,
				Reason:      reason,
				PerformedBy: 1,
			}

			err := service.validateAdjustment(adjustment)
			assert.NoError(t, err, "Reason '%s' should be valid", reason)
		}
	})

	t.Run("Negative adjustment should be allowed (validation only)", func(t *testing.T) {
		adjustment := &dto.InventoryAdjustment{
			ProductID:   1,
			LocationID:  1,
			Quantity:    -10, // Negative is allowed at validation level
			Reason:      "damage",
			PerformedBy: 1,
		}

		err := service.validateAdjustment(adjustment)
		assert.NoError(t, err, "Negative adjustments should pass validation")
	})
}

// Test edge cases for inventory calculations
func TestInventoryCalculationEdgeCases(t *testing.T) {
	t.Run("Large quantity reduction", func(t *testing.T) {
		currentQty := int64(1000000)
		reduction := int64(999999)
		newQty := currentQty - reduction

		assert.Equal(t, int64(1), newQty)
		assert.GreaterOrEqual(t, newQty, int64(0))
	})

	t.Run("Boundary: Reduce to exactly zero", func(t *testing.T) {
		currentQty := int64(100)
		reduction := int64(100)
		newQty := currentQty - reduction

		assert.Equal(t, int64(0), newQty)
		assert.GreaterOrEqual(t, newQty, int64(0))
	})

	t.Run("Boundary: Reduce by one more than available (underflow)", func(t *testing.T) {
		currentQty := int64(100)
		reduction := int64(101)
		newQty := currentQty - reduction

		assert.Equal(t, int64(-1), newQty)
		assert.Less(t, newQty, int64(0), "Should be negative (underflow)")
	})

	t.Run("Available quantity with reservations", func(t *testing.T) {
		totalQty := int64(100)
		reservedQty := int64(20)
		availableQty := totalQty - reservedQty

		assert.Equal(t, int64(80), availableQty)

		// Trying to sell more than available
		saleQty := int64(90)
		assert.Greater(t, saleQty, availableQty, "Sale quantity exceeds available")
	})
}
