package main

import (
	"database/sql"
	"log"

	"github.com/cukhoaimon/SimpleBank/api"
	db "github.com/cukhoaimon/SimpleBank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("The open connection to database process was encoutered an error", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(serverAddress); err != nil {
		log.Fatal("Cannot start server")
	}
}
