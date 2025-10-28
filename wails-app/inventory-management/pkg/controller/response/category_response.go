package response

import (
	"inventory-management/pkg/domain"
	"time"
)

// CategoryResponse represents a category in API responses
type CategoryResponse struct {
	CategoryID       int64     `json:"category_id"`
	Name             string    `json:"name"`
	Description      *string   `json:"description,omitempty"`
	ParentCategoryID *int64    `json:"parent_category_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

func toCategoryResponse(category *domain.Category) *CategoryResponse {
	return &CategoryResponse{
		CategoryID:       category.CategoryID,
		Name:             category.Name,
		Description:      category.Description,
		ParentCategoryID: category.ParentCategoryID,
		CreatedAt:        category.CreatedAt,
	}
}

func ToCategoryResponseList(categories []*domain.Category) []*CategoryResponse {
	responses := make([]*CategoryResponse, 0)
	for _, category := range categories {
		responses = append(responses, toCategoryResponse(category))
	}
	return responses
}
