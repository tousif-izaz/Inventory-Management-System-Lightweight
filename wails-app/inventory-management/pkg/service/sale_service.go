package service

import (
	"fmt"
	"inventory-management/pkg/domain"
	"inventory-management/pkg/repository"
	"inventory-management/pkg/service/dto"
	"time"
)

type SaleService struct {
	saleRepo        *repository.SaleRepository
	productRepo     repository.IProductRepository
	transactionRepo repository.ITransactionRepository
}

func NewSaleService(saleRepo *repository.SaleRepository, productRepo repository.IProductRepository, transactionRepo repository.ITransactionRepository) *SaleService {
	return &SaleService{
		saleRepo:        saleRepo,
		productRepo:     productRepo,
		transactionRepo: transactionRepo,
	}
}

// CreateSale creates a new sale with inventory adjustment
func (s *SaleService) CreateSale(createDTO *dto.SaleCreate, userID int64) (*domain.SaleWithItems, error) {
	// 1. Validate stock availability for all items
	for _, item := range createDTO.Items {
		product, err := s.productRepo.GetProductByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found: %w", item.ProductID, err)
		}

		if product.CurrentQuantity < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s (SKU: %s): available=%d, requested=%d",
				product.Name, product.SKU, product.CurrentQuantity, item.Quantity)
		}

		if product.IsActive == false {
			return nil, fmt.Errorf("product %s (SKU: %s) is inactive", product.Name, product.SKU)
		}
	}

	// 2. Calculate totals
	var totalAmount float64
	var calculatedTaxAmount float64
	var calculatedDiscountAmount float64

	for _, item := range createDTO.Items {
		// Line total calculation: (quantity * unit_price) * (1 - discount_percent/100) * (1 + tax_rate/100)
		subtotal := float64(item.Quantity) * item.UnitPrice
		discountAmount := subtotal * (item.DiscountPercent / 100)
		afterDiscount := subtotal - discountAmount
		taxAmount := afterDiscount * (item.TaxRate / 100)
		lineTotal := afterDiscount + taxAmount

		totalAmount += lineTotal
		calculatedTaxAmount += taxAmount
		calculatedDiscountAmount += discountAmount
	}

	netAmount := totalAmount

	// 3. Create Sale record
	sale := &domain.Sale{
		SaleDate:       createDTO.SaleDate,
		ReceiptNo:      createDTO.ReceiptNo,
		CustomerID:     createDTO.CustomerID,
		TotalAmount:    totalAmount,
		TaxAmount:      calculatedTaxAmount,
		DiscountAmount: calculatedDiscountAmount,
		NetAmount:      netAmount,
		PaymentStatus:  createDTO.PaymentStatus,
		PaymentMethod:  createDTO.PaymentMethod,
		Notes:          createDTO.Notes,
		SoldBy:         createDTO.SoldBy,
		CreatedAt:      time.Now(),
	}

	saleID, err := s.saleRepo.CreateSale(sale)
	if err != nil {
		return nil, fmt.Errorf("failed to create sale: %w", err)
	}

	sale.SaleID = saleID

	// 4. Create SaleItems and adjust inventory
	var items []domain.SaleItem
	for _, itemDTO := range createDTO.Items {
		// Calculate line total
		subtotal := float64(itemDTO.Quantity) * itemDTO.UnitPrice
		discountAmount := subtotal * (itemDTO.DiscountPercent / 100)
		afterDiscount := subtotal - discountAmount
		taxAmount := afterDiscount * (itemDTO.TaxRate / 100)
		lineTotal := afterDiscount + taxAmount

		item := &domain.SaleItem{
			SaleID:          saleID,
			ProductID:       itemDTO.ProductID,
			Quantity:        itemDTO.Quantity,
			UnitPrice:       itemDTO.UnitPrice,
			TaxRate:         itemDTO.TaxRate,
			DiscountPercent: itemDTO.DiscountPercent,
			LineTotal:       lineTotal,
		}

		err = s.saleRepo.CreateSaleItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to create sale item: %w", err)
		}

		// 5. Reduce product inventory
		product, _ := s.productRepo.GetProductByID(itemDTO.ProductID)
		previousQuantity := product.CurrentQuantity
		newQuantity := previousQuantity - itemDTO.Quantity

		err = s.productRepo.UpdateProductQuantity(itemDTO.ProductID, newQuantity)
		if err != nil {
			return nil, fmt.Errorf("failed to update product quantity: %w", err)
		}

		// 6. Create transaction audit record
		transactionRecord := &domain.Transaction{
			TransactionType:  "sale",
			TransactionDate:  time.Now(),
			ProductID:        itemDTO.ProductID,
			Quantity:         -itemDTO.Quantity, // Negative for sales
			ReferenceType:    strPtr("sale"),
			ReferenceID:      &saleID,
			PreviousQuantity: &previousQuantity,
			NewQuantity:      &newQuantity,
			PerformedBy:      createDTO.SoldBy,
			Notes:            strPtr(fmt.Sprintf("Sale: %s", createDTO.ReceiptNo)),
		}

		err = s.transactionRepo.Create(transactionRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction record: %w", err)
		}

		items = append(items, *item)
	}

	return &domain.SaleWithItems{
		Sale:  *sale,
		Items: items,
	}, nil
}

