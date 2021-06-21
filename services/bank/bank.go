package bank

import (
	"bank-api/models"
	"bank-api/web"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

// Define a AccountModel type which wraps a sql.DB connection pool.
type AccountModel struct {
	DB *sql.DB
}

// Withdraw ...
func (m *AccountModel) Withdraw(w http.ResponseWriter, amount decimal.Decimal, accNumber string) error {
	var customer models.Account

	if amount.IsZero() {
		return errors.New("amount must be greater than zero")
	}

	// create the select sql query
	sqlStatement := `SELECT * FROM customers WHERE accnumber=$1 `

	// execute the sql statement
	row := m.DB.QueryRow(sqlStatement, accNumber)

	err := row.Scan(&customer.Name, &customer.Address, &customer.Phone, &customer.AccNumber, &customer.Balance)
	if err != nil {
		return errors.New("account number can't be found")
	} else {
		sqlStatement := `UPDATE customers SET balance=$2 WHERE accnumber=$1 RETURNING *`
		// execute the sql statement
		total := customer.Balance.Sub(amount)
		err = m.DB.QueryRow(sqlStatement, accNumber, total).Scan(&customer.Name, &customer.Address, &customer.Phone, &customer.AccNumber, &customer.Balance)
		if err != nil {
			return errors.New(err.Error())
		}

	}

	json.NewEncoder(w).Encode(customer)
	return nil
}

// Deposit ...
func (m *AccountModel) Deposit(w http.ResponseWriter, amount decimal.Decimal, accNumber string) error {
	var customer models.Account

	if amount.IsZero() {
		return errors.New("amount must be greater than zero")
	}

	// create the select sql query
	sqlStatement := `SELECT * FROM customers WHERE accnumber=$1 `

	// execute the sql statement
	row := m.DB.QueryRow(sqlStatement, accNumber)

	err := row.Scan(&customer.Name, &customer.Address, &customer.Phone, &customer.AccNumber, &customer.Balance)
	if err != nil {
		return errors.New("account number can't be found")
	} else {
		sqlStatement := `UPDATE customers SET balance=$2 WHERE accnumber=$1 RETURNING *`
		// execute the sql statement
		total := customer.Balance.Add(amount)
		err = m.DB.QueryRow(sqlStatement, accNumber, total).Scan(&customer.Name, &customer.Address, &customer.Phone, &customer.AccNumber, &customer.Balance)
		if err != nil {
			return errors.New(err.Error())
		}

	}

	json.NewEncoder(w).Encode(customer)
	return nil
}

// Transfer ...
func (m *AccountModel) Transfer(w http.ResponseWriter, amount decimal.Decimal, accA string, accB string) error {

	if amount.IsZero() {
		return errors.New("amount must be greater than zero")
	}

	m.Withdraw(w, amount, accA)
	m.Deposit(w, amount, accB)
	return nil
}

// Statement ...
func (m *AccountModel) Statement(w http.ResponseWriter, accNumber string) error {
	var customer models.Account

	// create the select sql query
	sqlStatement := `SELECT * FROM customers WHERE accnumber=$1`

	// execute the sql statement
	row := m.DB.QueryRow(sqlStatement, accNumber)

	err := row.Scan(&customer.Name, &customer.Address, &customer.Phone, &customer.AccNumber, &customer.Balance)
	if err != nil {
		web.Error(w, fmt.Sprintf("Account with number %v can't be found!", accNumber))
	} else {
		// send the response
		json.NewEncoder(w).Encode(customer)

	}
	return nil
}
