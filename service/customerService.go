package service

import (
	"github.com/aerostatka/banking/domain"
	"github.com/aerostatka/banking/dto"
	"github.com/aerostatka/banking-lib/errs"
)

//go:generate mockgen -destination=../mock/service/mockCustomerService.go -package=service github.com/aerostatka/banking/service CustomerService

type CustomerService interface {
	GetAllCustomers(string) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func (s DefaultCustomerService) GetAllCustomers(status string) ([]dto.CustomerResponse, *errs.AppError) {
	if status == "active" {
		status = "1"
	} else if status == "inactive" {
		status = "0"
	} else {
		status = ""
	}

	var dtos []dto.CustomerResponse
	customers, err := s.repo.FindAll(status)

	if err != nil {
		return dtos, err
	}

	for _, c := range customers {
		dtos = append(dtos, c.ToDto())
	}

	return dtos, nil
}

func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.ById(id)

	if err != nil {
		return nil, err
	}

	response := c.ToDto()

	return &response, nil
}

func CreateCustomerService(r domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repo: r}
}
