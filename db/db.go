package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func ConnectDataBase() *sql.DB {
	connStr := "host=localhost port=3030 user=postgres password=abzal2005 dbname=go-expense-tracker sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Can not connect:", err)
	}

	fmt.Println("Connected to PostgresSQL")
	return db
}
