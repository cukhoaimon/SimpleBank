package main

import (
	"database/sql"
	"log"

	"github.com/cukhoaimon/SimpleBank/api"
	db "github.com/cukhoaimon/SimpleBank/db/sqlc"
	"github.com/cukhoaimon/SimpleBank/utils"
	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load configuration file")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("The open connection to database process was encountered an error", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("Cannot create server")
	}

	if err = server.Start(config.ServerAddress); err != nil {
		log.Fatal("Cannot start server")
	}
}
