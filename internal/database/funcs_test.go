package database

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	provider, err := NewPostgresContainerTestDBProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer provider.Close()

	db := provider.GetTestDB()
	// asert db is not null
	if db == nil {
		t.Fatal("db is nil")
	}
}
