package service

import (
	"errors"
	"fmt"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service/dto"
)

type ICategoryService interface {
	Add(categoryCreate *dto.CategoryCreate) error
	GetAllCategories() ([]*domain.Category, error)
	GetCategoryByID(categoryID int64) (*domain.Category, error)
	GetChildCategories(parentCategoryID int64) ([]*domain.Category, error)
	UpdateCategoryById(updatedCategory *dto.CategoryCreate, categoryId int64) error
	DeleteById(categoryId int64) error
}

type CategoryService struct {
	categoryRepository repository.ICategoryRepository
}

func NewCategoryService(categoryRepository repository.ICategoryRepository) ICategoryService {
	return &CategoryService{categoryRepository}
}

func (service *CategoryService) Add(categoryCreate *dto.CategoryCreate) error {
	err := validateCategoryCreate(categoryCreate)
	if err != nil {
		return err
	}

	category := categoryCreateToCategory(categoryCreate)
	return service.categoryRepository.AddCategory(category)
}

func (service *CategoryService) GetAllCategories() ([]*domain.Category, error) {
	return service.categoryRepository.GetAllCategories()
}

func (service *CategoryService) GetCategoryByID(categoryID int64) (*domain.Category, error) {
	if categoryID <= 0 {
		return nil, errors.New("invalid category ID")
	}
	return service.categoryRepository.GetCategoryByID(categoryID)
}

func (service *CategoryService) GetChildCategories(parentCategoryID int64) ([]*domain.Category, error) {
	if parentCategoryID <= 0 {
		return nil, errors.New("invalid parent category ID")
	}
	return service.categoryRepository.GetChildCategories(parentCategoryID)
}

func (service *CategoryService) UpdateCategoryById(updatedCategory *dto.CategoryCreate, categoryId int64) error {
	err := service.categoryRepository.CheckCategoryExistence(categoryId)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	err = validateCategoryCreate(updatedCategory)
	if err != nil {
		return err
	}

	category := categoryCreateToCategory(updatedCategory)
	return service.categoryRepository.UpdateCategoryById(category, categoryId)
}

func (service *CategoryService) DeleteById(categoryId int64) error {
	err := service.categoryRepository.CheckCategoryExistence(categoryId)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	return service.categoryRepository.DeleteCategoryById(categoryId)
}

func validateCategoryCreate(categoryCreate *dto.CategoryCreate) error {
	if categoryCreate.Name == "" {
		return errors.New("name can't be empty")
	}
	if len(categoryCreate.Name) > 255 {
		return errors.New("name can't exceed 255 characters")
	}

	return nil
}

func categoryCreateToCategory(categoryCreate *dto.CategoryCreate) *domain.Category {
	return &domain.Category{
		Name:             categoryCreate.Name,
		Description:      categoryCreate.Description,
		ParentCategoryID: categoryCreate.ParentCategoryID,
	}
}
