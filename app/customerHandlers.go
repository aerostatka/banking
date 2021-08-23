package app

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/aerostatka/banking/service"
	"github.com/gorilla/mux"
)

type CustomerHandlers struct {
	service service.CustomerService
}

func (ch *CustomerHandlers) getAllCustomers(rw http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	customers, err := ch.service.GetAllCustomers(status)

	if err != nil {
		writeResponse(rw, err.Code, "application/json", err.AsMessage())
	} else {
		writeResponse(rw, http.StatusOK, "application/json", customers)
	}
}

func (ch *CustomerHandlers) GetCustomer(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["customerId"]

	customer, err := ch.service.GetCustomer(id)

	if err != nil {
		writeResponse(rw, err.Code, "application/json", err.AsMessage())
	} else {
		writeResponse(rw, http.StatusOK, "application/json", customer)
	}
}

func writeResponse(rw http.ResponseWriter, code int, contentType string, data interface{}) {
	rw.Header().Add("Content-Type", contentType)
	rw.WriteHeader(code)
	var err error
	switch contentType {
	case "application/xml":
		err = xml.NewEncoder(rw).Encode(data)
	case "application/json":
		err = json.NewEncoder(rw).Encode(data)
	default:
	}

	if err != nil {
		panic(err)
	}
}
