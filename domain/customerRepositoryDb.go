package domain

import (
	"database/sql"

	"github.com/aerostatka/banking-lib/errs"
	"github.com/aerostatka/banking-lib/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

func (r CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {
	var err error
	customers := []Customer{}

	findQuery := "select customer_id, name, city, zipcode, date_of_birth, status from customers"
	if status != "" {
		findQuery += " where status = ?"
		err = r.client.Select(&customers, findQuery, status)
	} else {
		err = r.client.Select(&customers, findQuery)
	}

	if err != nil {
		logger.Error("Error while querying customers table: " + err.Error())

		return nil, errs.NewInternalServerError("Unexpected database error")
	}

	return customers, nil
}

func (r CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	var c Customer
	customerQuery := "select customer_id, name, city, zipcode, date_of_birth, status from customers where customer_id = ?"

	err := r.client.Get(&c, customerQuery, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer is not found")
		}

		logger.Error("Error while scanning customer " + err.Error())

		return nil, errs.NewInternalServerError("Unexpected database error")
	}

	return &c, nil
}

func NewCustomerRepositoryDb(dbClient *sqlx.DB) CustomerRepositoryDb {
	return CustomerRepositoryDb{client: dbClient}
}
