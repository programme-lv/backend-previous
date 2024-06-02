package service

import (
	"github.com/programme-lv/backend/internal/domain"
)

type UserService struct {
}

func (s *UserService) Login(username, password string) (*domain.User, error) {
	if username == "" || password == "" {

	}
	panic("not implemented")
}

func (s *UserService) Register(username, password, email, firstName, lastName string) (*domain.User, error) {

	//if username == "" || password == "" {
	//	return nil, ErrUsernameOrPasswordEmpty(getGQLReqLang(ctx))
	//}
	//if len(password) < 7 {
	//	return nil, ErrPasswordTooShort(getGQLReqLang(ctx), 7)
	//}
	//if len(password) > 31 {
	//	return nil, ErrPasswordTooLong(getGQLReqLang(ctx), 31)
	//}
	//if len(username) < 2 {
	//	return nil, ErrUsernameTooShort(getGQLReqLang(ctx), 2)
	//}
	//if len(username) > 14 {
	//	return nil, ErrUsernameTooLong(getGQLReqLang(ctx), 14)
	//}
	//
	//usernameExists, err := database.DoesUserExistByUsername(r.PostgresDB, username)
	//if err != nil {
	//	return nil, ErrInternalServer(getGQLReqLang(ctx))
	//}
	//if usernameExists {
	//	return nil, ErrUserWithThatUsernameExists(getGQLReqLang(ctx))
	//}
	//
	//emailExists, err := database.DoesUserExistByEmail(r.PostgresDB, email)
	//if err != nil {
	//	return nil, err
	//}
	//if emailExists {
	//	return nil, ErrUserWithThatEmailExists(getGQLReqLang(ctx))
	//}
	//
	//// validate email
	//_, err = mail.ParseAddress(email)
	//if err != nil {
	//	return nil, ErrInvalidEmail(getGQLReqLang(ctx))
	//}
	//
	//hashedPassword, err := r.HashPassword(password)
	//if err != nil {
	//	return nil, ErrInternalServer(getGQLReqLang(ctx))
	//}
	//
	//err = database.CreateUser(r.PostgresDB, username, hashedPassword, email, firstName, lastName)
	//if err != nil {
	//	return nil, ErrInternalServer(getGQLReqLang(ctx))
	//}
	//
	//user, err := database.SelectUserByUsername(r.PostgresDB, username)
	//if err != nil {
	//	return nil, ErrInternalServer(getGQLReqLang(ctx))
	//}
	panic("not implemented")
}
