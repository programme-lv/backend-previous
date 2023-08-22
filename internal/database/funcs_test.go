package database

import (
	"log"
	"os"
	"sync"
	"testing"

	"github.com/programme-lv/backend/internal/database/testdb"
)

var (
	once       sync.Once
	dbProvider testdb.TestDBProvider
)

func TestMain(m *testing.M) {
	var err error
	dbProvider, err = testdb.NewPostgresTestcontainerProvider()
	if err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	dbProvider.Close()
	os.Exit(exitCode)
}

func TestCreateUser(t *testing.T) {
	db := dbProvider.GetTestDB()
	// asert db is not null
	if db == nil {
		t.Fatal("db is nil")
	}
}
