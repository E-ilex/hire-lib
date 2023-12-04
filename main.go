package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func dbConnection() (*sql.DB, error) {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", "data.sqlite3")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func main() {
	_, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Golang microservice connects to DB!")
}
