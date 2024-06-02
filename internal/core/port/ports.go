package port

import (
	"github.com/programme-lv/backend/internal/domain"
)

type UserRepository interface {
	CreateUser(username, password, email, firstName, lastName string) error
	GetUserByUsername(username string) (*domain.User, error)
	DoesUserExistByUsername(username string) (bool, error)
	DoesUserExistByEmail(email string) (bool, error)
}
