package user

import (
	"github.com/programme-lv/backend/internal/domain"
)

type repoBase interface {
	DoesUserExistByUsername(username string) (bool, error)
	DoesUserExistByEmail(email string) (bool, error)
	CreateUser(username string, hashedPassword []byte, email, firstName, lastName string) (int64, error)
	GetUserByID(id int64) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
}

type Repo interface {
	repoBase
	BeginTx() (RepoTx, error)
}

type RepoTx interface {
	repoBase
	Commit() error
	Rollback() error
}
