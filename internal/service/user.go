package service

import (
	"fmt"
	"github.com/programme-lv/backend/internal"
	"github.com/programme-lv/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/mail"
)

type UserService struct {
	userRepo internal.UserRepo
	logger   *slog.Logger
}

var _ internal.UserService = &UserService{}

func NewUserService(userRepo internal.UserRepo, logger *slog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *UserService) Login(username, password string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		doesExist, err := s.userRepo.DoesUserExistByUsername(username)
		if err != nil {
			return nil, fmt.Errorf("checking if user exists by username: %w", err)
		}
		if !doesExist {
			return nil, domain.NewErrorUsernameOrPasswordIncorrect()
		} else {
			return nil, fmt.Errorf("getting user by username: %w", err)
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, domain.NewErrorUsernameOrPasswordIncorrect()
	}

	return user, nil
}

func (s *UserService) Register(username, password, email, firstName, lastName string) (*domain.User, error) {
	const minPasswordLength = 7
	const maxPasswordLength = 31

	const minUsernameLength = 2
	const maxUsernameLength = 14

	if username == "" || password == "" {
		return nil, domain.NewErrorUsernameOrPasswordEmpty()
	}

	if len(password) < minPasswordLength {
		return nil, domain.NewErrorPasswordTooShort(minPasswordLength)
	}

	if len(password) > maxPasswordLength {
		return nil, domain.NewErrorPasswordTooLong(maxPasswordLength)
	}

	if len(username) < minUsernameLength {
		return nil, domain.NewErrorUsernameTooShort(minUsernameLength)
	}

	if len(username) > maxUsernameLength {
		return nil, domain.NewErrorUsernameTooLong(maxUsernameLength)
	}

	usernameExists, err := s.userRepo.DoesUserExistByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("checking if user exists by username: %w", err)
	}
	if usernameExists {
		return nil, domain.NewErrorUserWithUsernameExists()
	}

	emailExists, err := s.userRepo.DoesUserExistByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("checking if user exists by email: %w", err)
	}
	if emailExists {
		return nil, domain.NewErrorUserWithEmailExists()
	}

	_, err = mail.ParseAddress(email)
	if err != nil {
		return nil, domain.NewErrorInvalidEmail()
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	userId, err := s.userRepo.CreateUser(username, hashedPass, email, firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	user, err := s.userRepo.GetUserByID(userId)
	if err != nil {
		return nil, fmt.Errorf("getting user by id: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUserByID(id int64) (*domain.User, error) {
	return s.userRepo.GetUserByID(id)
}
