package service

import (
	"time"

	"github.com/aerostatka/banking/domain"
	"github.com/aerostatka/banking/dto"
	"github.com/aerostatka/banking-lib/errs"
)

type AccountService interface {
	NewAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError)
	PerformTransaction(dto.NewTransactionRequest) (*dto.NewTransactionResponse, *errs.AppError)
}

type DefaultAccountService struct {
	repo domain.AccountRepository
}

func (s DefaultAccountService) NewAccount(request dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {
	if appErr := request.Validate(); appErr != nil {
		return nil, appErr
	}

	account := domain.NewAccount(request.CustomerId, request.AccountType, request.Amount)
	if newAccount, appError := s.repo.Save(account); appError != nil {
		return nil, appError
	} else {
		return newAccount.ToNewAccountResponseDto(), nil
	}
}

func (s DefaultAccountService) PerformTransaction(request dto.NewTransactionRequest) (*dto.NewTransactionResponse, *errs.AppError) {
	appErr := request.Validate()
	if appErr != nil {
		return nil, appErr
	}

	account, err := s.repo.ById(request.AccountId)
	if err != nil {
		return nil, err
	}

	if request.IsTypeWithdrawal() {
		if !account.CanWithdraw(request.Amount) {
			return nil, errs.NewValidationError("Amount is larger than account balance")
		}
	}

	if request.CustomerId != account.CustomerId {
		return nil, errs.NewValidationError("Invalid customer for the account")
	}

	transaction := domain.Transaction{
		AccountId: request.AccountId,
		Date:      time.Now().Format("2006-01-02 15:04:05"),
		Type:      request.TransactionType,
		Amount:    request.Amount,
	}
	newTransaction, newBalance, appError := s.repo.SaveTransaction(transaction)

	if appError != nil {
		return nil, appError
	}

	return newTransaction.ToNewTransactionResponseDto(newBalance), nil
}

func CreateAccountService(r domain.AccountRepository) DefaultAccountService {
	return DefaultAccountService{
		repo: r,
	}
}
