package bank

import (
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

	if accountA.Balance != decimal.NewFromInt(50) && accountB.Balance != decimal.NewFromInt(50) {
		t.Error("transfer from account A to account B is not working", err)
	}
}
