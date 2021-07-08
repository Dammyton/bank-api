package bank

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type DBI interface {
	sqlx.Queryer
	sqlx.Execer
}

// Define a AccountModel type which wraps a sqlx.DB connection pool.
type AccountModel struct {
	DB *sqlx.DB
}

type AccountStatus int

const (
	ActiveAccount AccountStatus = 0
	ClosedAccount AccountStatus = 1
)

type Account struct {
	ID   int64
	Code string
	Name string

	Status      AccountStatus
	DateCreated time.Time
}

func (a AccountModel) CreateAccount(name, code string, status int, balance decimal.Decimal, tx ...DBI) (id int64, err error) {
	var db DBI
	if tx != nil {
		db = tx[0]
	} else {
		tx, err := a.DB.Beginx()
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
	var db DBI = a.DB
	if tx != nil {
		db = tx[0]
	}
	_ = db

	l := LedgerModel{}
	err = l.AddTransaction(accountID, amount, time.Now(), CreditTransaction, db)

	return
}

func (a AccountModel) DebitAccount(accountID int64, amount decimal.Decimal, tx ...*sqlx.Tx) (err error) {
	var db DBI = a.DB
	if tx != nil {
		db = tx[0]
	}
	_ = db

	l := LedgerModel{}
	err = l.AddTransaction(accountID, amount, time.Now(), DebitTransaction, db)

	return
}

func (a AccountModel) GetAccountBalance(accountID int64, tx ...*sqlx.Tx) (balance decimal.Decimal, err error) {
	var db DBI = a.DB
	if tx != nil {
		db = tx[0]
	}
	_ = db

	l := LedgerModel{}
	balance, err = l.SumTransactions(accountID, AllTransactions, db)

	return
}

func (a AccountModel) CloseAccount(id int64, status AccountStatus, tx ...DBI) (err error) {
	var db DBI
	if tx != nil {
		db = tx[0]
	} else {
		tx, err := a.DB.Beginx()
		if err != nil {
			return err
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
	update account
	set status = $2
	where id = $1
	returning id
`, status, id)
	err = retv.Scan(&id)
	if err != nil {
		return err
	}

	return
}

func (a AccountModel) TransferAmount(accountID int64, destAccountID int64, amount decimal.Decimal) (err error) {
	if amount.IsZero() {
		return errors.New("amount must be greater than zero")
	}
	a.DebitAccount(accountID, amount)
	a.CreditAccount(destAccountID, amount)
	return
}

func (a AccountModel) UpdateAccountInfo(name, code string, status int, balance decimal.Decimal, tx ...DBI) (id int64, err error) {
	var db DBI
	if tx != nil {
		db = tx[0]
	} else {
		tx, err := a.DB.Beginx()
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
		update account
		set name = $2
		where id = $1
		returning id
	`, name, id)
	err = retv.Scan(&id)
	if err != nil {
		return 0, err
	}

	l := LedgerModel{}
	err = l.AddTransaction(id, balance, time.Now(), CreditTransaction, db)

	return
}

func (a AccountModel) GetAccount(accountID int64, tx ...*sqlx.Tx) (account decimal.Decimal, err error) {
	var db DBI = a.DB
	if tx != nil {
		db = tx[0]
	}
	_ = db

	sql := `
	select code, name, status from account
	where id = $1
`
	row := db.QueryRowx(sql, accountID)
	err = row.Scan(&account)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return
}

func (a AccountModel) GetAllAccounts(tx ...*sqlx.Tx) ([]Account, error) {
	var db DBI = a.DB
	if tx != nil {
		db = tx[0]
	}
	_ = db

	sql := `
	select code, name, status from account
	order by id
`
	rows, err := db.Queryx(sql)
	if err != nil {
		return []Account{}, err
	}

	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Code,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (a AccountModel) GetClosedAccounts(aType AccountStatus, tx ...*sqlx.Tx) ([]Account, error) {
	var db DBI = a.DB
	if tx != nil {
		db = tx[0]
	}
	_ = db

	sql := `
	select name,code,status from account
	where status = $1
	order by id
`
	rows, err := db.Queryx(sql, aType)
	if err != nil {
		return []Account{}, err
	}

	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Code,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
