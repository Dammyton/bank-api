package main

import (
	"bank-api/web"
	"encoding/json"
	"net/http"
)

func (a *Application) withdraw(w http.ResponseWriter, req *http.Request) {
	amount := web.QueryStrToDecimal(req, "amount")
	accNumber := web.QueryStr(req, "number")

	if accNumber == 0 || accNumber < 10 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	err := a.account.DebitAccount(accNumber, amount)
	if err != nil {
		web.Error(w, err)
		return
	}

}

func (a *Application) deposit(w http.ResponseWriter, req *http.Request) {
	amount := web.QueryStrToDecimal(req, "amount")
	accNumber := web.QueryStr(req, "number")

	if accNumber < 10 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	err := a.account.CreditAccount(accNumber, amount)
	if err != nil {
		web.Error(w, err)
		return
	}

}

func (a *Application) transfer(w http.ResponseWriter, req *http.Request) {
	accNumber := web.QueryStr(req, "number")
	destAccNumber := web.QueryStr(req, "dest")
	amount := web.QueryStrToDecimal(req, "amount")

	if accNumber == 0 || accNumber < 10 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	if destAccNumber == 0 || destAccNumber < 10 {
		web.Error(w, "Invalid account number: wrong destination!")
		return
	}

	err := a.account.TransferAmount(accNumber, destAccNumber, amount)
	if err != nil {
		web.Error(w, err)
		return
	}

}

func (a *Application) statement(w http.ResponseWriter, req *http.Request) {

	accNumber := web.QueryStr(req, "number")

	if accNumber == 0 || accNumber < 10 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	balance, err := a.account.GetAccountBalance(accNumber)
	if err != nil {
		web.Error(w, err)
		return
	}

	json.NewEncoder(w).Encode(balance)

}
