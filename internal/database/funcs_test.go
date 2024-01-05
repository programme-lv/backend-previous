package database

import (
	"log"
	"os"
	"sync"
	"testing"

	"github.com/programme-lv/backend/internal/database/testdb"
	"github.com/stretchr/testify/assert"
)

var (
	once       sync.Once
	dbProvider testdb.DBTestcontainer
)

func TestMain(m *testing.M) {
	var err error
	dbProvider, err = testdb.NewPostgresTestcontainer()
	if err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	dbProvider.Close()
	os.Exit(exitCode)
}

func TestCreateUser(t *testing.T) {
	db := dbProvider.GetConn()
	assert.NotNil(t, db)

	username := "username"
	password := "password"
	email := "email@gmail.com"
	firstName := "firstName"
	lastName := "lastName"

	assert.Nil(t, CreateUser(db, username, hashPassword(password), email, firstName, lastName))

	user, err := SelectUserByUsername(db, username)
	assert.Nil(t, err)

	assert.Equal(t, user.Username, username)
	assert.Equal(t, user.Email, email)
	assert.Equal(t, user.FirstName, firstName)
	assert.Equal(t, user.LastName, lastName)
	assert.Equal(t, user.IsAdmin, false)

	assert.True(t, passwordsMatch(user.HashedPassword, password))

	assert.Nil(t, DeleteUserById(db, user.ID))
}
