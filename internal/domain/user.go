package domain

type User struct {
	ID        int64
	Username  string
	Email     string
	FirstName string
	LastName  string
	EncPasswd []byte
	IsAdmin   bool
}
