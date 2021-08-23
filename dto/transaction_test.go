package dto

import (
	"net/http"
	"testing"
)

func Test_validation_error_for_transaction_type(t *testing.T) {
	req := NewTransactionRequest{
		TransactionType: "wrong type",
		Amount: 100,
	}

	appErr := req.Validate()

	if appErr.Message != "Transaction type should be deposit or withdrawal" {
		t.Error("Invalid validation of transaction type (message)")
	}

	if appErr.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid validation of transaction type (code)")
	}
}

func Test_validation_error_for_transaction_amount(t *testing.T) {
	req := NewTransactionRequest{
		TransactionType: "deposit",
		Amount: -100,
	}

	appErr := req.Validate()

	if appErr.Message != "Transaction should contain amount above 0" {
		t.Error("Invalid validation of transaction amount (message)")
	}

	if appErr.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid validation of transaction amount (code)")
	}
}
