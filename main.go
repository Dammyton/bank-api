package main

import (
	"bank-api/config"
	"bank-api/services/bank"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type application struct {
	account *bank.AccountModel
}

func main() {
	// create the postgres db connection
	db, err := config.SetUpDatabase()
	if err != nil {
		log.Fatal(err)
	}
	// close the db connection
	defer db.Close()

	app := &application{
		account: &bank.AccountModel{DB: db},
	}

	http.HandleFunc("/statement", app.statement)
	http.HandleFunc("/deposit", app.deposit)
	http.HandleFunc("/withdraw", app.withdraw)
	http.HandleFunc("/transfer", app.transfer)

	hostAddr := "localhost:8000"
	log.Println("running on ", hostAddr)
	log.Fatal(http.ListenAndServe(hostAddr, nil))
}
