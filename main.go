package main

import (
	"log"
	"net/http"
	"os"

	"final1/pkg/api"
	"final1/pkg/db"
)

const (
	defaultPort = "7540"
	defaultDB   = "scheduler.db"
)

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	port := envOr("TODO_PORT", defaultPort)
	dbFile := envOr("TODO_DBFILE", defaultDB)

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	api.Init()

	addr := ":" + port
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
