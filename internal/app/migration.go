package app

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"log"
)

func RunDBMigration(sourceURL string, databaseName string, databaseInstance database.Driver) {
	migration, err := migrate.NewWithDatabaseInstance(sourceURL, databaseName, databaseInstance)
	if err != nil {
		log.Fatal("fail to create migration instance: ", err)
	}

	if err = migration.Up(); err != nil && err.Error() != "no change" {
		log.Fatal("fail to run migrate up: ", err)
	}

	log.Println("migration is successfully")
}
