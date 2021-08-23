package domain

import (
	"github.com/aerostatka/banking/dto"
	"github.com/aerostatka/banking-lib/errs"
)

type Customer struct {
	Id      string `json:"uuid" db:"customer_id"`
	Name    string `json:"name"`
	City    string `json:"city"`
	ZipCode string `json:"postalCode"`
	DOB     string `json:"dob" db:"date_of_birth"`
	Status  bool   `json:"state"`
}

type CustomerRepository interface {
	FindAll(string) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
}

func (c Customer) statusAsText() string {
	status := "inactive"
	if c.Status {
		status = "active"
	}

	return status
}

func (c Customer) ToDto() dto.CustomerResponse {
	return dto.CustomerResponse{
		Id:      c.Id,
		Name:    c.Name,
		City:    c.City,
		ZipCode: c.ZipCode,
		DOB:     c.DOB,
		Status:  c.statusAsText(),
	}
}
