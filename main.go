package main

import (
	"log"
	"os"

	"github.com/danilomarques1/fidusserver/database"
)

func main() {
	db := database.Database()
	b, err := os.ReadFile("./database.sql")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = db.Exec(string(b)); err != nil {
		log.Fatal(err)
	}

	port := ":8080"
	fidusServer := NewFidusServer(port)
	if err := fidusServer.Start(); err != nil {
		log.Fatal(err)
	}
}
