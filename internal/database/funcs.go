package database

import "github.com/jmoiron/sqlx"

func SelectUserByUsername(db *sqlx.DB, username string) (*User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	return &user, err
}
