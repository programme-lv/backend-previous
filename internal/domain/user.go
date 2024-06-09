package domain

import "time"

type User struct {
	ID        int64
	Username  string
	Email     string
	FirstName string
	LastName  string
	EncPasswd []byte
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt *time.Time
}
