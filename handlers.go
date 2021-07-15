package main

import (
	"bank-api/web"
	"encoding/json"
	"net/http"
	"time"
)

func (a *Application) createaccount(w http.ResponseWriter, req *http.Request) {
	name := web.QueryStrg(req, "name")
	code := web.QueryStr(req, "code")
	status := web.QueryStr(req, "status")
	balance := web.QueryStrToDecimal(req, "balance")

	id, err := a.account.CreateAccount(name, code, status, balance, time.Now())
	if err != nil {
		web.Error(w, err)
		return
	}

	json.NewEncoder(w).Encode(id)
}

func (a *Application) debit(w http.ResponseWriter, req *http.Request) {
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

	balance, err := a.account.GetAccountBalance(accNumber)
	if err != nil {
		web.Error(w, err)
		return
	}

	json.NewEncoder(w).Encode(balance)

}

func (a *Application) credit(w http.ResponseWriter, req *http.Request) {
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

	balance, err := a.account.GetAccountBalance(accNumber)
	if err != nil {
		web.Error(w, err)
		return
	}

	json.NewEncoder(w).Encode(balance)

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

	balance, err := a.account.GetAccountBalance(accNumber)
	if err != nil {
		web.Error(w, err)
		return
	}

	json.NewEncoder(w).Encode(balance)

}

func (a *Application) getaccountbal(w http.ResponseWriter, req *http.Request) {

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

func (a *Application) closeaccount(w http.ResponseWriter, req *http.Request) {

	accNumber := web.QueryStr(req, "number")

	if accNumber == 0 || accNumber < 10 {
		web.Error(w, "Invalid account number: wrong length")
		return
	}

	err := a.account.CloseAccount(accNumber, 2)
	if err != nil {
		web.Error(w, err)
		return
	}

}

func (a *Application) getclosedaccts(w http.ResponseWriter, req *http.Request) {

	accounts, err := a.account.GetAllAccounts(2)
	if err != nil {
		web.Error(w, err)
		return
	}

	json.NewEncoder(w).Encode(accounts)

}

func (a *Application) getaccounts(w http.ResponseWriter, req *http.Request) {

	accounts, err := a.account.GetAllAccounts(0)
	if err != nil {
		web.Error(w, err)
		return
	}

	json.NewEncoder(w).Encode(accounts)

}

func (a *Application) getactiveaccts(w http.ResponseWriter, req *http.Request) {

	accounts, err := a.account.GetAllAccounts(1)
	if err != nil {
		web.Error(w, err)
		return
	}

	json.NewEncoder(w).Encode(accounts)

}
