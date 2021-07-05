package models

import "github.com/shopspring/decimal"

// Customer ...
type Customer struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

// Account ...
type Account struct {
	Customer
	AccNumber string          `json:"accnumber"`
	Balance   decimal.Decimal `json:"balance"`
}
