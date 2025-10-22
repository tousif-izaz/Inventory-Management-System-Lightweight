package repository

import (
	"database/sql"
	"fmt"
	"ims-intro/pkg/domain"

	"github.com/labstack/gommon/log"
)

type IProductRepository interface {
	// Basic CRUD
	GetAllProducts() ([]*domain.Product, error)
	GetAllProductsWithDetails() ([]*domain.ProductWithDetails, error)
	GetProductByID(productID int64) (*domain.Product, error)
	GetProductBySKU(sku string) (*domain.Product, error)
	GetProductsByCategory(categoryID int64) ([]*domain.Product, error)
	AddProduct(product *domain.Product) error
	UpdateProductById(updatedProduct *domain.Product, productId int64) error
	DeleteProductById(productId int64) error

	// Stock management queries
	GetLowStockProducts() ([]*domain.Product, error)
	GetExpiringProducts(days int) ([]*domain.Product, error)
	UpdateProductQuantity(productID int64, newQuantity int64) error

	// Utility
	CheckProductExistence(productId int64) error
}

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) IProductRepository {
	return &ProductRepository{db}
}

func (repository *ProductRepository) GetAllProducts() ([]*domain.Product, error) {
	query := `
		SELECT ProductID, Name, Description, SKU, CategoryID, BatchNo, ExpiryDate,
		       CostPrice, SellingPrice, CurrentQuantity, MinStockLevel, MaxStockLevel, ReorderPoint,
		       Unit, ShelfLocation, IsActive, CreatedAt, UpdatedAt
		FROM Products
		WHERE IsActive = 1
		ORDER BY ProductID
	`

	productRows, err := repository.db.Query(query)
	if err != nil {
		log.Errorf("error while getting all products: %v", err)
		return nil, err
	}
	defer productRows.Close()

	return extractProductsFromRows(productRows)
}

func (repository *ProductRepository) GetAllProductsWithDetails() ([]*domain.ProductWithDetails, error) {
	query := `
		SELECT
			p.ProductID, p.Name, p.Description, p.SKU, p.CategoryID,
			p.BatchNo, p.ExpiryDate, p.CostPrice, p.SellingPrice, p.CurrentQuantity,
			p.MinStockLevel, p.MaxStockLevel, p.ReorderPoint, p.Unit, p.ShelfLocation,
			p.IsActive, p.CreatedAt, p.UpdatedAt,
			c.Name as CategoryName
		FROM Products p
		LEFT JOIN Categories c ON p.CategoryID = c.CategoryID
		WHERE p.IsActive = 1
		ORDER BY p.Name
	`

	rows, err := repository.db.Query(query)
	if err != nil {
		log.Errorf("error while getting all products with details: %v", err)
		return nil, err
	}
	defer rows.Close()

	return extractProductsWithDetailsFromRows(rows)
}

func (repository *ProductRepository) GetProductByID(productID int64) (*domain.Product, error) {
	query := `
		SELECT ProductID, Name, Description, SKU, CategoryID, BatchNo, ExpiryDate,
		       CostPrice, SellingPrice, CurrentQuantity, MinStockLevel, MaxStockLevel, ReorderPoint,
		       Unit, ShelfLocation, IsActive, CreatedAt, UpdatedAt
		FROM Products
		WHERE ProductID = ?
	`

	row := repository.db.QueryRow(query, productID)
	return scanProductRow(row)
}

func (repository *ProductRepository) GetProductBySKU(sku string) (*domain.Product, error) {
	query := `
		SELECT ProductID, Name, Description, SKU, CategoryID, BatchNo, ExpiryDate,
		       CostPrice, SellingPrice, CurrentQuantity, MinStockLevel, MaxStockLevel, ReorderPoint,
		       Unit, ShelfLocation, IsActive, CreatedAt, UpdatedAt
		FROM Products
		WHERE SKU = ? AND IsActive = 1
	`

	row := repository.db.QueryRow(query, sku)
	return scanProductRow(row)
}

func (repository *ProductRepository) GetProductsByCategory(categoryID int64) ([]*domain.Product, error) {
	query := `
		SELECT ProductID, Name, Description, SKU, CategoryID, BatchNo, ExpiryDate,
		       CostPrice, SellingPrice, CurrentQuantity, MinStockLevel, MaxStockLevel, ReorderPoint,
		       Unit, ShelfLocation, IsActive, CreatedAt, UpdatedAt
		FROM Products
		WHERE CategoryID = ? AND IsActive = 1
		ORDER BY Name
	`

	productRows, err := repository.db.Query(query, categoryID)
	if err != nil {
		log.Errorf("error while getting products by category: %v", err)
		return nil, err
	}
	defer productRows.Close()

	return extractProductsFromRows(productRows)
}

