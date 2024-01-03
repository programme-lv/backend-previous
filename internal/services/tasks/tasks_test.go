package tasks

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/database/testdb"
	"github.com/programme-lv/backend/internal/services/objects"
	"github.com/stretchr/testify/assert"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	provider, err := testdb.NewPostgresTestcontainerProvider()
	if err != nil {
		panic(err)
	}
	db = provider.GetTestDB()
	defer provider.Close()

	code := m.Run()
	os.Exit(code)
}

func TestListPublishedTaskVersions(t *testing.T) {
	tx, err := db.BeginTxx(context.Background(), nil)
	assert.Nilf(t, err, "Failed to begin transaction: %v", err)
	defer tx.Rollback()

	userID, err := insertTempTestUser(tx)
	assert.Nilf(t, err, "Failed to create temp user: %v", err)

	taskVersions, err := ListPublishedTaskVersions(tx)
	assert.Nilf(t, err, "Failed to list published task versions: %v", err)
	assert.Equalf(t, 0, len(taskVersions), "Expected 0 tasks, got %d", len(taskVersions))

	target := initTargetTaskVersion()
	_, err = createTaskVersionTarget(tx, target, userID)
	assert.Nilf(t, err, "Failed to create task version target: %v", err)
	_, err = createTaskVersionTarget(tx, target, userID)
	assert.Nilf(t, err, "Failed to create task version target: %v", err)

	taskVersions, err = ListPublishedTaskVersions(tx)
	assert.Nilf(t, err, "Failed to list published task versions: %v", err)
	assert.Equalf(t, 2, len(taskVersions), "Expected 1 task, got %d", len(taskVersions))

	for _, received := range taskVersions {
		compareTaskVersion(t, target, received)
	}
}

func TestGetPublishedTaskVersionByCode(t *testing.T) {
	tx, err := db.BeginTxx(context.Background(), nil)
	assert.Nilf(t, err, "Failed to begin transaction: %v", err)
	defer tx.Rollback()

	userID, err := insertTempTestUser(tx)
	assert.Nilf(t, err, "Failed to create temp user: %v", err)

	target := initTargetTaskVersion()
	_, err = createTaskVersionTarget(tx, target, userID)
	assert.Nilf(t, err, "Failed to create task version target: %v", err)

	taskVersion, err := GetPublishedTaskVersionByCode(tx, target.Code)
	assert.Nilf(t, err, "Failed to get published task version by code: %v", err)
	compareTaskVersion(t, target, *taskVersion)
}

func TestListSolvedTasksByUserID(t *testing.T) {
	tx, err := db.BeginTxx(context.Background(), nil)
	assert.Nilf(t, err, "Failed to begin transaction: %v", err)
	defer tx.Rollback()

	userID, err := insertTempTestUser(tx)
	assert.Nilf(t, err, "Failed to create temp user: %v", err)

	solvedTasks, err := ListSolvedTasksByUserID(tx, userID)
	assert.Nilf(t, err, "Failed to list solved tasks by user id: %v", err)
	assert.Equalf(t, 0, len(solvedTasks), "Expected 0 solved tasks, got %d", len(solvedTasks))

	target := initTargetTaskVersion()
	taskVersions := make([]int64, 2)
	for i := 0; i < 2; i++ {
		taskVersionID, err := createTaskVersionTarget(tx, target, userID)
		assert.Nilf(t, err, "Failed to create task version target: %v", err)
		taskVersions[i] = taskVersionID
	}

	taskIds := make([]int64, 2)
	for i, taskVersionID := range taskVersions {
		taskVersion, err := selectTaskVersionByID(tx, taskVersionID)
		assert.Nilf(t, err, "Failed to select task version by id: %v", err)

		taskIds[i] = taskVersion.TaskID

		programmingLang, err := selectAnyProgrammingLangID(tx)
		assert.Nilf(t, err, "Failed to select any programming lang: %v", err)

		var maxScore int64 = 100
		// create an evaluation
		evaluation := model.Evaluations{
			EvalStatusID:            "F",
			EvalTotalScore:          maxScore,
			EvalPossibleScore:       &maxScore,
			TestRuntimeStatisticsID: nil,
			CompilationDataID:       nil,
			TaskVersionID:           taskVersionID,
		}
		evaluationID, err := insertEvaluation(tx, evaluation)
		assert.Nilf(t, err, "Failed to insert evaluation: %v", err)

		// create a task submission
		taskSubmission := model.TaskSubmissions{
			UserID:            userID,
			TaskID:            taskVersion.TaskID,
			ProgrammingLangID: programmingLang,
			Submission:        "123",
			VisibleEvalID:     &evaluationID,
		}
		_, err = insertTaskSubmission(tx, taskSubmission)
		assert.Nilf(t, err, "Failed to insert task submission: %v", err)
	}

	solvedTasks, err = ListSolvedTasksByUserID(tx, userID)
	assert.Nilf(t, err, "Failed to list solved tasks by user id: %v", err)
	assert.Equalf(t, 2, len(solvedTasks), "Expected 2 solved tasks, got %d", len(solvedTasks))

	assert.Equalf(t, taskIds[0], solvedTasks[0], "Expected task id %d, got %d", taskIds[0], solvedTasks[0])
	assert.Equalf(t, taskIds[1], solvedTasks[1], "Expected task id %d, got %d", taskIds[1], solvedTasks[1])
}

