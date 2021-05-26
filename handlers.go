package main

import (
	"bank-api/bank"
	"bank-api/web"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
		web.Error(w, "Amounts must be greater than zero!")
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
	amounts := web.QueryStr(req, "amount")

	if accNumber == "" {
		web.Error(w, "Account number is missing!")
		return
	}

	if number, err = strconv.ParseFloat(accNumber, 64); err != nil {
		web.Error(w, "Invalid account number!")
		return
	}

	if amount, err = strconv.ParseFloat(amounts, 64); err != nil {
		web.Error(w, "Invalid amount number!")
		return
	}

	account, ok := accounts[number]
	if !ok {
		web.Error(w, fmt.Sprintf("Account with number %v can't be found!", number))
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
	amounts := web.QueryStr(req, "amount")
	destAccNumber := web.QueryStr(req, "dest")

	if accNumber == "" {
		web.Error(w, "Account number is missing!")
		return
	}

	if number, err = strconv.ParseFloat(accNumber, 64); err != nil {
		web.Error(w, "Invalid account number!")
		return
	}

	if amount, err = strconv.ParseFloat(amounts, 64); err != nil {
		web.Error(w, "Invalid amount number!")
		return
	}

	if dest, err = strconv.ParseFloat(destAccNumber, 64); err != nil {
		web.Error(w, "Invalid account destination number!")
		return
	}

	accountA, ok := accounts[number]
	if !ok {
		web.Error(w, fmt.Sprintf("Account with number %v can't be found!", number))
		return
	}

	accountB, ok := accounts[dest]
	if !ok {
		web.Error(w, fmt.Sprintf("Account with number %v can't be found!", dest))
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

	if accNumber == "" {
		web.Error(w, "Account number is missing!")
		return
	}

	number, err := strconv.ParseFloat(accNumber, 64)
	if err != nil {
		web.Error(w, "Invalid account number!")
		return
	}

	account, ok := accounts[number]
	if !ok {
		web.Error(w, fmt.Sprintf("Account with number %v can't be found!", number))
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
