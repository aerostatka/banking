package dto

import (
	"strings"

	"github.com/aerostatka/banking-lib/errs"
)

type NewAccountRequest struct {
	CustomerId  string  `json:"customer_id"`
	AccountType string  `json:"account_type"`
	Amount      float64 `json:"amount"`
}

func (req NewAccountRequest) Validate() *errs.AppError {
	if req.Amount < 5000 {
		return errs.NewValidationError("Amount should be not less than 5000.00")
	}

	accountType := strings.ToLower(req.AccountType)

	if accountType != "savings" && accountType != "checking" {
		return errs.NewValidationError("Account type should be checking or savings")
	}

	return nil
}
