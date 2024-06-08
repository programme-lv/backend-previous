package user

import (
	"github.com/programme-lv/backend/internal/domain"
)

type Service interface {
	Login(username, password string) (*domain.User, error)
	Register(username, password, email, firstName, lastName string) (*domain.User, error)
	GetUserByID(id int64) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

func (s service) Login(username, password string) (*domain.User, error) {
	panic("implement me")
}

func (s service) Register(username, password, email, firstName, lastName string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetUserByID(id int64) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetUserByUsername(username string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

var _ Service = &service{}
