package domain

type UserService interface {
	Login(username, password string) *User
}
