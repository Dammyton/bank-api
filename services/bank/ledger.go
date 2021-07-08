package bank

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

// Define a AccountModel type which wraps a sql.DB connection pool.
type LedgerModel struct {
	db *sqlx.DB
}

type LegerTransactionType int

const (
	AllTransactions   LegerTransactionType = 0
	CreditTransaction LegerTransactionType = 1
	DebitTransaction  LegerTransactionType = 2
)

type Leger struct {
	ID              int64
	DateCreated     time.Time
	DateTransaction time.Time
	Account         int64
	Type            LegerTransactionType
	Amount          decimal.Decimal
}

func (a LedgerModel) AddTransaction(account int64, amount decimal.Decimal, dateTrx time.Time, tType LegerTransactionType, tx ...DBI) error {
	var db DBI = a.db
	if tx != nil {
		db = tx[0]
	}
	_ = db

	if tType == DebitTransaction {
		amount = amount.Mul(decimal.NewFromInt(-1))
	} else if tType == AllTransactions {
		tType = CreditTransaction
	}

	_, err := db.Exec(`
		insert into ledger
			(account, amount, datetransaction, type)
		values
			($1, $2, $3, $4)
	`, account, amount, dateTrx, tType)

	return err
}

func (a LedgerModel) SumTransactions(account int64, tType LegerTransactionType, tx ...DBI) (amount decimal.Decimal, err error) {
	var db DBI = a.db
	if tx != nil {
		db = tx[0]
	}
	_ = db

	sql := `
		select
			sum(amount) from ledger
		Where
			id = $1
	`
	if tType != AllTransactions {
		sql += `
			and type = $2
		`
	}

	row := db.QueryRowx(sql, account, tType)
	err = row.Scan(&amount)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return
}
