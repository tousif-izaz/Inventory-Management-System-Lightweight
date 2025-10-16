package request

import "ims-intro/pkg/service/dto"

// CreateSupplierRequest represents a request to create a supplier
type CreateSupplierRequest struct {
	Name          string  `json:"name" validate:"required,min=1,max=255"`
	ContactPerson *string `json:"contact_person,omitempty"`
	Email         *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone         *string `json:"phone,omitempty"`
	Address       *string `json:"address,omitempty"`
	City          *string `json:"city,omitempty"`
	Country       *string `json:"country,omitempty"`
	TaxID         *string `json:"tax_id,omitempty"`
	PaymentTerms  *string `json:"payment_terms,omitempty"`
}

// ToDTO converts request to DTO
func (r *CreateSupplierRequest) ToDTO() *dto.SupplierCreate {
	return &dto.SupplierCreate{
		Name:          r.Name,
		ContactPerson: r.ContactPerson,
		Email:         r.Email,
		Phone:         r.Phone,
		Address:       r.Address,
		City:          r.City,
		Country:       r.Country,
		TaxID:         r.TaxID,
		PaymentTerms:  r.PaymentTerms,
	}
}

// UpdateSupplierRequest represents a request to update a supplier
type UpdateSupplierRequest struct {
	Name          string  `json:"name" validate:"required,min=1,max=255"`
	ContactPerson *string `json:"contact_person,omitempty"`
	Email         *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone         *string `json:"phone,omitempty"`
	Address       *string `json:"address,omitempty"`
	City          *string `json:"city,omitempty"`
	Country       *string `json:"country,omitempty"`
	TaxID         *string `json:"tax_id,omitempty"`
	PaymentTerms  *string `json:"payment_terms,omitempty"`
}

// ToDTO converts request to DTO
func (r *UpdateSupplierRequest) ToDTO() *dto.SupplierCreate {
	return &dto.SupplierCreate{
		Name:          r.Name,
		ContactPerson: r.ContactPerson,
		Email:         r.Email,
		Phone:         r.Phone,
		Address:       r.Address,
		City:          r.City,
		Country:       r.Country,
		TaxID:         r.TaxID,
		PaymentTerms:  r.PaymentTerms,
	}
}
