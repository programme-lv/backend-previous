package langs

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/domain"
)

func ListEnabledProgrammingLanguages(db qrm.Queryable) ([]domain.ProgrammingLanguage, error) {
	stmt := postgres.SELECT(table.ProgrammingLanguages.AllColumns).
		FROM(table.ProgrammingLanguages).
		WHERE(table.ProgrammingLanguages.Enabled.EQ(postgres.Bool(true)))

	var langs []model.ProgrammingLanguages
	err := stmt.Query(db, &langs)
	if err != nil {
		return nil, err
	}

	res := make([]domain.ProgrammingLanguage, 0)
	for _, lang := range langs {
		res = append(res, domain.ProgrammingLanguage{
			ID:                lang.ID,
			Name:              lang.FullName,
			CodeFilename:      lang.CodeFilename,
			CompileCommand:    lang.CompileCmd,
			ExecuteCommand:    lang.ExecuteCmd,
			EnvVersionCommand: lang.EnvVersionCmd,
			HelloWorldCode:    lang.HelloWorldCode,
			MonacoID:          lang.MonacoID,
		})
	}

	return res, nil
}