func (repository *ProductRepository) AddProduct(product *domain.Product) error {
	insertStatement := `
		INSERT INTO Products (
			Name, Description, SKU, CategoryID, BatchNo, ExpiryDate,
			CostPrice, SellingPrice, CurrentQuantity, MinStockLevel, MaxStockLevel, ReorderPoint,
			Unit, ShelfLocation, IsActive
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1)
	`

	result, err := repository.db.Exec(
		insertStatement,
		product.Name,
		product.Description,
		product.SKU,
		product.CategoryID,
		product.BatchNo,
		product.ExpiryDate,
		product.CostPrice,
		product.SellingPrice,
		product.CurrentQuantity,
		product.MinStockLevel,
		product.MaxStockLevel,
		product.ReorderPoint,
		product.Unit,
		product.ShelfLocation,
	)
	if err != nil {
		log.Errorf("error while adding a new product: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Product added successfully: %v", result))
	return nil
}

func (repository *ProductRepository) UpdateProductById(updatedProduct *domain.Product, productId int64) error {
	updateStatement := `
		UPDATE Products
		SET Name = ?, Description = ?, SKU = ?, CategoryID = ?,
		    BatchNo = ?, ExpiryDate = ?, CostPrice = ?, SellingPrice = ?, CurrentQuantity = ?,
		    MinStockLevel = ?, MaxStockLevel = ?, ReorderPoint = ?, Unit = ?, ShelfLocation = ?
		WHERE ProductID = ?
	`

	result, err := repository.db.Exec(
		updateStatement,
		updatedProduct.Name,
		updatedProduct.Description,
		updatedProduct.SKU,
		updatedProduct.CategoryID,
		updatedProduct.BatchNo,
		updatedProduct.ExpiryDate,
		updatedProduct.CostPrice,
		updatedProduct.SellingPrice,
		updatedProduct.CurrentQuantity,
		updatedProduct.MinStockLevel,
		updatedProduct.MaxStockLevel,
		updatedProduct.ReorderPoint,
		updatedProduct.Unit,
		updatedProduct.ShelfLocation,
		productId,
	)
	if err != nil {
		log.Errorf("error while updating product: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Product updated successfully: %v", result))
	return nil
}

