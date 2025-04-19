package main

import (
	"backend/internal/db"
	"backend/internal/env"
	"backend/internal/store"
	"log"
)

func main() {

	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 30, 30, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store)
}
