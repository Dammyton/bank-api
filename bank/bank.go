package bank

import (
	"errors"

	"github.com/shopspring/decimal"
)

// Customer ...
type Customer struct {
	Name    string
	Address string
	Phone   string
}

// Account ...
type Account struct {
	Customer
	Number  int
	Balance decimal.Decimal
}

// Withdraw ...
func (a *Account) Withdraw(amount decimal.Decimal) error {

	if amount.IsZero() {
		return errors.New("amount must be greater than zero")
	}
	a.Balance = a.Balance.Sub(amount)
	return nil
}

// Deposit ...
func (a *Account) Deposit(amount decimal.Decimal) error {

	if amount.IsZero() {
		return errors.New("amount must be greater than zero")
	}

	a.Balance = a.Balance.Add(amount)
	return nil
}

// Transfer function
func (a *Account) Transfer(amount decimal.Decimal, dest *Account) error {

	if amount.IsZero() {
		return errors.New("amount must be greater than zero")
	}
	if a.Balance.LessThan(amount) {
		return errors.New("the amount to transfer should be greater than the account balance")
	}

	a.Withdraw(amount)
	dest.Deposit(amount)
	return nil
}

// Bank ...
type Bank interface {
	Statement() string
}

// Statement ...
func Statement(b Bank) string {
	return b.Statement()
}
