package user

import (
	"fmt"
	"github.com/programme-lv/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/mail"
)

type Service interface {
	Login(username, password string) (*domain.User, error)
	Register(username, password, email, firstName, lastName string) (*domain.User, error)
	GetUserByID(id int64) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
}

type userRepo interface {
	DoesUserExistByUsername(username string) (bool, error)
	DoesUserExistByEmail(email string) (bool, error)
	DoesUserExistByID(id int64) (bool, error)
	CreateUser(username string, hashedPassword []byte, email, firstName, lastName string) (int64, error)
	GetUserByID(id int64) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
}

type service struct {
	repo   userRepo
	logger *slog.Logger
}

func NewService(repo userRepo) Service {
	//create a logger that prefixes with user service
	logger := slog.Default().With("service", "user")
	return &service{repo: repo, logger: logger}
}

func (s service) Login(username, password string) (*domain.User, error) {
	usernameExists, err := s.repo.DoesUserExistByUsername(username)
	if err != nil {
		s.logger.Error(fmt.Sprintf("checking if user exists by username: %v", err))
		return nil, domain.NewErrorInternalServer()
	}

	if !usernameExists {
		return nil, newErrorUsernameOrPasswordIncorrect()
	}

	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		s.logger.Error(fmt.Sprintf("getting user by username: %v", err))
		return nil, domain.NewErrorInternalServer()
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncPasswd), []byte(password))
	if err != nil {
		return nil, newErrorUsernameOrPasswordIncorrect()
	}

	return user, nil
}

func (s service) Register(username, password, email, firstName, lastName string) (*domain.User, error) {
	const minPasswordLength = 7
	const maxPasswordLength = 31

	const minUsernameLength = 2
	const maxUsernameLength = 14

	if username == "" || password == "" {
		return nil, newErrorUsernameOrPasswordEmpty()
	}

	if len(password) < minPasswordLength {
		return nil, newErrorPasswordTooShort(minPasswordLength)
	}

	if len(password) > maxPasswordLength {
		return nil, newErrorPasswordTooLong(maxPasswordLength)
	}

	if len(username) < minUsernameLength {
		return nil, newErrorUsernameTooShort(minUsernameLength)
	}

	if len(username) > maxUsernameLength {
		return nil, newErrorUsernameTooLong(maxUsernameLength)
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, newErrorInvalidEmail()
	}

	usernameExists, err := s.repo.DoesUserExistByUsername(username)
	if err != nil {
		s.logger.Error(fmt.Sprintf("checking if user exists by username: %v", err))
		return nil, domain.NewErrorInternalServer()
	}

	if usernameExists {
		return nil, newErrorUsernameAlreadyExists()
	}

	emailExists, err := s.repo.DoesUserExistByEmail(email)
	if err != nil {
		s.logger.Error(fmt.Sprintf("checking if user exists by email: %v", err))
		return nil, domain.NewErrorInternalServer()
	}

	if emailExists {
		return nil, newErrorEmailAlreadyExists()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(fmt.Sprintf("hashing password: %v", err))
		return nil, domain.NewErrorInternalServer()
	}

	userID, err := s.repo.CreateUser(username, hashedPassword, email, firstName, lastName)
	if err != nil {
		s.logger.Error(fmt.Sprintf("creating user: %v", err))
		return nil, domain.NewErrorInternalServer()
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("getting user by ID: %v", err))
		return nil, domain.NewErrorInternalServer()
	}

	return user, nil
}

func (s service) GetUserByID(id int64) (*domain.User, error) {
	userExists, err := s.repo.DoesUserExistByID(id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("checking if user exists by ID: %v", err))
		return nil, domain.NewErrorInternalServer()
	}

	if !userExists {
		return nil, newErrorUserNotFound()
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("getting user by ID: %v", err))
		return nil, domain.NewErrorInternalServer()
	}

	return user, nil
}

func (s service) GetUserByUsername(username string) (*domain.User, error) {
	userExists, err := s.repo.DoesUserExistByUsername(username)
	if err != nil {
		s.logger.Error(fmt.Sprintf("checking if user exists by username: %v", err))
		return nil, domain.NewErrorInternalServer()
	}

	if !userExists {
		return nil, newErrorUserNotFound()
	}

	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		s.logger.Error(fmt.Sprintf("getting user by username: %v", err))
		return nil, domain.NewErrorInternalServer()
	}

	return user, nil
}

var _ Service = &service{}
