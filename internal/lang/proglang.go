package lang

import (
	"github.com/go-jet/jet/v2/qrm"
	"log/slog"
)

type Service interface {
	GetProgrammingLanguageByID(id string) (*ProgrammingLanguage, error)
	ListProgrammingLanguages() ([]*ProgrammingLanguage, error)
}

type progLangRepo interface {
	DoesLanguageExistByID(id string) (bool, error)
	GetProgrammingLanguageByID(id string) (*ProgrammingLanguage, error)
	GetAllProgrammingLanguages() ([]*ProgrammingLanguage, error)
}

type service struct {
	repo   progLangRepo
	logger *slog.Logger
}

var _ Service = service{}

func (s service) ListProgrammingLanguages() ([]*ProgrammingLanguage, error) {
	languages, err := s.repo.GetAllProgrammingLanguages()
	if err != nil {
		s.logger.Error("getting all programming languages", err)
		return nil, err
	}
	return languages, err
}

func (s service) GetProgrammingLanguageByID(id string) (*ProgrammingLanguage, error) {
	exists, err := s.repo.DoesLanguageExistByID(id)
	if err != nil {
		s.logger.Error("checking if language exists by id", err)
		return nil, err
	}

	if !exists {
		return nil, newErrorLanguageNotFound()
	}

	language, err := s.repo.GetProgrammingLanguageByID(id)
	if err != nil {
		s.logger.Error("getting language by id", err)
		return nil, err
	}

	return language, nil
}

var _ Service = service{}

func NewService(db qrm.DB) Service {
	repo := proglangRepoPostgresImpl{db: db}
	logger := slog.Default().With("service", "lang")
	return &service{repo: repo, logger: logger}
}
