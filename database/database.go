package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var once sync.Once
var db *sql.DB
var err error

func Database() *sql.DB {
	once.Do(func() {
		connectionString := "postgresql://fitz:fitz@localhost:5432/fidus?sslmode=disable"
		db, err = sql.Open("postgres", connectionString)
	})

	if err != nil {
		log.Fatal(err)
	}

	return db
}
