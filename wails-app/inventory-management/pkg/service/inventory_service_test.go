package service

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"inventory-management/pkg/domain"
	"inventory-management/pkg/service/dto"
)

// Mock repositories
type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) GetAll() ([]*domain.Inventory, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) GetByID(inventoryID int64) (*domain.Inventory, error) {
	args := m.Called(inventoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) GetByProductAndLocation(productID, locationID int64) (*domain.Inventory, error) {
	args := m.Called(productID, locationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) GetByProduct(productID int64) ([]*domain.Inventory, error) {
	args := m.Called(productID)
	return args.Get(0).([]*domain.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) GetByLocation(locationID int64) ([]*domain.Inventory, error) {
	args := m.Called(locationID)
	return args.Get(0).([]*domain.Inventory), args.Error(1)
}

func (m *MockInventoryRepository) CreateOrUpdate(productID, locationID, quantity int64) error {
	args := m.Called(productID, locationID, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) UpdateQuantity(productID, locationID int64, quantityChange int64) error {
	args := m.Called(productID, locationID, quantityChange)
	return args.Error(0)
}

func (m *MockInventoryRepository) SetQuantity(productID, locationID, newQuantity int64) error {
	args := m.Called(productID, locationID, newQuantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) GetAvailableQuantity(productID, locationID int64) (int64, error) {
	args := m.Called(productID, locationID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockInventoryRepository) HasSufficientStock(productID, locationID, requiredQuantity int64) (bool, error) {
	args := m.Called(productID, locationID, requiredQuantity)
	return args.Bool(0), args.Error(1)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(transaction *domain.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(transactionID int64) (*domain.Transaction, error) {
	args := m.Called(transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetAll(filters map[string]interface{}) ([]*domain.Transaction, error) {
	args := m.Called(filters)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByProduct(productID int64) ([]*domain.Transaction, error) {
	args := m.Called(productID)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByLocation(locationID int64) ([]*domain.Transaction, error) {
	args := m.Called(locationID)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByUser(userID int64) ([]*domain.Transaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByType(transactionType string) ([]*domain.Transaction, error) {
	args := m.Called(transactionType)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin() (*sql.Tx, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Tx), args.Error(1)
}

// Test Suite
type InventoryServiceTestSuite struct {
	suite.Suite
	inventoryRepo   *MockInventoryRepository
	transactionRepo *MockTransactionRepository
	service         IInventoryService
	mockDB          *sql.DB
}

func (suite *InventoryServiceTestSuite) SetupTest() {
	suite.inventoryRepo = new(MockInventoryRepository)
	suite.transactionRepo = new(MockTransactionRepository)
	// Use nil for db as we're mocking the repositories
	suite.service = NewInventoryService(suite.inventoryRepo, suite.transactionRepo, nil)
}

func TestInventoryServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryServiceTestSuite))
}

// Test: Successful inventory adjustment (increase)
func (suite *InventoryServiceTestSuite) TestAdjustInventory_Increase_Success() {
	adjustment := &dto.InventoryAdjustment{
		ProductID:   1,
		LocationID:  1,
		Quantity:    10,
		Reason:      "adjustment",
		Notes:       stringPtr("Adding stock"),
		PerformedBy: 1,
	}

	currentInventory := &domain.Inventory{
		InventoryID:       1,
		ProductID:         1,
		LocationID:        1,
		Quantity:          50,
		ReservedQuantity:  0,
		AvailableQuantity: 50,
		LastUpdated:       time.Now(),
	}

	suite.inventoryRepo.On("GetByProductAndLocation", int64(1), int64(1)).Return(currentInventory, nil)
	suite.inventoryRepo.On("UpdateQuantity", int64(1), int64(1), int64(10)).Return(nil)
	suite.transactionRepo.On("Create", mock.AnythingOfType("*domain.Transaction")).Return(nil)

	err := suite.service.AdjustInventory(adjustment)

	assert.NoError(suite.T(), err)
	suite.inventoryRepo.AssertExpectations(suite.T())
	suite.transactionRepo.AssertExpectations(suite.T())
}

// Test: Inventory underflow - attempting to reduce below zero
func (suite *InventoryServiceTestSuite) TestAdjustInventory_Underflow_Error() {
	adjustment := &dto.InventoryAdjustment{
		ProductID:   1,
		LocationID:  1,
		Quantity:    -60, // Trying to reduce by 60
		Reason:      "damage",
		Notes:       stringPtr("Damaged goods"),
		PerformedBy: 1,
	}

	currentInventory := &domain.Inventory{
		InventoryID:       1,
		ProductID:         1,
		LocationID:        1,
		Quantity:          50, // Only 50 available
		ReservedQuantity:  0,
		AvailableQuantity: 50,
		LastUpdated:       time.Now(),
	}

	suite.inventoryRepo.On("GetByProductAndLocation", int64(1), int64(1)).Return(currentInventory, nil)

	err := suite.service.AdjustInventory(adjustment)

	// Should return error about negative inventory
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "negative inventory")
	suite.inventoryRepo.AssertExpectations(suite.T())
	// Transaction repo should NOT be called
	suite.transactionRepo.AssertNotCalled(suite.T(), "Create")
}

// Test: Inventory underflow - exact boundary (reducing to exactly zero should work)
func (suite *InventoryServiceTestSuite) TestAdjustInventory_ReduceToZero_Success() {
	adjustment := &dto.InventoryAdjustment{
		ProductID:   1,
		LocationID:  1,
		Quantity:    -50, // Reduce to exactly zero
		Reason:      "damage",
		Notes:       stringPtr("All damaged"),
		PerformedBy: 1,
	}

	currentInventory := &domain.Inventory{
		InventoryID:       1,
		ProductID:         1,
		LocationID:        1,
		Quantity:          50,
		ReservedQuantity:  0,
		AvailableQuantity: 50,
		LastUpdated:       time.Now(),
	}

	suite.inventoryRepo.On("GetByProductAndLocation", int64(1), int64(1)).Return(currentInventory, nil)
	suite.inventoryRepo.On("UpdateQuantity", int64(1), int64(1), int64(-50)).Return(nil)
	suite.transactionRepo.On("Create", mock.AnythingOfType("*domain.Transaction")).Return(nil)

	err := suite.service.AdjustInventory(adjustment)

	assert.NoError(suite.T(), err)
	suite.inventoryRepo.AssertExpectations(suite.T())
	suite.transactionRepo.AssertExpectations(suite.T())
}

// Test: Transfer with insufficient stock
func (suite *InventoryServiceTestSuite) TestTransferInventory_InsufficientStock_Error() {
	transfer := &dto.InventoryTransfer{
		ProductID:      1,
		FromLocationID: 1,
		ToLocationID:   2,
		Quantity:       100, // Trying to transfer 100
		Reason:         stringPtr("Restocking"),
		PerformedBy:    1,
	}

	sourceInventory := &domain.Inventory{
		InventoryID:       1,
		ProductID:         1,
		LocationID:        1,
		Quantity:          50, // Only 50 available
		ReservedQuantity:  0,
		AvailableQuantity: 50,
		LastUpdated:       time.Now(),
	}

	suite.inventoryRepo.On("HasSufficientStock", int64(1), int64(1), int64(100)).Return(false, nil)
	suite.inventoryRepo.On("GetByProductAndLocation", int64(1), int64(1)).Return(sourceInventory, nil)

	err := suite.service.TransferInventory(transfer)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "insufficient inventory")
	assert.Contains(suite.T(), err.Error(), "available=50")
	assert.Contains(suite.T(), err.Error(), "required=100")
	suite.inventoryRepo.AssertExpectations(suite.T())
}

// Test: Successful inventory transfer
func (suite *InventoryServiceTestSuite) TestTransferInventory_Success() {
	transfer := &dto.InventoryTransfer{
		ProductID:      1,
		FromLocationID: 1,
		ToLocationID:   2,
		Quantity:       30,
		Reason:         stringPtr("Restocking"),
		PerformedBy:    1,
	}

	sourceInventory := &domain.Inventory{
		InventoryID:       1,
		ProductID:         1,
		LocationID:        1,
		Quantity:          50,
		ReservedQuantity:  0,
		AvailableQuantity: 50,
		LastUpdated:       time.Now(),
	}

	destInventory := &domain.Inventory{
		InventoryID:       2,
		ProductID:         1,
		LocationID:        2,
		Quantity:          20,
		ReservedQuantity:  0,
		AvailableQuantity: 20,
		LastUpdated:       time.Now(),
	}

	suite.inventoryRepo.On("HasSufficientStock", int64(1), int64(1), int64(30)).Return(true, nil)
	suite.inventoryRepo.On("GetByProductAndLocation", int64(1), int64(1)).Return(sourceInventory, nil)
	suite.inventoryRepo.On("GetByProductAndLocation", int64(1), int64(2)).Return(destInventory, nil)
	suite.inventoryRepo.On("UpdateQuantity", int64(1), int64(1), int64(-30)).Return(nil)
	suite.inventoryRepo.On("UpdateQuantity", int64(1), int64(2), int64(30)).Return(nil)
	suite.transactionRepo.On("Create", mock.AnythingOfType("*domain.Transaction")).Return(nil).Times(2)

	err := suite.service.TransferInventory(transfer)

	assert.NoError(suite.T(), err)
	suite.inventoryRepo.AssertExpectations(suite.T())
	suite.transactionRepo.AssertExpectations(suite.T())
}

// Test: Transfer to same location (should fail validation)
func (suite *InventoryServiceTestSuite) TestTransferInventory_SameLocation_Error() {
	transfer := &dto.InventoryTransfer{
		ProductID:      1,
		FromLocationID: 1,
		ToLocationID:   1, // Same as source
		Quantity:       10,
		PerformedBy:    1,
	}

	err := suite.service.TransferInventory(transfer)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "cannot be the same")
}

// Test: Validate sufficient stock - success
func (suite *InventoryServiceTestSuite) TestValidateSufficientStock_Success() {
	suite.inventoryRepo.On("HasSufficientStock", int64(1), int64(1), int64(10)).Return(true, nil)

	err := suite.service.ValidateSufficientStock(1, 1, 10)

	assert.NoError(suite.T(), err)
	suite.inventoryRepo.AssertExpectations(suite.T())
}

// Test: Validate sufficient stock - insufficient
func (suite *InventoryServiceTestSuite) TestValidateSufficientStock_Insufficient() {
	inventory := &domain.Inventory{
		InventoryID:       1,
		ProductID:         1,
		LocationID:        1,
		Quantity:          5,
		ReservedQuantity:  0,
		AvailableQuantity: 5,
		LastUpdated:       time.Now(),
	}

	suite.inventoryRepo.On("HasSufficientStock", int64(1), int64(1), int64(10)).Return(false, nil)
	suite.inventoryRepo.On("GetByProductAndLocation", int64(1), int64(1)).Return(inventory, nil)

	err := suite.service.ValidateSufficientStock(1, 1, 10)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "insufficient inventory")
	assert.Contains(suite.T(), err.Error(), "available=5")
	assert.Contains(suite.T(), err.Error(), "required=10")
	suite.inventoryRepo.AssertExpectations(suite.T())
}

// Test: Invalid adjustment - zero quantity
func (suite *InventoryServiceTestSuite) TestAdjustInventory_ZeroQuantity_Error() {
	adjustment := &dto.InventoryAdjustment{
		ProductID:   1,
		LocationID:  1,
		Quantity:    0, // Invalid
		Reason:      "adjustment",
		PerformedBy: 1,
	}

	err := suite.service.AdjustInventory(adjustment)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "cannot be zero")
}

// Test: Invalid adjustment - no reason
func (suite *InventoryServiceTestSuite) TestAdjustInventory_NoReason_Error() {
	adjustment := &dto.InventoryAdjustment{
		ProductID:   1,
		LocationID:  1,
		Quantity:    10,
		Reason:      "", // Invalid
		PerformedBy: 1,
	}

	err := suite.service.AdjustInventory(adjustment)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "reason is required")
}

// Test: Invalid adjustment - invalid reason
func (suite *InventoryServiceTestSuite) TestAdjustInventory_InvalidReason_Error() {
	adjustment := &dto.InventoryAdjustment{
		ProductID:   1,
		LocationID:  1,
		Quantity:    10,
		Reason:      "invalid_reason", // Not in allowed list
		PerformedBy: 1,
	}

	err := suite.service.AdjustInventory(adjustment)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid reason")
}

// Test: Transfer with negative quantity (should fail validation)
func (suite *InventoryServiceTestSuite) TestTransferInventory_NegativeQuantity_Error() {
	transfer := &dto.InventoryTransfer{
		ProductID:      1,
		FromLocationID: 1,
		ToLocationID:   2,
		Quantity:       -10, // Invalid
		PerformedBy:    1,
	}

	err := suite.service.TransferInventory(transfer)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "must be positive")
}

// Test: Repository error handling
func (suite *InventoryServiceTestSuite) TestAdjustInventory_RepositoryError() {
	adjustment := &dto.InventoryAdjustment{
		ProductID:   1,
		LocationID:  1,
		Quantity:    10,
		Reason:      "adjustment",
		PerformedBy: 1,
	}

	suite.inventoryRepo.On("GetByProductAndLocation", int64(1), int64(1)).
		Return(nil, errors.New("database connection failed"))

	err := suite.service.AdjustInventory(adjustment)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get current inventory")
	suite.inventoryRepo.AssertExpectations(suite.T())
}

// Test: Get by invalid product ID
func (suite *InventoryServiceTestSuite) TestGetByProduct_InvalidID_Error() {
	// GetByProduct returns ([]*domain.Inventory, error), so we need to handle both return values
	_, err := suite.service.GetByProduct(0)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid product ID")
}

// Test: Get by invalid location ID
func (suite *InventoryServiceTestSuite) TestGetByLocation_InvalidID_Error() {
	_, err := suite.service.GetByLocation(-1) // Invalid ID

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid location ID")
}
