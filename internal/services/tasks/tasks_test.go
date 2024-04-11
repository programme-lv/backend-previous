package tasks

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/testing/testdb"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	dbContainer, err := testdb.NewMigratedPostgresTestcontainer()
	if err != nil {
		panic(err)
	}
	db = dbContainer.GetConn()
	defer dbContainer.Close()

	code := m.Run()
	os.Exit(code)
}
