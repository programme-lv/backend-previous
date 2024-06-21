package eval

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

type submRepoImpl struct {
	db qrm.DB
}

func (s submRepoImpl) ListSolvedTaskIDs(userID int64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s submRepoImpl) ListPublicSubmissions() ([]*TaskSubmission, error) {
	stmt := postgres.SELECT(table.TaskSubmissions.AllColumns).
		FROM(table.TaskSubmissions).
		WHERE(table.TaskSubmissions.Hidden.EQ(postgres.Bool(false)))

	var records []*model.TaskSubmissions
	err := stmt.Query(s.db, &records)
	if err != nil {
		return nil, err
	}

	domainTaskSubmissions := make([]*TaskSubmission, 0, len(records))
	for _, record := range records {
		taskSubmission, errMapping := s.mapTaskSubmissionTableRowToDomainObject(record)
		if errMapping != nil {
			return nil, errMapping
		}
		domainTaskSubmissions = append(domainTaskSubmissions, taskSubmission)
	}
	return domainTaskSubmissions, nil
}

func (s submRepoImpl) mapTaskSubmissionTableRowToDomainObject(record *model.TaskSubmissions) (*TaskSubmission, error) {
	res := TaskSubmission{
		ID: record.ID,
		//Author:      record.UserID, TODO
		//Language:    record.ProgrammingLangID, TODO
		Content: record.Submission,
		//Task:        record.TaskID, TODO
		//VisibleEval: record.VisibleEvalID, TODO
		Hidden:    record.Hidden,
		CreatedAt: record.CreatedAt,
	}
	return &res, nil
}

var _ submissionRepo = submRepoImpl{}
