package service

import (
	realdomain "github.com/aerostatka/banking/domain"
	"github.com/aerostatka/banking/dto"
	"github.com/aerostatka/banking-lib/errs"
	"github.com/aerostatka/banking/mock/domain"
	"github.com/golang/mock/gomock"
	"testing"
)

var mockRepo *domain.MockAccountRepository
var s AccountService

func setup(t *testing.T) func ()  {
	ctrl := gomock.NewController(t)
	mockRepo = domain.NewMockAccountRepository(ctrl)
	s = CreateAccountService(mockRepo)

	return func () {
		s = nil
		defer ctrl.Finish()
	}
}

func Test_should_return_a_validation_error_for_new_account(t *testing.T) {
	req := dto.NewAccountRequest{
		CustomerId: "100",
		AccountType: "saving",
		Amount: 0,
	}

	s := CreateAccountService(nil)

	_, err := s.NewAccount(req)

	if err == nil {
		t.Error("Failed while testing account validation")
	}
}

func Test_should_return_an_error_from_the_server_side(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	req := dto.NewAccountRequest{
		CustomerId: "100",
		AccountType: "savings",
		Amount: 10000,
	}

	account := realdomain.NewAccount(req.CustomerId, req.AccountType, req.Amount)

	mockRepo.EXPECT().Save(account).Return(nil, errs.NewInternalServerError("Unexpected error"))

	_, err := s.NewAccount(req)

	if err == nil {
		t.Error("Test failed to return unexpected error")
	}
}

func Test_should_return_successfully_created_account(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	req := dto.NewAccountRequest{
		CustomerId: "100",
		AccountType: "savings",
		Amount: 10000,
	}

	account := realdomain.NewAccount(req.CustomerId, req.AccountType, req.Amount)

	accountWithId := account
	accountWithId.Id = "3001"

	mockRepo.EXPECT().Save(account).Return(&accountWithId, nil)

	newAccount, err := s.NewAccount(req)

	if err != nil {
		t.Error("Test failed while creating account")
	}

	if newAccount.AccountId != accountWithId.Id {
		t.Error("Test failed to match account id")
	}
}