func compareTaskVersion(t *testing.T, expected, received objects.TaskVersion) {
	assert.Equalf(t, expected.Code, received.Code, "Expected code %s, got %s", expected.Code, received.Code)
	assert.Equalf(t, expected.Name, received.Name, "Expected name %s, got %s", expected.Name, received.Name)
	assert.Equalf(t, expected.TimeLimitMs, received.TimeLimitMs,
		"Expected time limit %d, got %d", expected.TimeLimitMs, received.TimeLimitMs)
	assert.Equalf(t, expected.MemoryLimitKb, received.MemoryLimitKb,
		"Expected memory limit %d, got %d", expected.MemoryLimitKb, received.MemoryLimitKb)

	assert.Equalf(t, expected.Description.Story, received.Description.Story,
		"Expected story %s, got %s", expected.Description.Story, received.Description.Story)
	assert.Equalf(t, expected.Description.Input, received.Description.Input,
		"Expected input %s, got %s", expected.Description.Input, received.Description.Input)
	assert.Equalf(t, expected.Description.Output, received.Description.Output,
		"Expected output %s, got %s", expected.Description.Output, received.Description.Output)

	assert.Equalf(t, len(expected.Description.Examples), len(received.Description.Examples),
		"Expected %d examples, got %d", len(expected.Description.Examples), len(received.Description.Examples))
	for i, example := range expected.Description.Examples {
		assert.Equalf(t, example.Input, received.Description.Examples[i].Input,
			"Expected example input %s, got %s", example.Input, received.Description.Examples[i].Input)
		assert.Equalf(t, example.Answer, received.Description.Examples[i].Answer,
			"Expected examples %v, got %v", expected.Description.Examples, received.Description.Examples)
	}
}

func initTargetTaskVersion() objects.TaskVersion {
	notesStr := "Piezīmes"
	return objects.TaskVersion{
		Code: "summa",
		Name: "Summa",
		Description: &objects.Description{
			Story:    "Stāsts. Saskaiti skaitļus.",
			Input:    "Ievaddati",
			Output:   "Izvaddati",
			Examples: []objects.Example{{Input: "1 2", Answer: "3"}},
			Notes:    &notesStr,
		},
		TimeLimitMs:   1024,
		MemoryLimitKb: 262144,
	}
}

// Create a task version with the given target and return its ID.
func createTaskVersionTarget(tx *sqlx.Tx, target objects.TaskVersion, userID int64) (int64, error) {

	taskID, err := insertTask(tx, userID)
	if err != nil {
		return 0, fmt.Errorf("Failed to create task: %v", err)
	}

	taskVersion := model.TaskVersions{
		TaskID:          taskID,
		ShortCode:       target.Code,
		FullName:        target.Name,
		TimeLimMs:       target.TimeLimitMs,
		MemLimKibibytes: target.MemoryLimitKb,
		TestingTypeID:   "simple",
	}
	taskVersionID, err := insertTaskVersion(tx, taskVersion)
	if err != nil {
		return 0, fmt.Errorf("Failed to create task version: %v", err)
	}

	markdownStatement := model.MarkdownStatements{
		TaskVersionID: &taskVersionID,
		LangIso6391:   "lv",
		Story:         target.Description.Story,
		Input:         target.Description.Input,
		Output:        target.Description.Output,
	}
	_, err = insertMarkdownStatement(tx, markdownStatement)
	if err != nil {
		return 0, fmt.Errorf("Failed to create markdown statement: %v", err)
	}

	for _, example := range target.Description.Examples {
		statementExample := model.StatementExamples{
			TaskVersionID: taskVersionID,
			Input:         example.Input,
			Answer:        example.Answer,
		}
		_, err = insertStatementExample(tx, statementExample)
		if err != nil {
			return 0, fmt.Errorf("Failed to create statement example: %v", err)
		}
	}

	err = updateTaskRelevantAndPublishedVersionIds(tx, taskID, taskVersionID)
	if err != nil {
		return 0, fmt.Errorf("Failed to update task relevant and published version ids: %v", err)
	}

	return taskVersionID, nil
}

// Create a temporary test user and return its ID.
func insertTempTestUser(tx *sqlx.Tx) (int64, error) {
	userCreateStmt := table.Users.INSERT(
		table.Users.Username,
		table.Users.Email,
		table.Users.HashedPassword,
		table.Users.FirstName,
		table.Users.LastName,
		table.Users.IsAdmin).
		VALUES("test", "test@gmail.com", "test", "test", "test", false).
		RETURNING(table.Users.ID)
	userDest := &model.Users{}
	err := userCreateStmt.Query(tx, userDest)
	return userDest.ID, err
}

