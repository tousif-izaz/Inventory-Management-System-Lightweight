package domain

import "time"

// Supplier represents a vendor/supplier for purchases
type Supplier struct {
	SupplierID     int64      `json:"supplier_id"`
	Name           string     `json:"name"`
	ContactPerson  *string    `json:"contact_person,omitempty"`
	Email          *string    `json:"email,omitempty"`
	Phone          *string    `json:"phone,omitempty"`
	Address        *string    `json:"address,omitempty"`
	City           *string    `json:"city,omitempty"`
	Country        *string    `json:"country,omitempty"`
	TaxID          *string    `json:"tax_id,omitempty"`
	PaymentTerms   *string    `json:"payment_terms,omitempty"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
