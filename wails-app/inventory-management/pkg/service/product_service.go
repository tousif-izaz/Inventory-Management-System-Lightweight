package service

import (
	"errors"
	"fmt"
	"inventory-management/pkg/domain"
	"inventory-management/pkg/repository"
	"inventory-management/pkg/service/dto"
)

type IProductService interface {
	Add(productCreate *dto.ProductCreate) error
	GetAllProducts() ([]*domain.Product, error)
	GetAllProductsWithDetails() ([]*domain.ProductWithDetails, error)
	GetProductByID(productID int64) (*domain.Product, error)
	GetProductBySKU(sku string) (*domain.Product, error)
	GetAllProductsByCategory(categoryID int64) ([]*domain.Product, error)
	GetLowStockProducts() ([]*domain.Product, error)
	GetExpiringProducts(days int) ([]*domain.Product, error)
	UpdateProductById(updatedProduct *dto.ProductCreate, productId int64) error
	DeleteById(productId int64) error
}

type ProductService struct {
	productRepository repository.IProductRepository
}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &ProductService{productRepository}
}

func (service *ProductService) Add(productCreate *dto.ProductCreate) error {
	err := validateProductCreate(productCreate)
	if err != nil {
		return err
	}

	product := productCreateToProduct(productCreate)
	return service.productRepository.AddProduct(product)
}

func (service *ProductService) GetAllProducts() ([]*domain.Product, error) {
	return service.productRepository.GetAllProducts()
}

func (service *ProductService) GetAllProductsWithDetails() ([]*domain.ProductWithDetails, error) {
	return service.productRepository.GetAllProductsWithDetails()
}

func (service *ProductService) GetProductByID(productID int64) (*domain.Product, error) {
	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}
	return service.productRepository.GetProductByID(productID)
}

func (service *ProductService) GetProductBySKU(sku string) (*domain.Product, error) {
	if sku == "" {
		return nil, errors.New("SKU cannot be empty")
	}
	return service.productRepository.GetProductBySKU(sku)
}

func (service *ProductService) GetAllProductsByCategory(categoryID int64) ([]*domain.Product, error) {
	if categoryID <= 0 {
		return nil, errors.New("invalid category ID")
	}
	return service.productRepository.GetProductsByCategory(categoryID)
}

func (service *ProductService) GetLowStockProducts() ([]*domain.Product, error) {
	return service.productRepository.GetLowStockProducts()
}

func (service *ProductService) GetExpiringProducts(days int) ([]*domain.Product, error) {
	if days < 0 {
		return nil, errors.New("days must be non-negative")
	}
	if days > 365 {
		days = 365 // Cap at 1 year
	}
	return service.productRepository.GetExpiringProducts(days)
}

func (service *ProductService) UpdateProductById(updatedProduct *dto.ProductCreate, productId int64) error {
	err := service.productRepository.CheckProductExistence(productId)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	err = validateProductCreate(updatedProduct)
	if err != nil {
		return err
	}

	product := productCreateToProduct(updatedProduct)
	return service.productRepository.UpdateProductById(product, productId)
}

func (service *ProductService) DeleteById(productId int64) error {
	err := service.productRepository.CheckProductExistence(productId)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	return service.productRepository.DeleteProductById(productId)
}

func validateProductCreate(productCreate *dto.ProductCreate) error {
	if productCreate.Name == "" {
		return errors.New("name can't be empty")
	}
	if productCreate.SKU == "" {
		return errors.New("SKU can't be empty")
	}
	if productCreate.CategoryID <= 0 {
		return errors.New("valid category ID is required")
	}
	if productCreate.CostPrice < 0 {
		return errors.New("cost price can't be less than zero")
	}
	if productCreate.SellingPrice < 0 {
		return errors.New("selling price can't be less than zero")
	}
	if productCreate.MinStockLevel < 0 {
		return errors.New("min stock level can't be less than zero")
	}
	if productCreate.ReorderPoint < 0 {
		return errors.New("reorder point can't be less than zero")
	}
	if productCreate.Unit == "" {
		return errors.New("unit can't be empty")
	}

	// Validate unit is one of the allowed values
	validUnits := map[string]bool{
		"pcs":    true,
		"kg":     true,
		"liter":  true,
		"box":    true,
		"carton": true,
	}
	if !validUnits[productCreate.Unit] {
		return errors.New("unit must be one of: pcs, kg, liter, box, carton")
	}

	return nil
}

func productCreateToProduct(productCreate *dto.ProductCreate) *domain.Product {
	return &domain.Product{
		Name:            productCreate.Name,
		Description:     productCreate.Description,
		SKU:             productCreate.SKU,
		CategoryID:      productCreate.CategoryID,
		BatchNo:         productCreate.BatchNo,
		ExpiryDate:      productCreate.ExpiryDate,
		CostPrice:       productCreate.CostPrice,
		SellingPrice:    productCreate.SellingPrice,
		CurrentQuantity: productCreate.CurrentQuantity,
		MinStockLevel:   productCreate.MinStockLevel,
		MaxStockLevel:   productCreate.MaxStockLevel,
		ReorderPoint:    productCreate.ReorderPoint,
		Unit:            productCreate.Unit,
		ShelfLocation:   productCreate.ShelfLocation,
		IsActive:        true,
	}
}
