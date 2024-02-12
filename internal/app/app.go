package app

import (
	"database/sql"
	db "github.com/cukhoaimon/SimpleBank/internal/usecase/sqlc"
	"github.com/cukhoaimon/SimpleBank/pkg/grpc"
	"github.com/cukhoaimon/SimpleBank/utils"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/lib/pq"
	"log"
)

// Run will run simple bank app
func Run(config utils.Config) {
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("The open connection to database process was encountered an error", err)
	}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Fatal("Error when create postgres driver instance: ", err)
	}

	RunDBMigration(config.MigrationURL, config.PostgresDB, driver)

	store := db.NewStore(conn)

	go grpc.RunGatewayServer(store, config)
	grpc.Run(store, config)
	//http2.Run(store, config)
}
