package adapters

import (
	"context"
	"fmt"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/programme-lv/backend/internal/common/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/common/database/proglv/public/table"
	eval2 "github.com/programme-lv/backend/internal/eval"
	"github.com/programme-lv/backend/internal/eval/query"
	"time"
)

type EvaluationPostgresRepo struct {
	db qrm.DB
}

func (e EvaluationPostgresRepo) GetSubmissionByID(ctx context.Context, uuid uuid.UUID) (*query.Submission, error) {
	submissions, err := e.AllSubmissions(ctx)
	if err != nil {
		return nil, err
	}
	if submissions == nil {
		return nil, fmt.Errorf("no submissions found")
	}

	for _, submission := range submissions {
		if submission.UUID == uuid {
			return submission, nil
		}
	}
	return nil, fmt.Errorf("submission not found")
}

func (e EvaluationPostgresRepo) AddSubmission(ctx context.Context, submission eval2.Submission) error {
	stmt := table.TaskSubmissions.INSERT(table.TaskSubmissions.AllColumns).
		MODEL(&model.TaskSubmissions{
			UserID:            submission.AuthorID(),
			TaskID:            submission.TaskID(),
			ProgrammingLangID: submission.ProgrammingLanguageID(),
			Submission:        submission.MessageBody(),
			CreatedAt:         time.Now(),
			Hidden:            false,
			VisibleEvalID:     nil,
			ID:                submission.UUID(),
		})
	_, err := stmt.ExecContext(ctx, e.db)
	if err != nil {
		return err
	}
	return nil
}

func NewEvaluationPostgresRepo(db qrm.DB) EvaluationPostgresRepo {
	return EvaluationPostgresRepo{db: db}
}

func (e EvaluationPostgresRepo) NextSubmissionID(ctx context.Context) (int64, error) {
	var id struct {
		ID int64
	}
	err := postgres.SELECT(postgres.Raw("nextval('task_submissions_id_seq'::regclass)").AS("id")).
		Query(e.db, &id)
	if err != nil {
		return 0, err
	}

	return id.ID, nil
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

func (e EvaluationPostgresRepo) allRuntimeDataRecords(ctx context.Context) ([]model.RuntimeData, error) {
	stmt := postgres.SELECT(table.RuntimeData.AllColumns).FROM(table.RuntimeData)
	var records []model.RuntimeData
	err := stmt.QueryContext(ctx, e.db, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (e EvaluationPostgresRepo) AllSubmissions(ctx context.Context) ([]*query.Submission, error) {
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

	mapProgLangIDToProgLang := make(map[string]*model.ProgrammingLanguages)
	for _, programmingLanguageRecord := range programmingLanguageRecords {
		mapProgLangIDToProgLang[programmingLanguageRecord.ID] = &programmingLanguageRecord
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

	var runtimeDataRecords []model.RuntimeData
	runtimeDataRecords, err = e.allRuntimeDataRecords(ctx)
	if err != nil {
		return nil, err
	}

	mapEvalIDToRuntimeData := make(map[int64]*model.RuntimeData)
	for _, runtimeDataRecord := range runtimeDataRecords {
		mapEvalIDToRuntimeData[runtimeDataRecord.ID] = &runtimeDataRecord
	}

	var submissions []*query.Submission
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

		progLang, progLangFound := mapProgLangIDToProgLang[submission.ProgrammingLangID]
		if !progLangFound {
			continue
		}

		var evaluationRes *query.Evaluation = nil
		if submission.VisibleEvalID != nil {
			visibleEval, evalFound := mapEvalIDToEvaluation[*submission.VisibleEvalID]
			if evalFound {
				evaluationRes = &query.Evaluation{
					ID:           visibleEval.ID,
					Status:       visibleEval.EvalStatusID,
					TotalScore:   visibleEval.EvalTotalScore,
					MaxScore:     visibleEval.EvalPossibleScore,
					CompileRData: nil,
				}
				if visibleEval.CompilationDataID != nil {
					compileData, compileDataFound := mapEvalIDToRuntimeData[*visibleEval.CompilationDataID]
					if compileDataFound {
						evaluationRes.CompileRData = &query.RuntimeData{
							TimeMillis: int(*compileData.TimeMillis),
							MemoryKB:   int(*compileData.MemoryKibibytes),
							ExitCode:   int(*compileData.ExitCode),
							Stdout:     *compileData.Stdout,
							Stderr:     *compileData.Stderr,
						}
					}
				}
			} else {
				continue
			}
		}

		username, ok := mapUserIDToUsername[submission.UserID]
		if !ok {
			continue
		}

		submissions = append(submissions, &query.Submission{
			UUID:             submission.ID,
			TaskFullName:     taskVersion.FullName,
			TaskCode:         taskVersion.ShortCode,
			AuthorUsername:   username,
			ProgLangID:       progLang.ID,
			ProgLangFullName: progLang.FullName,
			SubmissionCode:   submission.Submission,
			EvaluationRes:    evaluationRes,
			CreatedAt:        submission.CreatedAt,
		})
	}

	return submissions, nil
}

var _ query.AllSubmissionsReadModel = (*EvaluationPostgresRepo)(nil)
var _ eval2.Repository = (*EvaluationPostgresRepo)(nil)
