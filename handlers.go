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
	} else {
		err := account.Deposit(amount)
		if err != nil {
			fmt.Fprintf(w, "%v", err)

			return
		} else {
			web.Error(w, account.Statement())
			return
		}

	}
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
	} else {
		err := account.Withdraw(amount)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		} else {
			web.Error(w, account.Statement())
			return

		}

	}

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
	} else {
		err := accountA.Transfer(amount, accountB.Account)
		if err != nil {
			fmt.Fprintf(w, "%v", err)

			return
		} else {
			web.Error(w, accountA.Statement())
			return
		}
	}
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
	} else {
		json.NewEncoder(w).Encode(bank.Statement(account))

	}
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
