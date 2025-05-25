package db

import (
	"database/sql"
	"log"
	"os"

	"testing"

	"github.com/ofer-sin/Courses/BackendCourse/simplebank/util"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// Load configuration
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Set up the database connection and other necessary resources
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB)
	// Run the tests
	exitCode := m.Run()

	// Clean up resources
	// ...

	// Exit with the appropriate code
	os.Exit(exitCode)
}
