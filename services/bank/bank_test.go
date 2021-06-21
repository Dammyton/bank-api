package bank

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func TestTransfer(t *testing.T) {
	accountA := Account{
		Customer: Customer{
			Name:    "John",
			Address: "Los Angeles, California",
			Phone:   "(213) 555 0147",
		},
		Number:  100000000001,
		Balance: decimal.NewFromInt(0),
	}

	accountB := Account{
		Customer: Customer{
			Name:    "Mark",
			Address: "Irvine, California",
			Phone:   "(949) 555 0198",
		},
		Number:  100000000002,
		Balance: decimal.NewFromInt(0),
	}

	accountA.Deposit(decimal.NewFromInt(100))
	err := accountA.Transfer(decimal.NewFromInt(50), &accountB)

	if accountA.Balance.Cmp(decimal.NewFromInt(50)) != 0 {
		t.Error("transfer from account A to account B is not working", err)
	}
}

func TestDeposit(t *testing.T) {
	account := Account{
		Customer: Customer{
			Name:    "John",
			Address: "Los Angeles, California",
			Phone:   "(213) 555 0147",
		},
		Number:  100000000001,
		Balance: decimal.NewFromInt(0),
	}

	account.Deposit(decimal.NewFromInt(500))

	if account.Balance.Cmp(decimal.NewFromInt(500)) != 0 {
		t.Error("balance is not being updated after a deposit")
	}
}

func TestWithdraw(t *testing.T) {
	account := Account{
		Customer: Customer{
			Name:    "John",
			Address: "Los Angeles, California",
			Phone:   "(213) 555 0147",
		},
		Number:  100000000001,
		Balance: decimal.NewFromInt(0),
	}
	account.Deposit(decimal.NewFromInt(1000))
	account.Withdraw(decimal.NewFromInt(100))

	if account.Balance.Cmp(decimal.NewFromInt(900)) != 0 {
		t.Error("balance is not being updated after withdraw")
	}

}

func TestStatement(t *testing.T) {
	account := Account{
		Customer: Customer{
			Name:    "John",
			Address: "Los Angeles, California",
			Phone:   "(213) 555 0147",
		},
		Number:  100000000001,
		Balance: decimal.NewFromInt(0),
	}

	account.Deposit(decimal.NewFromInt(100))
	statement := account.Statement()

	if statement != "100000000001 - John - 100" {
		t.Error("statement doesn't have the proper format")
	}
}

// Statement ...
func (a *Account) Statement() string {
	return fmt.Sprintf("%v - %v - %v", a.Number, a.Name, a.Balance)
}
