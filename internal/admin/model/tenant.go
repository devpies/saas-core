package model

import "github.com/go-playground/validator/v10"

var tenantValidator *validator.Validate

func init() {
	v := NewValidator()
	tenantValidator = v
}

// NewTenant represents a new tenant.
type NewTenant struct {
	Email    string `json:"email" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Company  string `json:"companyName" validate:"required"`
	Plan     string `json:"plan" validate:"required,oneof=basic premium"`
}

// Tenant represents a tenant.
type Tenant struct {
	ID       string `json:"ID"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Company  string `json:"companyName"`
	Plan     string `json:"plan"`
}

// Validate validates the NewTenant.
func (nt *NewTenant) Validate() error {
	return tenantValidator.Struct(nt)
}
