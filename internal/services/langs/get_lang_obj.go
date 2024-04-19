package langs

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
)

func GetLangObj(db qrm.DB, langID string) (*objects.ProgrammingLanguage, error) {
	langRecord, err := selectLangRecord(db, langID)
	if err != nil {
		return nil, err
	}
	res := objects.ProgrammingLanguage{
		ID:                langID,
		Name:              "",
		CodeFilename:      "",
		CompileCommand:    new(string),
		ExecuteCommand:    "",
		EnvVersionCommand: new(string),
		HelloWorldCode:    new(string),
		MonacoID:          new(string),
		Enabled:           true,
	}
	if langRecord != nil {
		res.Name = langRecord.FullName
		res.CodeFilename = langRecord.CodeFilename
		res.CompileCommand = langRecord.CompileCmd
		res.ExecuteCommand = langRecord.ExecuteCmd
		res.EnvVersionCommand = langRecord.EnvVersionCmd
		res.HelloWorldCode = langRecord.HelloWorldCode
		res.MonacoID = langRecord.MonacoID
		res.Enabled = langRecord.Enabled
	}
	return &res, nil
}

func selectLangRecord(db qrm.DB, langID string) (*model.ProgrammingLanguages, error) {
	stmt := postgres.SELECT(table.ProgrammingLanguages.AllColumns).
		FROM(table.ProgrammingLanguages).
		WHERE(table.ProgrammingLanguages.ID.EQ(postgres.String(langID)))
	var record model.ProgrammingLanguages
	err := stmt.Query(db, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
