package submissions

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func ListVisibleTaskSubmissionRowsWithEvaluation(db qrm.DB) ([]model.TaskSubmissions, error) {
	stmt := postgres.SELECT(table.TaskSubmissions.AllColumns).
		FROM(table.TaskSubmissions).
		WHERE(table.TaskSubmissions.VisibleEvalID.IS_NOT_NULL().AND(
			table.TaskSubmissions.Hidden.EQ(postgres.Bool(false))))

	var records []model.TaskSubmissions
	err := stmt.Query(db, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}
