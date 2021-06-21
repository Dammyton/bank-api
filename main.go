package main

import (
	"log"
	"net/http"

	"bank-api/bank"
)

func main() {
	accounts["100000000001"] = &CustomAccount{
		Account: &bank.Account{
			Customer: bank.Customer{
				Name:    "John",
				Address: "Los Angeles, California",
				Phone:   "(213) 555 0147",
			},
			Number: 100000000001,
		},
	}

	accounts["100000000002"] = &CustomAccount{
		Account: &bank.Account{
			Customer: bank.Customer{
				Name:    "Mark",
				Address: "Irvine, California",
				Phone:   "(949) 555 0198",
			},
			Number: 100000000002,
		},
	}

	http.HandleFunc("/statement", statement)
	http.HandleFunc("/deposit", deposit)
	http.HandleFunc("/withdraw", withdraw)
	http.HandleFunc("/transfer", transfer)

	hostAddr := "localhost:8000"
	log.Println("running on ", hostAddr)
	log.Fatal(http.ListenAndServe(hostAddr, nil))
}
