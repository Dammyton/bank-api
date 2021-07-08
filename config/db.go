package config

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "test1234"
	dbname   = "bankapi"
)

func SetUpDatabase() (*sqlx.DB, error) {
	// setup database connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	DB, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	fmt.Println("Database Successfully connected!")
	return DB, nil

}
