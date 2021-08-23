package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aerostatka/banking/domain"
	"github.com/aerostatka/banking-lib/logger"
	"github.com/aerostatka/banking/service"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func sanityCheck() {
	if os.Getenv("SERVER_ADDRESS") == "" ||
		os.Getenv("SERVER_PORT") == "" {
		logger.Fatal("Server environment variables are not defined.")
	}

	if os.Getenv("DB_HOST") == "" ||
		os.Getenv("DB_PORT") == "" ||
		os.Getenv("DB_USER") == "" ||
		os.Getenv("DB_USER") == "" ||
		os.Getenv("DB_NAME") == "" {
		logger.Fatal("DB environment variables are not defined.")
	}
}

func Start() {
	sanityCheck()

	router := mux.NewRouter()
	dbClient := getDbClient()
	customerRepositoryDb := domain.NewCustomerRepositoryDb(dbClient)
	accountRepositoryDb := domain.NewAccountRepositoryDb(dbClient)
	authRepository := domain.NewAuthRepository()
	ch := CustomerHandlers{service: service.CreateCustomerService(customerRepositoryDb)}
	ah := AccountHandlers{service: service.CreateAccountService(accountRepositoryDb)}
	am := AuthMiddleware{repo: authRepository}
	router.
		HandleFunc("/customers", ch.getAllCustomers).
		Methods(http.MethodGet).
		Name("GetAllCustomers")
	router.
		HandleFunc("/customers/{customerId:[0-9]+}", ch.GetCustomer).
		Methods(http.MethodGet).
		Name("GetCustomer")
	router.
		HandleFunc("/customers/{customerId:[0-9]+}/account", ah.CreateAccount).
		Methods(http.MethodPost).
		Name("NewAccount")
	router.
		HandleFunc("/customers/{customerId:[0-9]+}/account/{accountId:[0-9]+}", ah.PerformTransaction).
		Methods(http.MethodPost).
		Name("NewTransaction")

	router.Use(am.authorizationHandler())

	host := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), router))
}

func getDbClient() *sqlx.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	client, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		panic(err)
	}

	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client
}
