package internal

import "github.com/programme-lv/backend/internal/domain"

type UserService interface {
	Login(username, password string) (*domain.User, error)
	Register(username, password, email, firstName, lastName string) (*domain.User, error)
	GetUserByID(id int64) (*domain.User, error)
}

type AuthSessionManager interface {
	PutUserIDIntoCtx(ctx interface{}, userID int64)
	GetUserIDFromCtx(ctx interface{}) (int64, error)
	PopUserIDFromCtx(ctx interface{}) (int64, error)
}

type UserRepo interface {
	DoesUserExistByUsername(username string) (bool, error)
	DoesUserExistByEmail(email string) (bool, error)
	CreateUser(username string, hashedPassword []byte, email, firstName, lastName string) (int64, error)
	GetUserByID(id int64) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
}