package langs

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/database/testdb"
	"github.com/stretchr/testify/assert"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	provider, err := testdb.NewPostgresTestcontainer()
	if err != nil {
		panic(err)
	}
	db = provider.GetConn()
	defer provider.Close()

	code := m.Run()
	os.Exit(code)
}

func TestListProgrammingLanguages(t *testing.T) {
	tx, err := db.BeginTxx(context.Background(), nil)
	assert.Nilf(t, err, "Failed to begin transaction: %v", err)
	defer tx.Rollback()

	langs, err := ListEnabledProgrammingLanguages(tx)
	if err != nil {
		t.Fatal(err)
	}

	if len(langs) == 0 {
		t.Fatal("Expected at least one language")
	}

	for _, lang := range langs {
		assert.NotEmptyf(t, lang.ID, "Expected non-empty ID")
		assert.NotEmptyf(t, lang.Name, "Expected non-empty Name")
		assert.NotEmptyf(t, lang.CodeFilename, "Expected non-empty CodeFilename")
		assert.NotEmptyf(t, lang.ExecuteCommand, "Expected non-empty ExecuteCommand")
	}
}
