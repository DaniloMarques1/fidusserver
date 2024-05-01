package main

import (
	"fmt"
	"log"
	"os"

	"github.com/danilomarques1/fidusserver/database"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Not enable to load env variable %v\n", err)
	}

	db := database.Database()
	b, err := os.ReadFile("./database.sql")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = db.Exec(string(b)); err != nil {
		log.Fatal(err)
	}

	port := fmt.Sprintf(":%v", os.Getenv("PORT"))
	fidusServer := NewFidusServer(port)
	if err := fidusServer.Start(); err != nil {
		log.Fatal(err)
	}
}
