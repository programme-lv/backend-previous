package adapters

import (
	"context"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/common/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/common/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/eval/query"
)

type EvaluationPostgresRepo struct {
	db qrm.DB
}

func NewEvaluationPostgresRepo(db qrm.DB) EvaluationPostgresRepo {
	return EvaluationPostgresRepo{db: db}
}

func (e EvaluationPostgresRepo) allEvaluationRecords(ctx context.Context) ([]model.Evaluations, error) {
	stmt := postgres.SELECT(table.Evaluations.AllColumns).FROM(table.Evaluations)
	var records []model.Evaluations
	err := stmt.QueryContext(ctx, e.db, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (e EvaluationPostgresRepo) allSubmissionRecords(ctx context.Context) ([]model.TaskSubmissions, error) {
	stmt := postgres.SELECT(table.TaskSubmissions.AllColumns).FROM(table.TaskSubmissions)
	var records []model.TaskSubmissions
	err := stmt.QueryContext(ctx, e.db, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (e EvaluationPostgresRepo) allProgrammingLanguageRecords(ctx context.Context) ([]model.ProgrammingLanguages, error) {
	stmt := postgres.SELECT(table.ProgrammingLanguages.AllColumns).FROM(table.ProgrammingLanguages)
	var records []model.ProgrammingLanguages
	err := stmt.QueryContext(ctx, e.db, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (e EvaluationPostgresRepo) allUserRecords(ctx context.Context) ([]model.Users, error) {
	stmt := postgres.SELECT(table.Users.AllColumns).FROM(table.Users)
	var records []model.Users
	err := stmt.QueryContext(ctx, e.db, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (e EvaluationPostgresRepo) allTaskRecords(ctx context.Context) ([]model.Tasks, error) {
	stmt := postgres.SELECT(table.Tasks.AllColumns).FROM(table.Tasks)
	var records []model.Tasks
	err := stmt.QueryContext(ctx, e.db, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (e EvaluationPostgresRepo) allTaskVersionRecords(ctx context.Context) ([]model.TaskVersions, error) {
	stmt := postgres.SELECT(table.TaskVersions.AllColumns).FROM(table.TaskVersions)
	var records []model.TaskVersions
	err := stmt.QueryContext(ctx, e.db, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (e EvaluationPostgresRepo) AllSubmissions(ctx context.Context) ([]query.Submission, error) {
	var err error

	var submissionRecords []model.TaskSubmissions
	submissionRecords, err = e.allSubmissionRecords(ctx)
	if err != nil {
		return nil, err
	}
	var evaluationRecords []model.Evaluations
	evaluationRecords, err = e.allEvaluationRecords(ctx)
	if err != nil {
		return nil, err
	}
	mapEvalIDToEvaluation := make(map[int64]*model.Evaluations)
	for _, evaluationRecord := range evaluationRecords {
		mapEvalIDToEvaluation[evaluationRecord.ID] = &evaluationRecord
	}

	var programmingLanguageRecords []model.ProgrammingLanguages
	programmingLanguageRecords, err = e.allProgrammingLanguageRecords(ctx)
	if err != nil {
		return nil, err
	}

	var userRecords []model.Users
	userRecords, err = e.allUserRecords(ctx)
	if err != nil {
		return nil, err
	}
	mapUserIDToUsername := make(map[int64]string)
	for _, userRecord := range userRecords {
		mapUserIDToUsername[userRecord.ID] = userRecord.Username
	}

	mapProgLangIDToFullName := make(map[string]string)
	for _, programmingLanguageRecord := range programmingLanguageRecords {
		mapProgLangIDToFullName[programmingLanguageRecord.ID] = programmingLanguageRecord.FullName
	}

	var taskRecords []model.Tasks
	taskRecords, err = e.allTaskRecords(ctx)
	if err != nil {
		return nil, err
	}
	mapTaskIDToTask := make(map[int64]*model.Tasks)
	for _, taskRecord := range taskRecords {
		mapTaskIDToTask[taskRecord.ID] = &taskRecord
	}

	var taskVersionRecords []model.TaskVersions
	taskVersionRecords, err = e.allTaskVersionRecords(ctx)
	if err != nil {
		return nil, err
	}
	mapTaskVersionIDToTaskVersion := make(map[int64]*model.TaskVersions)
	for _, taskVersionRecord := range taskVersionRecords {
		mapTaskVersionIDToTaskVersion[taskVersionRecord.ID] = &taskVersionRecord
	}

	var submissions []query.Submission
	for _, submission := range submissionRecords {
		task, taskFound := mapTaskIDToTask[submission.TaskID]
		if !taskFound {
			continue
		}
		if task.StableVersionID == nil {
			continue
		}
		taskVersion, taskVersionFound := mapTaskVersionIDToTaskVersion[*task.StableVersionID]
		if !taskVersionFound {
			continue
		}

		progLangFullName, progLangFound := mapProgLangIDToFullName[submission.ProgrammingLangID]
		if !progLangFound {
			continue
		}

		var evaluationRes *query.Evaluation = nil
		if submission.VisibleEvalID != nil {
			visibleEval, evalFound := mapEvalIDToEvaluation[*submission.VisibleEvalID]
			if evalFound {
				evaluationRes = &query.Evaluation{
					ID:         visibleEval.ID,
					Status:     visibleEval.EvalStatusID,
					TotalScore: visibleEval.EvalTotalScore,
					MaxScore:   visibleEval.EvalPossibleScore,
				}
			} else {
				continue
			}
		}

		username, ok := mapUserIDToUsername[submission.UserID]
		if !ok {
			continue
		}

		submissions = append(submissions, query.Submission{
			ID:               submission.ID,
			TaskFullName:     taskVersion.FullName,
			TaskCode:         taskVersion.ShortCode,
			AuthorUsername:   username,
			ProgLangFullName: progLangFullName,
			SubmissionCode:   submission.Submission,
			EvaluationRes:    evaluationRes,
			CreatedAt:        submission.CreatedAt,
		})
	}

	return submissions, nil
}

var _ query.AllSubmissionsReadModel = (*EvaluationPostgresRepo)(nil)
