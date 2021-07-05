package bank

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type DBI interface {
	sqlx.Queryer
	sqlx.Execer
}

// Define a AccountModel type which wraps a sql.DB connection pool.
type AccountModel struct {
	db *sqlx.DB
}

type Account struct {
	ID   int64
	Code string
	Name string

	Status      int
	DateCreated time.Time
}

func (a AccountModel) CreateAccount(name, code string, status int, balance decimal.Decimal, tx ...DBI) (id int64, err error) {
	var db DBI
	if tx != nil {
		db = tx[0]
	} else {
		tx, err := a.db.Beginx()
		if err != nil {
			return 0, err
		}
		db = tx

		defer func() {
			if err == nil {
				tx.Commit()
			} else {
				tx.Rollback()
			}
		}()
	}

	retv := db.QueryRowx(`
		insert into account
			(code, name, status)
		values
			($1, $2, $3)
		returning id
	`, code, name, status)
	err = retv.Scan(&id)
	if err != nil {
		return 0, err
	}

	l := LedgerModel{}
	err = l.AddTransaction(id, balance, time.Now(), CreditTransaction, db)

	return
}

func (a AccountModel) CreditAccount(accountID int64, amount decimal.Decimal, tx ...*sqlx.Tx) (err error) {
	var db DBI = a.db
	if tx != nil {
		db = tx[0]
	}
	_ = db

	l := LedgerModel{}
	err = l.AddTransaction(accountID, amount, time.Now(), CreditTransaction, db)

	return
}

func (a AccountModel) DebitAccount(accountID int64, amount decimal.Decimal, tx ...*sqlx.Tx) (err error) {
	var db DBI = a.db
	if tx != nil {
		db = tx[0]
	}
	_ = db

	l := LedgerModel{}
	err = l.AddTransaction(accountID, amount, time.Now(), DebitTransaction, db)

	return
}

func (a AccountModel) GetAccountBalance(accountID int64, tx ...*sqlx.Tx) (balance decimal.Decimal, err error) {
	var db DBI = a.db
	if tx != nil {
		db = tx[0]
	}
	_ = db

	l := LedgerModel{}
	balance, err = l.SumTransactions(accountID, AllTransactions, db)

	return
}
