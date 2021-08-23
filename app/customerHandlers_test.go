package app

import (
	"github.com/aerostatka/banking/dto"
	"github.com/aerostatka/banking-lib/errs"
	"github.com/aerostatka/banking/mock/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router *mux.Router
var ch CustomerHandlers
var mockService *service.MockCustomerService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = service.NewMockCustomerService(ctrl)
	ch = CustomerHandlers{service: mockService}
	router = mux.NewRouter()
	router.HandleFunc("/customers", ch.getAllCustomers)

	return func() {
		router = nil
		defer ctrl.Finish()
	}
}

func TestCustomerHandlers_GetCustomer_return_ok_status(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	dummyCustomers := []dto.CustomerResponse{
		{Id: "uuid1111", Name: "Alex Johnson", City: "New York", ZipCode: "10050", DOB: "05/25/1985", Status: "active"},
		{Id: "uuid1112", Name: "Dakota Fanning", City: "Jersey City", ZipCode: "07306", DOB: "12/31/1990", Status: "active"},
		{Id: "uuid1113", Name: "Alexis Castle", City: "Bensalem", ZipCode: "19020", DOB: "01/09/1991", Status: "active"},
	}
	mockService.EXPECT().GetAllCustomers("").Return(dummyCustomers, nil)

	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Error("Failed to receive correct status code")
	}
}
func TestCustomerHandlers_GetCustomer_return_unhandled_exception(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	mockService.EXPECT().GetAllCustomers("").Return(nil, errs.NewInternalServerError("Test error"))
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed to receive correct status code")
	}
}