func (repository *ProductRepository) DeleteProductById(productId int64) error {
	// Hard delete - permanently remove from database
	deleteStatement := "DELETE FROM Products WHERE ProductID = ?"
	result, err := repository.db.Exec(deleteStatement, productId)
	if err != nil {
		log.Errorf("error while deleting product: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Info("Product deleted successfully (hard delete)")
	log.Info(fmt.Sprintf("%v rows affected", rowsAffected))

	return nil
}

func (repository *ProductRepository) GetLowStockProducts() ([]*domain.Product, error) {
	query := `
		SELECT ProductID, Name, Description, SKU, CategoryID,
		       BatchNo, ExpiryDate, CostPrice, SellingPrice, CurrentQuantity, MinStockLevel,
		       MaxStockLevel, ReorderPoint, Unit, ShelfLocation, IsActive, CreatedAt, UpdatedAt
		FROM Products
		WHERE IsActive = 1 AND CurrentQuantity <= ReorderPoint
		ORDER BY CurrentQuantity ASC
	`

	productRows, err := repository.db.Query(query)
	if err != nil {
		log.Errorf("error while getting low stock products: %v", err)
		return nil, err
	}
	defer productRows.Close()

	return extractProductsFromRows(productRows)
}

func (repository *ProductRepository) GetExpiringProducts(days int) ([]*domain.Product, error) {
	query := `
		SELECT ProductID, Name, Description, SKU, CategoryID, BatchNo, ExpiryDate,
		       CostPrice, SellingPrice, CurrentQuantity, MinStockLevel, MaxStockLevel, ReorderPoint,
		       Unit, ShelfLocation, IsActive, CreatedAt, UpdatedAt
		FROM Products
		WHERE IsActive = 1
		  AND ExpiryDate IS NOT NULL
		  AND ExpiryDate <= date('now', '+' || ? || ' days')
		ORDER BY ExpiryDate ASC
	`

	productRows, err := repository.db.Query(query, days)
	if err != nil {
		log.Errorf("error while getting expiring products: %v", err)
		return nil, err
	}
	defer productRows.Close()

	return extractProductsFromRows(productRows)
}

func (repository *ProductRepository) UpdateProductQuantity(productID int64, newQuantity int64) error {
	updateStatement := `
		UPDATE Products
		SET CurrentQuantity = ?, UpdatedAt = CURRENT_TIMESTAMP
		WHERE ProductID = ? AND IsActive = 1
	`

	result, err := repository.db.Exec(updateStatement, newQuantity, productID)
	if err != nil {
		log.Errorf("error while updating product quantity: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("error while checking rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found or is inactive", productID)
	}

	log.Info(fmt.Sprintf("Product quantity updated successfully: ProductID=%d, NewQuantity=%d", productID, newQuantity))
	return nil
}

func (repository *ProductRepository) CheckProductExistence(productId int64) error {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM Products WHERE ProductID = ? AND IsActive = 1)"
	err := repository.db.QueryRow(query, productId).Scan(&exists)
	if err != nil {
		log.Errorf("error while checking product existence: %v", err)
		return err
	}

	if !exists {
		return fmt.Errorf("product with id %d does not exist or is inactive", productId)
	}

	return nil
}

// Helper functions

func scanProductRow(row *sql.Row) (*domain.Product, error) {
	product := &domain.Product{}
	var description sql.NullString
	var batchNo sql.NullString
	var expiryDate sql.NullTime
	var maxStockLevel sql.NullInt64
	var shelfLocation sql.NullString

	err := row.Scan(
		&product.ProductID,
		&product.Name,
		&description,
		&product.SKU,
		&product.CategoryID,
		&batchNo,
		&expiryDate,
		&product.CostPrice,
		&product.SellingPrice,
		&product.CurrentQuantity,
		&product.MinStockLevel,
		&maxStockLevel,
		&product.ReorderPoint,
		&product.Unit,
		&shelfLocation,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}

	if err != nil {
		log.Errorf("error while scanning product: %v", err)
		return nil, err
	}

	// Handle nullable fields
	if description.Valid {
		product.Description = &description.String
	}
	if batchNo.Valid {
		product.BatchNo = &batchNo.String
	}
	if expiryDate.Valid {
		product.ExpiryDate = &expiryDate.Time
	}
	if maxStockLevel.Valid {
		product.MaxStockLevel = &maxStockLevel.Int64
	}
	if shelfLocation.Valid {
		product.ShelfLocation = &shelfLocation.String
	}

	return product, nil
}

func extractProductsFromRows(productRows *sql.Rows) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)

	for productRows.Next() {
		product := &domain.Product{}
		var description sql.NullString
		var batchNo sql.NullString
		var expiryDate sql.NullTime
		var maxStockLevel sql.NullInt64
		var shelfLocation sql.NullString

		err := productRows.Scan(
			&product.ProductID,
			&product.Name,
			&description,
			&product.SKU,
			&product.CategoryID,
			&batchNo,
			&expiryDate,
			&product.CostPrice,
			&product.SellingPrice,
			&product.CurrentQuantity,
			&product.MinStockLevel,
			&maxStockLevel,
			&product.ReorderPoint,
			&product.Unit,
			&shelfLocation,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			log.Errorf("error while scanning product row: %v", err)
			continue
		}

		// Handle nullable fields
		if description.Valid {
			product.Description = &description.String
		}
		if batchNo.Valid {
			product.BatchNo = &batchNo.String
		}
		if expiryDate.Valid {
			product.ExpiryDate = &expiryDate.Time
		}
		if maxStockLevel.Valid {
			product.MaxStockLevel = &maxStockLevel.Int64
		}
		if shelfLocation.Valid {
			product.ShelfLocation = &shelfLocation.String
		}

		products = append(products, product)
	}

	if err := productRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
}

func extractProductsWithDetailsFromRows(rows *sql.Rows) ([]*domain.ProductWithDetails, error) {
	products := make([]*domain.ProductWithDetails, 0)

	for rows.Next() {
		product := &domain.ProductWithDetails{}
		var description sql.NullString
		var batchNo sql.NullString
		var expiryDate sql.NullTime
		var maxStockLevel sql.NullInt64
		var shelfLocation sql.NullString
		var categoryName sql.NullString

		err := rows.Scan(
			&product.ProductID,
			&product.Name,
			&description,
			&product.SKU,
			&product.CategoryID,
			&batchNo,
			&expiryDate,
			&product.CostPrice,
			&product.SellingPrice,
			&product.CurrentQuantity,
			&product.MinStockLevel,
			&maxStockLevel,
			&product.ReorderPoint,
			&product.Unit,
			&shelfLocation,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
			&categoryName,
		)

		if err != nil {
			log.Errorf("error while scanning product with details row: %v", err)
			continue
		}

		// Handle nullable fields
		if description.Valid {
			product.Description = &description.String
		}
		if batchNo.Valid {
			product.BatchNo = &batchNo.String
		}
		if expiryDate.Valid {
			product.ExpiryDate = &expiryDate.Time
		}
		if maxStockLevel.Valid {
			product.MaxStockLevel = &maxStockLevel.Int64
		}
		if shelfLocation.Valid {
			product.ShelfLocation = &shelfLocation.String
		}
		if categoryName.Valid {
			product.CategoryName = &categoryName.String
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product with details rows: %w", err)
	}

	return products, nil
}
