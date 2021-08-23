package app

import (
	"encoding/json"
	"net/http"

	"github.com/aerostatka/banking/dto"
	"github.com/aerostatka/banking/service"
	"github.com/gorilla/mux"
)

type AccountHandlers struct {
	service service.AccountService
}

func (ah *AccountHandlers) CreateAccount(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var request dto.NewAccountRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		writeResponse(rw, http.StatusBadRequest, "application/json", err.Error())
	} else {
		request.CustomerId = vars["customerId"]
		account, appErr := ah.service.NewAccount(request)

		if appErr != nil {
			writeResponse(rw, appErr.Code, "application/json", appErr.Message)
		} else {
			writeResponse(rw, http.StatusCreated, "application/json", account)
		}
	}
}

func (ah *AccountHandlers) PerformTransaction(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var request dto.NewTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		writeResponse(rw, http.StatusBadRequest, "application/json", err.Error())
	} else {
		request.CustomerId = vars["customerId"]
		request.AccountId = vars["accountId"]
		account, appErr := ah.service.PerformTransaction(request)

		if appErr != nil {
			writeResponse(rw, appErr.Code, "application/json", appErr.Message)
		} else {
			writeResponse(rw, http.StatusCreated, "application/json", account)
		}
	}
}
