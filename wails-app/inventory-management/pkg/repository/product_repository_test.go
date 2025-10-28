package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"inventory-management/pkg/testutil"
)

func TestUpdateProductQuantity_Success(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(db)
	productID := int64(1)
	newQuantity := int64(75)

	testutil.MockUpdateProductQuantity(mock, productID, newQuantity, false)

	err = repo.UpdateProductQuantity(productID, newQuantity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProductQuantity_ProductNotFound(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(db)
	productID := int64(999)
	newQuantity := int64(50)

	// Mock the update to return 0 rows affected (product not found)
	mock.ExpectExec("UPDATE Products SET CurrentQuantity = \\?, UpdatedAt = CURRENT_TIMESTAMP WHERE ProductID = \\? AND IsActive = 1").
		WithArgs(newQuantity, productID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateProductQuantity(productID, newQuantity)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found or is inactive")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProductQuantity_DatabaseError(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(db)
	productID := int64(1)
	newQuantity := int64(100)

	testutil.MockUpdateProductQuantity(mock, productID, newQuantity, true)

	err = repo.UpdateProductQuantity(productID, newQuantity)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProductQuantity_ZeroQuantity(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(db)
	productID := int64(1)
	newQuantity := int64(0)

	testutil.MockUpdateProductQuantity(mock, productID, newQuantity, false)

	err = repo.UpdateProductQuantity(productID, newQuantity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProductQuantity_LargeQuantity(t *testing.T) {
	db, mock, err := testutil.SetupMockDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(db)
	productID := int64(1)
	newQuantity := int64(999999)

	testutil.MockUpdateProductQuantity(mock, productID, newQuantity, false)

	err = repo.UpdateProductQuantity(productID, newQuantity)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
