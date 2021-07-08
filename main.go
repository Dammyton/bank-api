package main

import (
	"bank-api/config"
	"bank-api/services/bank"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Application struct {
	account *bank.AccountModel
}

func main() {
	// create the postgres db connection
	app, db, err := setup()
	if err != nil {
		log.Fatal(err)
	}
	// close the db connection
	defer db.Close()

	http.HandleFunc("/statement", app.statement)
	http.HandleFunc("/deposit", app.deposit)
	http.HandleFunc("/withdraw", app.withdraw)
	http.HandleFunc("/transfer", app.transfer)

	hostAddr := "localhost:8000"
	log.Println("running on ", hostAddr)
	log.Fatal(http.ListenAndServe(hostAddr, nil))
}

func setup() (app *Application, db *sqlx.DB, err error) {
	db, err = config.SetUpDatabase()
	if err != nil {
		return
	}

	app = &Application{
		account: &bank.AccountModel{DB: db},
	}

	return
}
