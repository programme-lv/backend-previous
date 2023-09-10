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

func DeleteUserById(db *sqlx.DB, id int64) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func CreateTaskSubmission(db sqlx.Ext, userId int64, taskId int64, programmingLangId string, submission string) error {
	_, err := db.Exec("INSERT INTO task_submissions (user_id, task_id, programming_lang_id, submission, created_at) VALUES ($1, $2, $3, $4, now())", userId, taskId, programmingLangId, submission)
	return err
}
