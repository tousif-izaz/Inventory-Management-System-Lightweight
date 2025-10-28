package testutil

import (
	"database/sql"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"inventory-management/pkg/domain"
)

// SetupMockDB creates a mock database connection for testing
func SetupMockDB() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

// Helper functions to create test domain objects

func CreateTestProduct(id int64) *domain.Product {
	desc := "Test product description"
	batchNo := "BATCH-001"
	expiryDate := time.Now().AddDate(1, 0, 0)
	maxStock := int64(100)
	shelfLocation := "A1-B2"

	return &domain.Product{
		ProductID:       id,
		Name:            "Test Product",
		Description:     &desc,
		SKU:             "TEST-SKU-001",
		CategoryID:      1,
		BatchNo:         &batchNo,
		ExpiryDate:      &expiryDate,
		CostPrice:       50.00,
		SellingPrice:    89.99,
		CurrentQuantity: 50,
		MinStockLevel:   10,
		MaxStockLevel:   &maxStock,
		ReorderPoint:    20,
		Unit:            "pcs",
		ShelfLocation:   &shelfLocation,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func CreateTestCategory(id int64, name string, parentID *int64) *domain.Category {
	desc := "Test category"
	return &domain.Category{
		CategoryID:       id,
		Name:             name,
		Description:      &desc,
		ParentCategoryID: parentID,
		CreatedAt:        time.Now(),
	}
}

func CreateTestSupplier(id int64) *domain.Supplier {
	contact := "John Doe"
	email := "supplier@test.com"
	phone := "+1234567890"
	address := "123 Test St"
	city := "Test City"
	country := "Test Country"
	taxID := "TAX123"
	terms := "Net 30"

	return &domain.Supplier{
		SupplierID:    id,
		Name:          "Test Supplier",
		ContactPerson: &contact,
		Email:         &email,
		Phone:         &phone,
		Address:       &address,
		City:          &city,
		Country:       &country,
		TaxID:         &taxID,
		PaymentTerms:  &terms,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func CreateTestCustomer(id int64) *domain.Customer {
	email := "customer@test.com"
	phone := "+1234567890"
	address := "789 Customer St"
	city := "Test City"
	country := "Test Country"

	return &domain.Customer{
		CustomerID:    id,
		Name:          "Test Customer",
		Email:         &email,
		Phone:         &phone,
		Address:       &address,
		City:          &city,
		Country:       &country,
		LoyaltyPoints: 100,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func CreateTestInventory(id, productID, quantity int64) *domain.Inventory {
	return &domain.Inventory{
		InventoryID:       id,
		ProductID:         productID,
		Quantity:          quantity,
		ReservedQuantity:  0,
		AvailableQuantity: quantity,
		LastUpdated:       time.Now(),
	}
}

func CreateTestSale(id int64) *domain.Sale {
	customerID := int64(1)
	method := "card"
	notes := "Test sale"
	soldBy := int64(1)

	return &domain.Sale{
		SaleID:         id,
		SaleDate:       time.Now(),
		ReceiptNo:      "RCP-001",
		CustomerID:     &customerID,
		TotalAmount:    179.98,
		TaxAmount:      16.36,
		DiscountAmount: 0,
		NetAmount:      196.34,
		PaymentStatus:  "paid",
		PaymentMethod:  &method,
		Notes:          &notes,
		SoldBy:         &soldBy,
		CreatedAt:      time.Now(),
	}
}

func CreateTestSaleItem(id, saleID, productID int64) domain.SaleItem {
	locationID := int64(1)

	return domain.SaleItem{
		SaleItemID:      id,
		SaleID:          saleID,
		ProductID:       productID,
		Quantity:        2,
		UnitPrice:       89.99,
		TaxRate:         10.0,
		DiscountPercent: 0,
		LineTotal:       179.98,
		LocationID:      &locationID,
	}
}

func CreateTestPurchase(id int64) *domain.Purchase {
	ref := "REF-001"
	method := "bank_transfer"
	notes := "Test purchase"
	receivedBy := int64(1)

	return &domain.Purchase{
		PurchaseID:      id,
		PurchaseDate:    time.Now(),
		SupplierID:      1,
		InvoiceNumber:   "INV-001",
		ReferenceNumber: &ref,
		TotalAmount:     5000.00,
		TaxAmount:       500.00,
		DiscountAmount:  0,
		NetAmount:       5500.00,
		PaymentStatus:   "pending",
		PaymentMethod:   &method,
		Notes:           &notes,
		ReceivedBy:      &receivedBy,
		CreatedAt:       time.Now(),
	}
}

func CreateTestPurchaseItem(id, purchaseID, productID int64) domain.PurchaseItem {
	locationID := int64(1)
	batchNo := "BATCH-001"
	expiryDate := time.Now().AddDate(1, 0, 0)

	return domain.PurchaseItem{
		PurchaseItemID:  id,
		PurchaseID:      purchaseID,
		ProductID:       productID,
		Quantity:        100,
		UnitCost:        50.00,
		TaxRate:         10.0,
		DiscountPercent: 0,
		LineTotal:       5000.00,
		LocationID:      &locationID,
		BatchNo:         &batchNo,
		ExpiryDate:      &expiryDate,
	}
}

func CreateTestTransaction(id, productID int64, txType string, quantity int64) *domain.Transaction {
	locationID := int64(1)
	refType := "sale"
	refID := int64(1)
	prevQty := int64(100)
	newQty := prevQty + quantity
	reason := "Test transaction"
	performedBy := int64(1)
	notes := "Test notes"

	return &domain.Transaction{
		TransactionID:    id,
		TransactionType:  txType,
		TransactionDate:  time.Now(),
		ProductID:        productID,
		LocationID:       &locationID,
		Quantity:         quantity,
		ReferenceType:    &refType,
		ReferenceID:      &refID,
		PreviousQuantity: &prevQty,
		NewQuantity:      &newQty,
		Reason:           &reason,
		PerformedBy:      &performedBy,
		Notes:            &notes,
		CreatedAt:        time.Now(),
	}
}

// Pointer helpers for optional fields
func StringPtr(s string) *string {
	return &s
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

// MockUpdateProductQuantity creates a mock expectation for UpdateProductQuantity
func MockUpdateProductQuantity(mock sqlmock.Sqlmock, productID, newQuantity int64, shouldFail bool) {
	if shouldFail {
		mock.ExpectExec("UPDATE Products SET CurrentQuantity = \\?, UpdatedAt = CURRENT_TIMESTAMP WHERE ProductID = \\? AND IsActive = 1").
			WithArgs(newQuantity, productID).
			WillReturnError(sql.ErrConnDone)
	} else {
		mock.ExpectExec("UPDATE Products SET CurrentQuantity = \\?, UpdatedAt = CURRENT_TIMESTAMP WHERE ProductID = \\? AND IsActive = 1").
			WithArgs(newQuantity, productID).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
}