// Create a task and return its ID.
func insertTask(tx *sqlx.Tx, userID int64) (int64, error) {
	createTaskStmt := table.Tasks.INSERT(
		table.Tasks.CreatedByID).VALUES(userID).RETURNING(table.Tasks.ID)
	taskDest := &model.Tasks{}
	err := createTaskStmt.Query(tx, taskDest)
	return taskDest.ID, err
}

// Create a task version "summa" for the given task and return its ID.
func insertTaskVersion(tx *sqlx.Tx, taskVersion model.TaskVersions) (int64, error) {
	createTaskVersionStmt := table.TaskVersions.INSERT(
		table.TaskVersions.TaskID,
		table.TaskVersions.ShortCode,
		table.TaskVersions.FullName,
		table.TaskVersions.TimeLimMs,
		table.TaskVersions.TestingTypeID,
		table.TaskVersions.MemLimKibibytes).MODEL(taskVersion).
		RETURNING(table.TaskVersions.ID)
	taskVersionDest := &model.TaskVersions{}
	err := createTaskVersionStmt.Query(tx, taskVersionDest)
	return taskVersionDest.ID, err
}

func insertMarkdownStatement(tx *sqlx.Tx, markdownStatement model.MarkdownStatements) (int64, error) {
	createMarkdownStatementStmt := table.MarkdownStatements.INSERT(
		table.MarkdownStatements.TaskVersionID,
		table.MarkdownStatements.LangIso6391,
		table.MarkdownStatements.Story,
		table.MarkdownStatements.Input,
		table.MarkdownStatements.Output).MODEL(markdownStatement).
		RETURNING(table.MarkdownStatements.ID)
	markdownStatementDest := &model.MarkdownStatements{}
	err := createMarkdownStatementStmt.Query(tx, markdownStatementDest)
	return markdownStatementDest.ID, err
}

func insertStatementExample(tx *sqlx.Tx, statementExample model.StatementExamples) (int64, error) {
	createStatementExampleStmt := table.StatementExamples.INSERT(
		table.StatementExamples.TaskVersionID,
		table.StatementExamples.Input,
		table.StatementExamples.Answer).MODEL(statementExample).
		RETURNING(table.StatementExamples.ID)
	statementExampleDest := &model.StatementExamples{}
	err := createStatementExampleStmt.Query(tx, statementExampleDest)
	return statementExampleDest.ID, err
}

// Update the task's relevant and published version IDs to the given task version ID.
func updateTaskRelevantAndPublishedVersionIds(tx *sqlx.Tx, taskID, taskVersionID int64) error {
	updateTaskStmt := table.Tasks.UPDATE(
		table.Tasks.RelevantVersionID,
		table.Tasks.PublishedVersionID).SET(
		taskVersionID, taskVersionID).WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))
	_, err := updateTaskStmt.Exec(tx)
	return err
}

func insertTaskSubmission(tx *sqlx.Tx, taskSubmission model.TaskSubmissions) (int64, error) {
	createTaskSubmissionStmt := table.TaskSubmissions.INSERT(
		table.TaskSubmissions.UserID,
		table.TaskSubmissions.TaskID,
		table.TaskSubmissions.ProgrammingLangID,
		table.TaskSubmissions.Submission,
		table.TaskSubmissions.VisibleEvalID,
		table.TaskSubmissions.Hidden,
	).MODEL(taskSubmission).
		RETURNING(table.TaskSubmissions.ID)
	taskSubmissionDest := &model.TaskSubmissions{}
	err := createTaskSubmissionStmt.Query(tx, taskSubmissionDest)
	return taskSubmissionDest.ID, err
}

func selectTaskVersionByID(tx *sqlx.Tx, taskVersionID int64) (*model.TaskVersions, error) {
	selectTaskVersionStmt := table.TaskVersions.SELECT(table.TaskVersions.AllColumns).
		WHERE(table.TaskVersions.ID.EQ(postgres.Int64(taskVersionID)))
	taskVersionDest := &model.TaskVersions{}
	err := selectTaskVersionStmt.Query(tx, taskVersionDest)
	return taskVersionDest, err
}

func selectAnyProgrammingLangID(tx *sqlx.Tx) (string, error) {
	selectProgrammingLangStmt := table.ProgrammingLanguages.SELECT(table.ProgrammingLanguages.ID).
		LIMIT(1)
	programmingLangDest := &model.ProgrammingLanguages{}
	err := selectProgrammingLangStmt.Query(tx, programmingLangDest)
	return programmingLangDest.ID, err
}

func insertEvaluation(tx *sqlx.Tx, evaluation model.Evaluations) (int64, error) {
	createEvaluationStmt := table.Evaluations.INSERT(
		table.Evaluations.EvalStatusID,
		table.Evaluations.EvalTotalScore,
		table.Evaluations.EvalPossibleScore,
		table.Evaluations.TestRuntimeStatisticsID,
		table.Evaluations.CompilationDataID,
		table.Evaluations.TaskVersionID).MODEL(evaluation).
		RETURNING(table.Evaluations.ID)
	evaluationDest := &model.Evaluations{}
	err := createEvaluationStmt.Query(tx, evaluationDest)
	return evaluationDest.ID, err
}
