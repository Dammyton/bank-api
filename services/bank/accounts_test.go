package bank

import (
	"fmt"
	"testing"
	"time"

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
	db, err := sqlx.Connect("psql_txdb", "account")
	if err != nil {
		t.Fatal(err)
	}

	var name string
	var code int
	var status int
	var balance decimal.Decimal
	var date_created time.Time

	testModel := &AccountModel{db}
	_, err = testModel.CreateAccount("Adeola", 1901234567, 1, decimal.NewFromInt(250), time.Now())
	if err != nil {
		t.Fatal(err)
	}

	sql := `
	select 
		code, name, status, balance, date_created
	from 
		account
	where 
		code = $1;
`
	row := testModel.DB.QueryRowx(sql, 1901234567)
	err = row.Scan(&code, &name, &status, &balance, &date_created)
	if err != nil {
		t.Fatal("create account is not working", err)
	}

	if code != 1901234567 && name != "Adeola" && status != 1 && balance != decimal.NewFromInt(250) {
		t.Fatal("create account is not working", err)
	}

}

func TestCreditAccount(t *testing.T) {
	db, err := sqlx.Connect("psql_txdb", "account")
	if err != nil {
		t.Fatal(err)
	}

	testModel := &AccountModel{db}

	err = testModel.CreditAccount(1098765432, decimal.NewFromInt(200))
	if err != nil {
		t.Fatal(err)
	}

	var i Leger

	sql := `
	select account, amount from ledger
	where account = $1;
`
	row := testModel.DB.QueryRowx(sql, 1098765432)
	err = row.Scan(&i.Account, &i.Amount)
	if err != nil {
		t.Fatal(err)
	}

	if i.Amount.Cmp(decimal.NewFromInt(200)) != 0 || i.Account != 1098765432 {
		t.Fatal("credit account is not working", err)
	}

}

func TestDebitAccount(t *testing.T) {
	db, err := sqlx.Connect("psql_txdb", "account")
	if err != nil {
		t.Fatal(err)
	}

	testModel := &AccountModel{db}

	err = testModel.DebitAccount(1098765432, decimal.NewFromInt(200))
	if err != nil {
		t.Fatal(err)
	}
	var i Leger

	sql := `
	select account, amount from ledger
	where account = $1;
`
	row := testModel.DB.QueryRowx(sql, 1098765432)
	err = row.Scan(&i.Account, &i.Amount)
	if err != nil {
		t.Fatal(err)
	}

	if i.Amount.Cmp(decimal.NewFromInt(-200)) != 0 && i.Account != 1098765432 {
		t.Fatal("credit account is not working", err)
	}

}

func TestGetAccountBalance(t *testing.T) {
	db, err := sqlx.Connect("psql_txdb", "account")
	if err != nil {
		t.Fatal(err)
	}
	testModel := &AccountModel{db}

	err = testModel.CreditAccount(1098765432, decimal.NewFromInt(900))
	if err != nil {
		t.Fatal(err)
	}

	err = testModel.DebitAccount(1098765432, decimal.NewFromInt(500))
	if err != nil {
		t.Fatal(err)
	}

	balance, err := testModel.GetAccountBalance(1098765432)
	if err != nil {
		t.Fatal(err)
	}

	if balance.Cmp(decimal.NewFromInt(400)) != 0 {
		t.Fatal("get account balance is not working", err, balance)
	}

}

func TestTransferAmount(t *testing.T) {
	db, err := sqlx.Connect("psql_txdb", "account")
	if err != nil {
		t.Fatal(err)
	}
	testModel := &AccountModel{db}

	accountA, err := testModel.CreateAccount("John", 1000000011, 1, decimal.NewFromInt(250), time.Now())
	if err != nil {
		t.Fatal("create account is not working", err)
	}
	accountB, err := testModel.CreateAccount("Janet", 1000000012, 1, decimal.NewFromInt(300), time.Now())
	if err != nil {
		t.Fatal("create account is not working", err)
	}

	err = testModel.TransferAmount(accountA, accountB, decimal.NewFromInt(200))
	if err != nil {
		t.Fatal(err)
	}

	var amountA, amountB decimal.Decimal

	sql := `
	select 
		 sum(amount)
	from 
		ledger
	where 
		account = $1;
`
	row := testModel.DB.QueryRowx(sql, 1000000011)
	err = row.Scan(&amountA)
	if err != nil {
		t.Fatal(err)
	}

	sql_ := `
	select 
		 sum(amount)
	from 
		ledger
	where 
		account = $1;
`
	row = testModel.DB.QueryRowx(sql_, 1000000012)
	err = row.Scan(&amountB)
	if err != nil {
		t.Fatal(err)
	}

	if amountA.Cmp(decimal.NewFromInt(50)) != 0 && amountB.Cmp(decimal.NewFromInt(500)) != 0 {
		t.Fatal("transfer amount from account A to B is not working", err)
	}
}

func TestGetAccount(t *testing.T) {
	db, err := sqlx.Connect("psql_txdb", "account")
	if err != nil {
		t.Fatal(err)
	}
	testModel := &AccountModel{db}

	_, err = testModel.CreateAccount("John", 2098765431, 1, decimal.NewFromInt(250), time.Now())
	if err != nil {
		t.Fatal(err)
	}

	account, err := testModel.GetAccount(2098765431)
	if err != nil {
		t.Fatal(err)
	}

	if account.Code != 2098765431 && account.Name != "John" && account.Status != 1 {
		t.Fatal("get account is not working", err)
	}
}

func TestCloseAccount(t *testing.T) {
	db, err := sqlx.Connect("psql_txdb", "account")
	if err != nil {
		t.Fatal(err)
	}
	testModel := &AccountModel{db}

	var status int

	_, err = testModel.CreateAccount("Adeola", 2098765431, 1, decimal.NewFromInt(250), time.Now())
	if err != nil {
		t.Fatal(err)
	}

	err = testModel.CloseAccount(2098765431, 2)
	if err != nil {
		t.Fatal(err)
	}

	sql := `
	select 
		 status 
	from 
		account
	where 
		code = $1;
`
	row := testModel.DB.QueryRowx(sql, 2098765431)
	err = row.Scan(&status)
	if err != nil {
		t.Fatal(err)
	}

	if status != 2 {
		t.Fatal("close account is not working", err)
	}

}

func TestGetAllAccounts(t *testing.T) {
	db, err := sqlx.Connect("psql_txdb", "account")
	if err != nil {
		t.Fatal(err)
	}
	testModel := &AccountModel{db}

	_, err = testModel.CreateAccount("John", 2098765431, 1, decimal.NewFromInt(250), time.Now())
	if err != nil {
		t.Fatal(err)
	}

	_, err = testModel.CreateAccount("Doe", 1098765431, 1, decimal.NewFromInt(250), time.Now())
	if err != nil {
		t.Fatal(err)
	}

	account, err := testModel.GetAllAccounts(0)
	if err != nil {
		t.Fatal(err)
	}

	if len(account) != 2 {
		t.Fatal("get all accounts is not working", err, account)
	}
}
