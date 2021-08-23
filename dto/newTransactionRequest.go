package dto

import (
	"strings"

	"github.com/aerostatka/banking-lib/errs"
)

const (
	TRANSACTION_WITHDRAWAL = "withdrawal"
	TRANSACTION_DEPOSIT    = "deposit"
)

type NewTransactionRequest struct {
	AccountId       string  `json:"account_id"`
	CustomerId      string  `json:"customer_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}

func (req NewTransactionRequest) Validate() *errs.AppError {
	if req.Amount < 0 {
		return errs.NewValidationError("Transaction should contain amount above 0")
	}

	if !req.IsTypeWithdrawal() && !req.IsTypeDeposit() {
		return errs.NewValidationError("Transaction type should be deposit or withdrawal")
	}

	return nil
}

func (req NewTransactionRequest) IsTypeWithdrawal() bool {
	transactionType := strings.ToLower(req.TransactionType)
	return transactionType == TRANSACTION_WITHDRAWAL
}

func (req NewTransactionRequest) IsTypeDeposit() bool {
	transactionType := strings.ToLower(req.TransactionType)
	return transactionType == TRANSACTION_DEPOSIT
}
