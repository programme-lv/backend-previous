package submissions

import (
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/langs"
	"github.com/programme-lv/backend/internal/services/objects"
	"github.com/programme-lv/backend/internal/services/tasks"
	"github.com/programme-lv/backend/internal/services/users"
)

func GetSubmissionObject(db qrm.DB, submissionID int64) (*objects.TaskSubmission, error) {
	res := objects.TaskSubmission{
		ID:          0,
		Author:      &objects.User{},
		Language:    &objects.ProgrammingLanguage{},
		Content:     "",
		Task:        &objects.Task{},
		VisibleEval: &objects.Evaluation{}, // TODO
		Hidden:      false,
		CreatedAt:   time.Time{},
	}

	submRecord, err := selectSubmissionRecord(db, submissionID)
	if err != nil {
		return nil, err
	}

	res.ID = submRecord.ID
	res.Content = submRecord.Submission
	res.CreatedAt = submRecord.CreatedAt
	res.Hidden = submRecord.Hidden

	authorObj, err := users.GetUserObj(db, submRecord.UserID)
	if err != nil {
		return nil, err
	}
	res.Author = authorObj

	langObj, err := langs.GetLangObj(db, submRecord.ProgrammingLangID)
	if err != nil {
		return nil, err
	}
	res.Language = langObj

	taskObj, err := tasks.GetTaskObjByTaskID(db, submRecord.TaskID, 2, 2)
	if err != nil {
		return nil, err
	}
	res.Task = taskObj

	evalObj, err := GetEvaluationObj(db, *submRecord.VisibleEvalID, true)
	if err != nil {
		return nil, err
	}
	res.VisibleEval = evalObj

	return &res, nil
}

func selectSubmissionRecord(db qrm.DB, submissionID int64) (*model.TaskSubmissions, error) {
	stmt := postgres.SELECT(table.TaskSubmissions.AllColumns).
		FROM(table.TaskSubmissions).
		WHERE(table.TaskSubmissions.ID.EQ(postgres.Int64(submissionID)))

	var record model.TaskSubmissions
	err := stmt.Query(db, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
