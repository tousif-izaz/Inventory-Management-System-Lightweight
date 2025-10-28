package repository

import (
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
	"inventory-management/pkg/domain"
)

type ICategoryRepository interface {
	GetAllCategories() ([]*domain.Category, error)
	GetCategoryByID(categoryID int64) (*domain.Category, error)
	GetChildCategories(parentCategoryID int64) ([]*domain.Category, error)
	AddCategory(category *domain.Category) error
	UpdateCategoryById(updatedCategory *domain.Category, categoryId int64) error
	DeleteCategoryById(categoryId int64) error
	CheckCategoryExistence(categoryId int64) error
}

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) ICategoryRepository {
	return &CategoryRepository{db}
}

func (repository *CategoryRepository) GetAllCategories() ([]*domain.Category, error) {
	query := `
		SELECT CategoryID, Name, Description, ParentCategoryID, CreatedAt
		FROM Categories
		ORDER BY Name
	`

	categoryRows, err := repository.db.Query(query)
	if err != nil {
		log.Errorf("error while getting all categories: %v", err)
		return nil, err
	}
	defer categoryRows.Close()

	return extractCategoriesFromRows(categoryRows)
}

func (repository *CategoryRepository) GetCategoryByID(categoryID int64) (*domain.Category, error) {
	query := `
		SELECT CategoryID, Name, Description, ParentCategoryID, CreatedAt
		FROM Categories
		WHERE CategoryID = ?
	`

	row := repository.db.QueryRow(query, categoryID)
	return scanCategoryRow(row)
}

func (repository *CategoryRepository) GetChildCategories(parentCategoryID int64) ([]*domain.Category, error) {
	query := `
		SELECT CategoryID, Name, Description, ParentCategoryID, CreatedAt
		FROM Categories
		WHERE ParentCategoryID = ?
		ORDER BY Name
	`

	categoryRows, err := repository.db.Query(query, parentCategoryID)
	if err != nil {
		log.Errorf("error while getting child categories: %v", err)
		return nil, err
	}
	defer categoryRows.Close()

	return extractCategoriesFromRows(categoryRows)
}

func (repository *CategoryRepository) AddCategory(category *domain.Category) error {
	insertStatement := `
		INSERT INTO Categories (Name, Description, ParentCategoryID)
		VALUES (?, ?, ?)
	`

	result, err := repository.db.Exec(
		insertStatement,
		category.Name,
		category.Description,
		category.ParentCategoryID,
	)
	if err != nil {
		log.Errorf("error while adding a new category: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Category added successfully: %v", result))
	return nil
}

func (repository *CategoryRepository) UpdateCategoryById(updatedCategory *domain.Category, categoryId int64) error {
	updateStatement := `
		UPDATE Categories
		SET Name = ?, Description = ?, ParentCategoryID = ?
		WHERE CategoryID = ?
	`

	result, err := repository.db.Exec(
		updateStatement,
		updatedCategory.Name,
		updatedCategory.Description,
		updatedCategory.ParentCategoryID,
		categoryId,
	)
	if err != nil {
		log.Errorf("error while updating category: %v", err)
		return err
	}

	log.Info(fmt.Sprintf("Category updated successfully: %v", result))
	return nil
}

func (repository *CategoryRepository) DeleteCategoryById(categoryId int64) error {
	// Hard delete for categories (as per schema)
	// Note: This will fail if there are products using this category (RESTRICT constraint)
	deleteStatement := "DELETE FROM Categories WHERE CategoryID = ?"
	result, err := repository.db.Exec(deleteStatement, categoryId)
	if err != nil {
		log.Errorf("error while deleting category: %v", err)
		return fmt.Errorf("cannot delete category: it may be in use by products")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("category with id %d not found", categoryId)
	}

	log.Info("Category deleted successfully")
	log.Info(fmt.Sprintf("%v rows affected", rowsAffected))

	return nil
}

func (repository *CategoryRepository) CheckCategoryExistence(categoryId int64) error {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM Categories WHERE CategoryID = ?)"
	err := repository.db.QueryRow(query, categoryId).Scan(&exists)
	if err != nil {
		log.Errorf("error while checking category existence: %v", err)
		return err
	}

	if !exists {
		return fmt.Errorf("category with id %d does not exist", categoryId)
	}

	return nil
}

// Helper functions

func scanCategoryRow(row *sql.Row) (*domain.Category, error) {
	category := &domain.Category{}
	var description sql.NullString
	var parentCategoryID sql.NullInt64

	err := row.Scan(
		&category.CategoryID,
		&category.Name,
		&description,
		&parentCategoryID,
		&category.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}

	if err != nil {
		log.Errorf("error while scanning category: %v", err)
		return nil, err
	}

	// Handle nullable fields
	if description.Valid {
		category.Description = &description.String
	}
	if parentCategoryID.Valid {
		category.ParentCategoryID = &parentCategoryID.Int64
	}

	return category, nil
}

func extractCategoriesFromRows(categoryRows *sql.Rows) ([]*domain.Category, error) {
	categories := make([]*domain.Category, 0)

	for categoryRows.Next() {
		category := &domain.Category{}
		var description sql.NullString
		var parentCategoryID sql.NullInt64

		err := categoryRows.Scan(
			&category.CategoryID,
			&category.Name,
			&description,
			&parentCategoryID,
			&category.CreatedAt,
		)

		if err != nil {
			log.Errorf("error while scanning category row: %v", err)
			continue
		}

		// Handle nullable fields
		if description.Valid {
			category.Description = &description.String
		}
		if parentCategoryID.Valid {
			category.ParentCategoryID = &parentCategoryID.Int64
		}

		categories = append(categories, category)
	}

	if err := categoryRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating category rows: %w", err)
	}

	return categories, nil
}
