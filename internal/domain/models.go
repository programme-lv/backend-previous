package domain

type User struct {
	ID             int64
	Username       string
	Email          string
	FirstName      string
	LastName       string
	IsAdmin        bool
	HashedPassword []byte
}
