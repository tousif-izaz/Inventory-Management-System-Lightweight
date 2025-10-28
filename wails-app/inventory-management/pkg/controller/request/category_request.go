package request

import "inventory-management/pkg/service/dto"

// CreateCategoryRequest represents a request to create a new category
type CreateCategoryRequest struct {
	Name             string  `json:"name" validate:"required,min=1,max=255"`
	Description      *string `json:"description,omitempty"`
	ParentCategoryID *int64  `json:"parent_category_id,omitempty"`
}

// ToDTO converts request to DTO
func (r *CreateCategoryRequest) ToDTO() *dto.CategoryCreate {
	return &dto.CategoryCreate{
		Name:             r.Name,
		Description:      r.Description,
		ParentCategoryID: r.ParentCategoryID,
	}
}

// UpdateCategoryRequest represents a request to update a category
type UpdateCategoryRequest struct {
	Name             string  `json:"name" validate:"required,min=1,max=255"`
	Description      *string `json:"description,omitempty"`
	ParentCategoryID *int64  `json:"parent_category_id,omitempty"`
}

// ToDTO converts request to DTO
func (r *UpdateCategoryRequest) ToDTO() *dto.CategoryCreate {
	return &dto.CategoryCreate{
		Name:             r.Name,
		Description:      r.Description,
		ParentCategoryID: r.ParentCategoryID,
	}
}
