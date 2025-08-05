package db

import (
	_ "github.com/lib/pq"

	"database/sql"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
