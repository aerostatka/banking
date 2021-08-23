package domain

import (
	"database/sql"
	"strconv"

	"github.com/aerostatka/banking/dto"
	"github.com/aerostatka/banking-lib/errs"
	"github.com/aerostatka/banking-lib/logger"
	"github.com/jmoiron/sqlx"
)

type AccountRepositoryDb struct {
	client *sqlx.DB
}

func (r AccountRepositoryDb) Save(a Account) (*Account, *errs.AppError) {
	sqlInsert := "INSERT INTO accounts(customer_id, opening_date, account_type, amount, status) VALUES(?, ?, ?, ?, ?);"

	result, err := r.client.Exec(sqlInsert, a.CustomerId, a.OpeningDate, a.Type, a.Amount, a.Status)
	if err != nil {
		logger.Error("Error while creating new account: " + err.Error())

		return nil, errs.NewInternalServerError("Unexpected error from database")
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting last insert id new account: " + err.Error())

		return nil, errs.NewInternalServerError("Unexpected error from database")
	}

	a.Id = strconv.FormatInt(id, 10)

	return &a, nil
}

func (r AccountRepositoryDb) ById(id string) (*Account, *errs.AppError) {
	var a Account
	findQuery := "select * from accounts where account_id = ?"

	err := r.client.Get(&a, findQuery, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Account is not found")
		}

		logger.Error("Error while scanning account " + err.Error())

		return nil, errs.NewInternalServerError("Unexpected database error")
	}

	return &a, nil
}

func (r AccountRepositoryDb) SaveTransaction(t Transaction) (*Transaction, float64, *errs.AppError) {
	sqlInsert := "INSERT INTO transactions(account_id, transaction_type, amount, transaction_date) VALUES(?, ?, ?, ?);"
	sqlSelect := "SELECT amount FROM accounts WHERE account_id = ?;"
	sqlUpdate := "UPDATE accounts SET amount = ? WHERE account_id = ?;"

	tx, err := r.client.Begin()

	if err != nil {
		logger.Error("Error while starting transaction: " + err.Error())

		return nil, 0, errs.NewInternalServerError("Unexpected error from database")
	}

	insertResult, err := r.client.Exec(sqlInsert, t.AccountId, t.Type, t.Amount, t.Date)
	if err != nil {
		tx.Rollback()
		logger.Error("Error while creating new transaction: " + err.Error())

		return nil, 0, errs.NewInternalServerError("Unexpected error from database")
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		tx.Rollback()
		logger.Error("Error while getting last insert id transaction: " + err.Error())

		return nil, 0, errs.NewInternalServerError("Unexpected error from database")
	}

	t.Id = strconv.FormatInt(id, 10)

	var balance float64
	err = r.client.Get(&balance, sqlSelect, t.AccountId)
	if err != nil {
		tx.Rollback()
		logger.Error("Error while getting old balance: " + err.Error())

		return nil, 0, errs.NewInternalServerError("Unexpected error from database")
	}

	var newBalance float64
	if t.Type == dto.TRANSACTION_DEPOSIT {
		newBalance = balance + t.Amount
	} else {
		newBalance = balance - t.Amount
	}

	_, err = r.client.Exec(sqlUpdate, newBalance, t.AccountId)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while creating new transaction: " + err.Error())

		return nil, 0, errs.NewInternalServerError("Unexpected error from database")
	}

	tx.Commit()

	return &t, newBalance, nil
}

func NewAccountRepositoryDb(dbClient *sqlx.DB) AccountRepositoryDb {
	return AccountRepositoryDb{client: dbClient}
}
