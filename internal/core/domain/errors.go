package domain

import (
	"errors"
)

var (
	ErrPasswordTooLong = errors.New("password is too long")

	ErrUsernameTooShort = errors.New("username is too short")
	ErrUsernameTooLong  = errors.New("username is too long")
)

type I18NError struct {
}
