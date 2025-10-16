package repository

import (
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
	"ims-intro/pkg/domain"
)

type IProductRepository interface {
	GetAllProducts() []*domain.Product
	GetProductsByCategory(category string) []*domain.Product
	AddProduct(product *domain.Product) error
	CheckProductExistence(productId int64) error
	UpdateProductById(updatedProduct *domain.Product, productId int64) error
	DeleteProductById(productId int64) error
}

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) IProductRepository {
	return &ProductRepository{db}
}

func (repository *ProductRepository) GetAllProducts() []*domain.Product {
	productRows, err := repository.db.Query("SELECT id, name, price, quantity, category FROM products")
	if err != nil {
		log.Errorf("error while getting all products: %v", err)
		return make([]*domain.Product, 0)
	}
	defer productRows.Close()

	return extractProductsFromRows(productRows)
}

func (repository *ProductRepository) GetProductsByCategory(category string) []*domain.Product {
	productRows, err := repository.db.Query("SELECT id, name, price, quantity, category FROM products WHERE category = ?", category)
	if err != nil {
		log.Errorf("error while getting all products by category: %v", err)
		return make([]*domain.Product, 0)
	}
	defer productRows.Close()

	return extractProductsFromRows(productRows)
}

func (repository *ProductRepository) AddProduct(product *domain.Product) error {
	insertStatement := "INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)"

	result, err := repository.db.Exec(insertStatement, product.Name, product.Price, product.Quantity, product.Category)
	if err != nil {
		log.Errorf("error while adding a new product: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Product added successfully: %v", result))
	return nil
}

func (repository *ProductRepository) CheckProductExistence(productId int64) error {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM products WHERE id = ?)"
	err := repository.db.QueryRow(query, productId).Scan(&exists)
	if err != nil {
		log.Errorf("error while checking product existence: %v", err)
		return err
	}

	if !exists {
		return fmt.Errorf("product with id %d does not exist", productId)
	}

	return nil
}

func (repository *ProductRepository) UpdateProductById(updatedProduct *domain.Product, productId int64) error {
	updateStatement := "UPDATE products SET name = ?, price = ?, quantity = ?, category = ? WHERE id = ?"
	result, err := repository.db.Exec(updateStatement, updatedProduct.Name, updatedProduct.Price, updatedProduct.Quantity, updatedProduct.Category, productId)
	if err != nil {
		log.Errorf("error while updating product: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Product updated successfully: %v", result))
	return nil
}

func (repository *ProductRepository) DeleteProductById(productId int64) error {
	result, err := repository.db.Exec("DELETE FROM products WHERE id = ?", productId)
	if err != nil {
		log.Errorf("error while deleting product: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Info("Product deleted successfully")
	log.Info(fmt.Sprintf("%v rows affected", rowsAffected))

	return nil
}

func extractProductsFromRows(productRows *sql.Rows) []*domain.Product {
	products := make([]*domain.Product, 0)

	for productRows.Next() {
		product := &domain.Product{}
		err := productRows.Scan(&product.Id, &product.Name, &product.Price, &product.Quantity, &product.Category)
		if err != nil {
			log.Errorf("error while scanning product row: %v", err)
			continue
		}
		products = append(products, product)
	}

	return products
}
