package usecase

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cukhoaimon/SimpleBank/utils"
	_ "github.com/lib/pq"
)

var testQuery *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../../../.")
	if err != nil {
		log.Fatal("Cannot load configuration file")
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("The open connection to database process was encoutered an error", err)
	}

	testQuery = New(testDB)

	os.Exit(m.Run())
}
