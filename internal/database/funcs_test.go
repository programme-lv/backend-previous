package database

import (
	"testing"

	"github.com/programme-lv/backend/internal/database/testdb"
)

func TestCreateUser(t *testing.T) {
	dbProvider, err := testdb.NewPostgresTestcontainerProvider()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { dbProvider.Close() })

	db := dbProvider.GetTestDB()
	// asert db is not null
	if db == nil {
		t.Fatal("db is nil")
	}
}
