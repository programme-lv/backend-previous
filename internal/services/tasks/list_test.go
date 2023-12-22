package tasks

import (
	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/database/testdb"
	"os"
	"testing"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	provider, err := testdb.NewPostgresTestcontainerProvider()
	if err != nil {
		panic(err)
	}
	db = provider.GetTestDB()
	defer provider.Close()

	code := m.Run()
	os.Exit(code)
}

func TestListPublishedTasks(t *testing.T) {
	tasks, err := ListPublishedTasks(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) == 0 {
		t.Fatal("no tasks returned")
	}
}
