package testdb

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestNewMigratedPostgresTestcontainer(t *testing.T) {
	tc, err := NewMigratedPostgresTestcontainer()
	if err != nil {
		t.Errorf("NewPostgresTestcontainer() error = %v", err)
		return
	}
	defer tc.Close()
	db := tc.GetConn()
	// check with sqlx whether table "flyway_schema_history" exists
	var tableExists bool
	err = db.Get(&tableExists, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'flyway_schema_history')")
	if err != nil {
		t.Errorf("db.Get() error = %v", err)
		return
	}
	if !tableExists {
		t.Errorf("table 'flyway_schema_history' does not exist")
		return
	}
}
