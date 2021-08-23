package domain

import (
	"github.com/aerostatka/banking/dto"
	"github.com/aerostatka/banking-lib/errs"
)

type Account struct {
	Id          string  `db:"account_id"`
	CustomerId  string  `db:"customer_id"`
	OpeningDate string  `db:"opening_date"`
	Type        string  `db:"account_type"`
	Amount      float64 `db:"amount"`
	Status      string  `db:"status"`
}

type Transaction struct {
	Id        string
	AccountId string
	Amount    float64
	Type      string
	Date      string
}

const dbTSLayout = "2006-01-02 15:04:05"

//go:generate mockgen -destination=../mock/domain/mockAccountRepository.go -package=domain github.com/aerostatka/banking/domain AccountRepository
type AccountRepository interface {
	Save(Account) (*Account, *errs.AppError)
	ById(string) (*Account, *errs.AppError)
	SaveTransaction(Transaction) (*Transaction, float64, *errs.AppError)
}

func (a Account) ToNewAccountResponseDto() *dto.NewAccountResponse {
	return &dto.NewAccountResponse{
		AccountId: a.Id,
	}
}

func (a Account) CanWithdraw(amount float64) bool {
	return (a.Amount - amount) >= 0
}

func (t Transaction) ToNewTransactionResponseDto(balance float64) *dto.NewTransactionResponse {
	return &dto.NewTransactionResponse{
		Id:              t.Id,
		AccountId:       t.AccountId,
		NewBalance:      balance,
		TransactionType: t.Type,
		TransactionDate: t.Date,
	}
}

func NewAccount(customerId, accountType string, amount float64) Account {
	return Account{
		CustomerId:  customerId,
		OpeningDate: dbTSLayout,
		Type:        accountType,
		Amount:      amount,
		Status:      "1",
	}
}
