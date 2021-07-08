package bank

import (
	"bank-api/models"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "test1234"
	dbname   = "bankapi"
)

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	txdb.Register("psql_txdb", "postgres", psqlInfo)
}

func TestCreateAccount(t *testing.T) {

}

func TestCreditAccount(t *testing.T) {
	db, err := sqlx.Connect("psql_txdb", "account")
	if err != nil {
		fmt.Println(err)
	}
	testModel := &AccountModel{db}
	account := &models.Account{
		Customer: models.Customer{
			Name:    "John",
			Address: "Los Angeles, California",
			Phone:   "(213) 555 0147",
		},

		AccNumber: 1000000005,
		Balance:   decimal.NewFromInt(300),
	}

	err = testModel.CreditAccount(account.AccNumber, decimal.NewFromInt(250))
	if err != nil {
		t.Error("credit account is not working", err)
	}
}
