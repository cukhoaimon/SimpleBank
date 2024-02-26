package app

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/rs/zerolog/log"
)

func RunDBMigration(sourceURL string, databaseName string, databaseInstance database.Driver) {
	migration, err := migrate.NewWithDatabaseInstance(sourceURL, databaseName, databaseInstance)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create migration instance: ")
	}

	if err = migration.Up(); err != nil && err.Error() != "no change" {
		log.Fatal().Err(err).Msg("fail to run migrate up: ")
	}

	log.Info().Msg("migration is successfully")
}
