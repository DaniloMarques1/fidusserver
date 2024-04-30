package database

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var once sync.Once
var db *sql.DB
var err error

func Database() *sql.DB {
	once.Do(func() {
		connectionString := os.Getenv("DATABASE_URI")
		db, err = sql.Open("postgres", connectionString)
	})

	if err != nil {
		log.Fatal(err)
	}

	return db
}
