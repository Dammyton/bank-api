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
	BankAccountAll    AccountStatus = 0
	BankAccountActive AccountStatus = 1
	BankAccountClosed AccountStatus = 2
)

type Account struct {
	ID   int64
	Code int64
	Name string

	Status      AccountStatus
	DateCreated time.Time
}

func (a AccountModel) CreateAccount(name string, code int64, status int64, balance decimal.Decimal, dateCreated time.Time, tx ...DBI) (id int64, err error) {
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

	err = db.QueryRowx(`
		insert into account
			(code, name, status,balance, date_created)
		values
			($1, $2, $3, $4, $5)
		returning id;
	`, code, name, status, balance, dateCreated).Scan(&id)
	if err != nil {
		return 0, err
	}

	l := LedgerModel{}
	err = l.AddTransaction(code, balance, time.Now(), CreditTransaction, db)

	return code, nil
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

func (a AccountModel) CloseAccount(accountID int64, status AccountStatus, tx ...DBI) (err error) {
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
	update 
		account
	set 
		status = $1
	where 
		code = $2
	returning code
`, status, accountID)
	err = retv.Scan(&status)
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

func (a AccountModel) GetAccount(accountID int64, tx ...*sqlx.Tx) (data Account, err error) {
	var db DBI = a.DB
	if tx != nil {
		db = tx[0]
	}
	_ = db

	sql := `
	select code, name, status from account
	where code = $1
`
	row := db.QueryRowx(sql, accountID)
	err = row.Scan(&data.Code, &data.Name, &data.Status)
	if err != nil {
		return Account{}, err
	}

	return
}

func (a AccountModel) GetAllAccounts(aType AccountStatus, tx ...*sqlx.Tx) ([]Account, error) {
	var db DBI = a.DB
	if tx != nil {
		db = tx[0]
	}
	_ = db

	var sql string

	if aType == 0 {

		sql = `
	select code, name,status from account
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

	sql = `
	select 
		code, name, status from account
	where 
		status = $1
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
