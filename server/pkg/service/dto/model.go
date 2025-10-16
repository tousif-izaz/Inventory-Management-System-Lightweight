package dto

import "time"

// User DTOs
type UserCreate struct {
	Username string
	Name     string
	Email    *string
	Password string
	Role     string
}

// Product DTOs
type ProductCreate struct {
	Name          string
	Description   *string
	SKU           string
	Barcode       *string
	CategoryID    int64
	BatchNo       *string
	ExpiryDate    *time.Time
	CostPrice     float64
	SellingPrice  float64
	MinStockLevel int64
	MaxStockLevel *int64
	ReorderPoint  int64
	Unit          string
}

// Category DTOs
type CategoryCreate struct {
	Name             string
	Description      *string
	ParentCategoryID *int64
}

// Supplier DTOs
type SupplierCreate struct {
	Name          string
	ContactPerson *string
	Email         *string
	Phone         *string
	Address       *string
	City          *string
	Country       *string
	TaxID         *string
	PaymentTerms  *string
}

// Location DTOs
type LocationCreate struct {
	Name    string
	Type    *string
	Address *string
}

// Customer DTOs
type CustomerCreate struct {
	Name    string
	Email   *string
	Phone   *string
	Address *string
	City    *string
	Country *string
}

// Inventory DTOs
type InventoryAdjustment struct {
	ProductID  int64
	LocationID int64
	Quantity   int64 // Can be negative for reductions
	Reason     string
	Notes      *string
	PerformedBy int64
}

type InventoryTransfer struct {
	ProductID         int64
	FromLocationID    int64
	ToLocationID      int64
	Quantity          int64
	Reason            *string
	PerformedBy       int64
}

// Purchase DTOs
type PurchaseCreate struct {
	PurchaseDate    time.Time
	SupplierID      int64
	InvoiceNumber   string
	ReferenceNumber *string
	TaxAmount       float64
	DiscountAmount  float64
	PaymentStatus   string
	PaymentMethod   *string
	Notes           *string
	ReceivedBy      *int64
	Items           []PurchaseItemCreate
}

type PurchaseItemCreate struct {
	ProductID       int64
	Quantity        int64
	UnitCost        float64
	TaxRate         float64
	DiscountPercent float64
	LocationID      *int64
	BatchNo         *string
	ExpiryDate      *time.Time
}

// Sale DTOs
type SaleCreate struct {
	SaleDate       time.Time
	ReceiptNo      string
	CustomerID     *int64
	TaxAmount      float64
	DiscountAmount float64
	PaymentStatus  string
	PaymentMethod  *string
	Notes          *string
	SoldBy         *int64
	Items          []SaleItemCreate
}

type SaleItemCreate struct {
	ProductID       int64
	Quantity        int64
	UnitPrice       float64
	TaxRate         float64
	DiscountPercent float64
	LocationID      *int64
}

// Transaction DTOs
type TransactionCreate struct {
	TransactionType string
	ProductID       int64
	LocationID      *int64
	Quantity        int64
	ReferenceType   *string
	ReferenceID     *int64
	Reason          *string
	PerformedBy     *int64
	Notes           *string
}
