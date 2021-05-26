package main

import (
	"bank-api/bank"
	"bank-api/web"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

var accounts = map[float64]*CustomAccount{}

var (
	number, amount, dest float64
	err                  error
)

func deposit(w http.ResponseWriter, req *http.Request) {
	accNumber := web.QueryStr(req, "number")
	amounts := web.QueryStr(req, "amount")

	if accNumber == "" {
		web.Error(w, "Account number is missing!")
		return
	}

	if number, err = decimal.NewFromString(accNumber); err != nil {
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
