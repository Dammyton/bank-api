package main

import (
	"bank-api/web"
	"net/http"
)

func (a *Application) withdraw(w http.ResponseWriter, req *http.Request) {
	amount := web.QueryStrToDecimal(req, "amount")
	accNumber := web.QueryStr(req, "number")

	if len(accNumber) == 0 || len(accNumber) < 10 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	err := a.account.Withdraw(w, amount, accNumber)
	if err != nil {
		web.Error(w, err)
		return
	}

}

func (a *Application) deposit(w http.ResponseWriter, req *http.Request) {
	amount := web.QueryStrToDecimal(req, "amount")
	accNumber := web.QueryStr(req, "number")

	if len(accNumber) < 10 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	err := a.account.Deposit(w, amount, accNumber)
	if err != nil {
		web.Error(w, err)
		return
	}

}

func (a *Application) transfer(w http.ResponseWriter, req *http.Request) {
	accNumber := web.QueryStr(req, "number")
	destAccNumber := web.QueryStr(req, "dest")
	amount := web.QueryStrToDecimal(req, "amount")

	if len(accNumber) == 0 || len(accNumber) < 10 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	if len(destAccNumber) == 0 || len(destAccNumber) < 10 {
		web.Error(w, "Invalid account number: wrong destination!")
		return
	}

	err := a.account.Transfer(w, amount, accNumber, destAccNumber)
	if err != nil {
		web.Error(w, err)
		return
	}

}

func (a *Application) statement(w http.ResponseWriter, req *http.Request) {

	accNumber := web.QueryStr(req, "number")

	if len(accNumber) == 0 || len(accNumber) < 10 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	err := a.account.Statement(w, accNumber)
	if err != nil {
		web.Error(w, err)
		return
	}

}
