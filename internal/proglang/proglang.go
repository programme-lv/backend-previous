package proglang

import (
	"github.com/programme-lv/backend/internal/domain"
	"log/slog"
)

type Service interface {
	GetProgrammingLanguageByID(id string) (*domain.ProgrammingLanguage, error)
}

type progLangRepo interface {
	DoesLanguageExistByID(id string) (bool, error)
	GetProgrammingLanguageByID(id string) (*domain.ProgrammingLanguage, error)
}

type service struct {
	repo   progLangRepo
	logger *slog.Logger
}

func (s service) GetProgrammingLanguageByID(id string) (*domain.ProgrammingLanguage, error) {
	exists, err := s.repo.DoesLanguageExistByID(id)
	if err != nil {
		s.logger.Error("checking if language exists by id", err)
		return nil, domain.NewErrorInternalServer()
	}

	if !exists {
		return nil, newErrorLanguageNotFound()
	}

	language, err := s.repo.GetProgrammingLanguageByID(id)
	if err != nil {
		s.logger.Error("getting language by id", err)
		return nil, domain.NewErrorInternalServer()
	}

	return language, nil
}

var _ Service = service{}

func NewService(repo progLangRepo) Service {
	logger := slog.Default().With("service", "proglang")
	return &service{repo: repo, logger: logger}
}