// GetSaleByID retrieves a sale by ID
func (s *SaleService) GetSaleByID(id int64) (*domain.SaleWithItems, error) {
	return s.saleRepo.GetSaleByID(id)
}

// GetAllSales retrieves all sales with optional filters
func (s *SaleService) GetAllSales(customerID *int64, paymentStatus *string, fromDate, toDate *time.Time) ([]domain.Sale, error) {
	return s.saleRepo.GetAllSales(customerID, paymentStatus, fromDate, toDate)
}

// UpdatePaymentStatus updates the payment status of a sale
func (s *SaleService) UpdatePaymentStatus(id int64, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"pending":  true,
		"partial":  true,
		"paid":     true,
		"refunded": true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid payment status: %s", status)
	}

	return s.saleRepo.UpdatePaymentStatus(id, status)
}

// ProcessRefund processes a refund and restores inventory
func (s *SaleService) ProcessRefund(id int64, userID int64) error {
	// Get the sale with items
	saleWithItems, err := s.saleRepo.GetSaleByID(id)
	if err != nil {
		return fmt.Errorf("failed to get sale: %w", err)
	}

	// Check if already refunded
	if saleWithItems.PaymentStatus == "refunded" {
		return fmt.Errorf("sale is already refunded")
	}

	// Restore inventory for each item
	for _, item := range saleWithItems.Items {
		product, err := s.productRepo.GetProductByID(item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product: %w", err)
		}

		previousQuantity := product.CurrentQuantity
		newQuantity := previousQuantity + item.Quantity

		err = s.productRepo.UpdateProductQuantity(item.ProductID, newQuantity)
		if err != nil {
			return fmt.Errorf("failed to restore product quantity: %w", err)
		}

		// Create transaction audit record for return
		transactionRecord := &domain.Transaction{
			TransactionType:  "return",
			TransactionDate:  time.Now(),
			ProductID:        item.ProductID,
			Quantity:         item.Quantity, // Positive for returns
			ReferenceType:    strPtr("sale"),
			ReferenceID:      &id,
			PreviousQuantity: &previousQuantity,
			NewQuantity:      &newQuantity,
			PerformedBy:      &userID,
			Notes:            strPtr(fmt.Sprintf("Refund for sale: %s", saleWithItems.ReceiptNo)),
		}

		err = s.transactionRepo.Create(transactionRecord)
		if err != nil {
			return fmt.Errorf("failed to create transaction record: %w", err)
		}
	}

	// Update sale payment status to refunded
	err = s.saleRepo.UpdatePaymentStatus(id, "refunded")
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

// GetSalesByCustomerID retrieves all sales for a customer
func (s *SaleService) GetSalesByCustomerID(customerID int64) ([]domain.Sale, error) {
	return s.saleRepo.GetSalesByCustomerID(customerID)
}

// Helper function to create string pointer
func strPtr(s string) *string {
	return &s
}
