package bank

import (
	"bank-api/models"
	"database/sql"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-txdb"
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

func TestWithdraw(t *testing.T) {
	db, err := sql.Open("psql_txdb", "customers")
	if err != nil {
		fmt.Println(err)
	}
	testModel := &AccountModel{db}

	w := httptest.NewRecorder()
	account := &models.Account{
		Customer: models.Customer{
			Name:    "John",
			Address: "Los Angeles, California",
			Phone:   "(213) 555 0147",
		},

		AccNumber: "1000000005",
		Balance:   decimal.NewFromInt(300),
	}

	sqlStatement := `
	INSERT INTO customers (name, address, phone, accnumber, balance)
	VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, account.Name, account.Address, account.Phone, account.AccNumber, account.Balance)
	if err != nil {
		t.Error(err)
	}

	err = testModel.Withdraw(w, decimal.NewFromInt(250), account.AccNumber)
	if err != nil {
		t.Error("withdraw from account is not working", err)
	}

}

func TestDeposit(t *testing.T) {
	db, err := sql.Open("psql_txdb", "customers")
	if err != nil {
		fmt.Println(err)
	}
	testModel := &AccountModel{db}

	w := httptest.NewRecorder()
	account := &models.Account{
		Customer: models.Customer{
			Name:    "John",
			Address: "Los Angeles, California",
			Phone:   "(213) 555 0147",
		},

		AccNumber: "1000000005",
		Balance:   decimal.NewFromInt(300),
	}

	sqlStatement := `
	INSERT INTO customers (name, address, phone, accnumber, balance)
	VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, account.Name, account.Address, account.Phone, account.AccNumber, account.Balance)
	if err != nil {
		t.Error(err)
	}

	err = testModel.Deposit(w, decimal.NewFromInt(250), account.AccNumber)
	if err != nil {
		t.Error("deposit into account is not working", err)
	}

}

func TestTransfer(t *testing.T) {
	db, err := sql.Open("psql_txdb", "customers")
	if err != nil {
		fmt.Println(err)
	}
	testModel := &AccountModel{db}

	w := httptest.NewRecorder()
	accountA := &models.Account{
		Customer: models.Customer{
			Name:    "John",
			Address: "Los Angeles, California",
			Phone:   "(213) 555 0147",
		},

		AccNumber: "1000000005",
		Balance:   decimal.NewFromInt(300),
	}

	sqlStatementA := `
	INSERT INTO customers (name, address, phone, accnumber, balance)
	VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatementA, accountA.Name, accountA.Address, accountA.Phone, accountA.AccNumber, accountA.Balance)
	if err != nil {
		t.Error(err)
	}

	accountB := &models.Account{
		Customer: models.Customer{
			Name:    "John",
			Address: "Los Angeles, California",
			Phone:   "(213) 555 0147",
		},

		AccNumber: "1000000005",
		Balance:   decimal.NewFromInt(500),
	}

	sqlStatementB := `
	INSERT INTO customers (name, address, phone, accnumber, balance)
	VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatementB, accountB.Name, accountB.Address, accountB.Phone, accountB.AccNumber, accountB.Balance)
	if err != nil {
		t.Error(err)
	}

	err = testModel.Transfer(w, decimal.NewFromInt(250), accountA.AccNumber, accountB.AccNumber)
	if err != nil {
		t.Error("transfer from account A to account B is not working", err)
	}

}

func TestStatement(t *testing.T) {
	db, err := sql.Open("psql_txdb", "customers")
	if err != nil {
		fmt.Println(err)
	}
	testModel := &AccountModel{db}

	w := httptest.NewRecorder()
	account := &models.Account{
		Customer: models.Customer{
			Name:    "John",
			Address: "Los Angeles, California",
			Phone:   "(213) 555 0147",
		},

		AccNumber: "1000000005",
		Balance:   decimal.NewFromInt(300),
	}

	sqlStatement := `
	INSERT INTO customers (name, address, phone, accnumber, balance)
	VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(sqlStatement, account.Name, account.Address, account.Phone, account.AccNumber, account.Balance)
	if err != nil {
		t.Error(err)
	}

	err = testModel.Statement(w, account.AccNumber)
	if err != nil {
		t.Error("statement of account is not working", err)
	}

}
