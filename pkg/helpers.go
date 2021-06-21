package pkg

import (
	"github.com/shopspring/decimal"
)

func StringToDecimal(val string) decimal.Decimal {
	retv, err := decimal.NewFromString(val)
	if err != nil {
		return decimal.NewFromInt(0)
	}

	return retv
}
