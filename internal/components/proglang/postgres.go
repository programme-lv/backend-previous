package proglang

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/domain"
)

type proglangRepoPostgresImpl struct {
	db qrm.DB
}

func (p proglangRepoPostgresImpl) DoesLanguageExistByID(id string) (bool, error) {
	stmt := postgres.SELECT(table.ProgrammingLanguages.AllColumns).
		FROM(table.ProgrammingLanguages).
		WHERE(table.ProgrammingLanguages.ID.EQ(postgres.String(id))).
		LIMIT(1)

	var record model.ProgrammingLanguages
	err := stmt.Query(p.db, &record)
	if err != nil {
		return false, err
	}

	return record.ID != "", nil
}

func (p proglangRepoPostgresImpl) GetProgrammingLanguageByID(id string) (*domain.ProgrammingLanguage, error) {
	stmt := postgres.SELECT(table.ProgrammingLanguages.AllColumns).
		FROM(table.ProgrammingLanguages).
		WHERE(table.ProgrammingLanguages.ID.EQ(postgres.String(id)))

	var record model.ProgrammingLanguages
	err := stmt.Query(p.db, &record)
	if err != nil {
		return nil, err
	}

	return &domain.ProgrammingLanguage{
		ID:                record.ID,
		Name:              record.FullName,
		CodeFilename:      record.CodeFilename,
		CompileCommand:    record.CompileCmd,
		ExecuteCommand:    record.ExecuteCmd,
		EnvVersionCommand: record.EnvVersionCmd,
		HelloWorldCode:    record.HelloWorldCode,
		MonacoID:          record.MonacoID,
		Enabled:           record.Enabled,
	}, nil
}

var _ progLangRepo = proglangRepoPostgresImpl{}
