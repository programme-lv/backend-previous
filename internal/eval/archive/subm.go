package archive

//
//import (
//	"fmt"
//	"github.com/go-jet/jet/v2/postgres"
//	"github.com/jmoiron/sqlx"
//	"github.com/programme-lv/backend/internal/common/database/proglv/public/model"
//	"github.com/programme-lv/backend/internal/common/database/proglv/public/table"
//	"github.com/programme-lv/backend/internal/task"
//	"log/slog"
//)
//
//type Service interface {
//	// ListUserSolvedPublishedTasks returns a list of all tasks that the user has solved
//	// and that are published. As of now the user is considered to have solved the task
//	// if they have submitted a solution that has maximum score.
//	ListUserSolvedPublishedTasks(userID int64) ([]*task.Task, error)
//
//	//ListSubmissionsWithMaxScore(userID int64) ([]*domain.TaskSubmission, error)
//}
//
//type submissionRepo interface {
//	ListSolvedTaskIDs(userID int64) ([]int64, error)
//	ListPublicSubmissions() ([]*TaskSubmission, error)
//}
//
//type service struct {
//	repo    submissionRepo
//	db      *sqlx.DB
//	logger  *slog.Logger
//	taskSrv task.Service
//}
//
//var _ Service = &service{}
//
//func NewService(db *sqlx.DB, taskSrv task.Service) Service {
//	logger := slog.Default().With("service", "submission")
//	return &service{db: db, logger: logger, taskSrv: taskSrv}
//}
//
//func (s *service) ListUserSolvedPublishedTasks(userID int64) ([]*task.Task, error) {
//	publishedTasks, err := s.taskSrv.ListPublishedTasks()
//	if err != nil {
//		s.logger.Error(fmt.Sprintf("listing published tasks: %v", err))
//		return nil, err
//	}
//
//	solvedTaskIDs, err := s.listSolvedTaskIDs(userID)
//	if err != nil {
//		return nil, err
//	}
//
//	solvedPublishedTasks := make([]*task.Task, 0, len(publishedTasks))
//	for _, pTask := range publishedTasks {
//		for _, solvedTaskID := range solvedTaskIDs {
//			if pTask.id == solvedTaskID {
//				solvedPublishedTasks = append(solvedPublishedTasks, pTask)
//				break
//			}
//		}
//	}
//
//	return solvedPublishedTasks, nil
//}
//
//func (s *service) listSolvedTaskIDs(userID int64) ([]int64, error) {
//	// select submissions (task_submissions table) that belong to a user
//	// join them with evaluations (evaluations table) on submission.visible_eval_id = evaluation.id
//	// where evaluation.eval_total_score = eval_possible_score
//	// project onto task_id of submissions (leave unique values)
//
//	query := postgres.SELECT(postgres.DISTINCT(table.TaskSubmissions.taskID)).
//		FROM(table.TaskSubmissions.
//			INNER_JOIN(table.Evaluations, table.TaskSubmissions.VisibleEvalID.EQ(table.Evaluations.id))).
//		WHERE(table.TaskSubmissions.UserID.EQ(postgres.Int64(userID)).
//			AND(table.Evaluations.EvalTotalScore.EQ(table.Evaluations.EvalTotalScore)))
//
//	var taskSubmissionRecords []model.TaskSubmissions
//	err := query.Query(s.db, &taskSubmissionRecords)
//	if err != nil {
//		return nil, fmt.Errorf("failed to execute query: %w", err)
//	}
//
//	var taskIDs []int64 = make([]int64, len(taskSubmissionRecords))
//	for i, record := range taskSubmissionRecords {
//		taskIDs[i] = record.taskID
//	}
//
//	return taskIDs, nil
//}
