package main

import (
	"bank-api/bank"
	"bank-api/web"
	"encoding/json"
	"fmt"
	"net/http"
)

var accounts = map[string]*CustomAccount{}

func deposit(w http.ResponseWriter, req *http.Request) {
	accNumber := web.QueryStr(req, "number")

	if len(accNumber) == 0 || len(accNumber) < 12 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	amount := web.QueryStrToDecimal(req, "amount")
	if amount.IsZero() {
		web.Error(w, "Amount must be greater than zero!")
		return
	}

	account, ok := accounts[accNumber]
	if !ok {
		web.Error(w, fmt.Sprintf("Account with number %s can't be found!", accNumber))
		return
	}

	err := account.Deposit(amount)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	web.Error(w, account.Statement())

}

func withdraw(w http.ResponseWriter, req *http.Request) {
	accNumber := web.QueryStr(req, "number")

	if len(accNumber) == 0 || len(accNumber) < 12 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	amount := web.QueryStrToDecimal(req, "amount")
	if amount.IsZero() {
		web.Error(w, "Amount must be greater than zero!")
		return
	}

	account, ok := accounts[accNumber]
	if !ok {
		web.Error(w, fmt.Sprintf("Account with number %v can't be found!", accNumber))
		return
	}

	err := account.Withdraw(amount)
	if err != nil {
		web.Error(w, err)
		// web.Error(w, "Hello!")
		// fmt.Fprintf(w, "%v", err)
		return
	}

	web.Response(w, account.Statement())

}

func transfer(w http.ResponseWriter, req *http.Request) {
	accNumber := web.QueryStr(req, "number")
	destAccNumber := web.QueryStr(req, "dest")

	if len(accNumber) == 0 || len(accNumber) < 12 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	amount := web.QueryStrToDecimal(req, "amount")
	if amount.IsZero() {
		web.Error(w, "Amount must be greater than zero!")
		return
	}

	if len(destAccNumber) == 0 || len(destAccNumber) < 12 {
		web.Error(w, "Invalid account number: wrong destination!")
		return
	}

	accountA, ok := accounts[accNumber]
	if !ok {
		web.Error(w, fmt.Sprintf("Account with number %v can't be found!", accNumber))
		return
	}

	accountB, ok := accounts[destAccNumber]
	if !ok {
		web.Error(w, fmt.Sprintf("Account with number %v can't be found!", destAccNumber))
		return
	}

	err := accountA.Transfer(amount, accountB.Account)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	web.Error(w, accountA.Statement())

}

func statement(w http.ResponseWriter, req *http.Request) {
	accNumber := web.QueryStr(req, "number")

	if len(accNumber) == 0 || len(accNumber) < 12 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	amount := web.QueryStrToDecimal(req, "amount")
	if amount.IsZero() {
		web.Error(w, "Amount must be greater than zero!")
		return
	}

	account, ok := accounts[accNumber]
	if !ok {
		web.Error(w, fmt.Sprintf("Account with number %v can't be found!", accNumber))
		return
	}
	json.NewEncoder(w).Encode(bank.Statement(account))

}

// CustomAccount ...
type CustomAccount struct {
	*bank.Account
}

// Statement ...
func (c *CustomAccount) Statement() string {
	json, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}

	return string(json)
}
