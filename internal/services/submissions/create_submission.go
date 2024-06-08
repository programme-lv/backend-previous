package submissions

import (
	"github.com/go-jet/jet/qrm"
	"github.com/programme-lv/backend/internal/database/postgres/proglv/public/table"
)

type CreateSubmissionParams struct {
	UserID            int64
	TaskID            int64
	ProgrammingLangID string
	Submission        string
}

func CreateSubmission(db qrm.DB, params CreateSubmissionParams) (int64, error) {
	stmt := table.TaskSubmissions.INSERT(
		table.TaskSubmissions.UserID,
		table.TaskSubmissions.TaskID,
		table.TaskSubmissions.ProgrammingLangID,
		table.TaskSubmissions.Submission,
		table.TaskSubmissions.CreatedAt,
	).VALUES(
		params.UserID,
		params.TaskID,
		params.ProgrammingLangID,
		params.Submission,
		"now()",
	).RETURNING(table.TaskSubmissions.ID)

	var record model.TaskSubmissions
	err := stmt.Query(db, &record)
	if err != nil {
		return 0, err
	}

	return record.ID, nil
}
