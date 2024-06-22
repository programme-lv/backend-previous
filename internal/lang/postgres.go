package lang

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/common/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/common/database/proglv/public/table"
)

type proglangRepoPostgresImpl struct {
	db qrm.DB
}

func (p proglangRepoPostgresImpl) GetAllProgrammingLanguages() ([]*ProgrammingLanguage, error) {
	stmt := postgres.SELECT(table.ProgrammingLanguages.AllColumns).
		FROM(table.ProgrammingLanguages)

	var records []model.ProgrammingLanguages
	err := stmt.Query(p.db, &records)
	if err != nil {
		return nil, err
	}

	var languages []*ProgrammingLanguage
	for _, record := range records {
		languages = append(languages, p.mapProgLangTableRowToDomainObject(record))
	}

	return languages, nil
}

func (p proglangRepoPostgresImpl) DoesLanguageExistByID(id string) (bool, error) {
	stmt := postgres.SELECT(table.ProgrammingLanguages.AllColumns).
		FROM(table.ProgrammingLanguages).
		WHERE(table.ProgrammingLanguages.ID.EQ(postgres.String(id))).
		LIMIT(1)

	var record model.ProgrammingLanguages
	err := stmt.Query(p.db, &record)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return false, nil
		}
		return false, err
	}

	return record.ID != "", nil
}

func (p proglangRepoPostgresImpl) GetProgrammingLanguageByID(id string) (*ProgrammingLanguage, error) {
	stmt := postgres.SELECT(table.ProgrammingLanguages.AllColumns).
		FROM(table.ProgrammingLanguages).
		WHERE(table.ProgrammingLanguages.ID.EQ(postgres.String(id)))

	var record model.ProgrammingLanguages
	err := stmt.Query(p.db, &record)
	if err != nil {
		return nil, err
	}

	return p.mapProgLangTableRowToDomainObject(record), nil
}

func (p proglangRepoPostgresImpl) mapProgLangTableRowToDomainObject(record model.ProgrammingLanguages) *ProgrammingLanguage {
	return &ProgrammingLanguage{
		ID:                record.ID,
		Name:              record.FullName,
		CodeFilename:      record.CodeFilename,
		CompileCommand:    record.CompileCmd,
		ExecuteCommand:    record.ExecuteCmd,
		EnvVersionCommand: record.EnvVersionCmd,
		HelloWorldCode:    record.HelloWorldCode,
		MonacoID:          record.MonacoID,
		Enabled:           record.Enabled,
	}
}

var _ progLangRepo = proglangRepoPostgresImpl{}
