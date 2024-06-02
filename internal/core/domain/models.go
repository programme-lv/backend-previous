package model

type User struct {
	ID        int64
	Username  string
	Email     string
	FirstName string
	LastName  string
	IsAdmin   bool
}
