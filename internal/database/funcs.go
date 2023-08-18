package database

import "github.com/jmoiron/sqlx"

func SelectUserByUsername(db *sqlx.DB, username string) (*User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	return &user, err
}

func DoesUserExistByUsername(db *sqlx.DB, username string) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE username = $1", username)
	return count > 0, err
}

func DoesUserExistByEmail(db *sqlx.DB, email string) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", email)
	return count > 0, err
}

func CreateUser(db *sqlx.DB, username string, hashed_password string, email string,
	firstName string, lastName string) error {
	_, err := db.Exec("INSERT INTO users (username, hashed_password, email, first_name, last_name, created_at) VALUES ($1, $2, $3, $4, $5, now())", username, hashed_password, email, firstName, lastName)
	return err
}
