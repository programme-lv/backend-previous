package models

import "time"

type User struct {
	ID             int64      `db:"id"`
	Username       string     `db:"username"`
	HashedPassword string     `db:"hashed_password"`
	Email          string     `db:"email"`
	FirstName      string     `db:"first_name"`
	LastName       string     `db:"last_name"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}
