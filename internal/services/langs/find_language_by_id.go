package langs

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func FindLanguageByID(db qrm.DB, langID string) (*model.ProgrammingLanguages, error) {
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
