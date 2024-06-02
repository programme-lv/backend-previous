package service

import (
	"github.com/programme-lv/backend/internal/core/domain"
	"github.com/programme-lv/backend/internal/core/port"
	domain2 "github.com/programme-lv/backend/internal/domain"
)

type UserService struct {
	repo port.UserRepository
}

func (s *UserService) Register(username, password, email, firstName, lastName string) error {
	if username == "" || password == "" {
		return domain.ErrUsernameOrPasswordEmpty{}
	}
	if len(password) < 8 {
		return domain.ErrPasswordTooShort{Min: 8, Current: len(password)}
	}
	if len(password) > 32 {
		return domain2.ErrPasswordTooLong
	}
	if len(username) < 3 {
		return domain2.ErrUsernameTooShort
	}
	if len(username) > 15 {
		return domain2.ErrUsernameTooLong
	}
	return s.repo.CreateUser(username, password, email, firstName, lastName)
}
