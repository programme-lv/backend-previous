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

	taskVersions, err := ListPublishedTaskVersions(tx)
	assert.Nilf(t, err, "Failed to list published task versions: %v", err)
	assert.Equal(t, 0, len(taskVersions), "Expected 0 tasks, got %d", len(taskVersions))

	target := initTargetTaskVersion()
	err = createTaskVersionTarget(tx, target)
	assert.Nilf(t, err, "Failed to create task version target: %v", err)

	taskVersions, err = ListPublishedTaskVersions(tx)
	assert.Nilf(t, err, "Failed to list published task versions: %v", err)
	assert.Equal(t, 1, len(taskVersions), "Expected 1 task, got %d", len(taskVersions))

	received := taskVersions[0]
	assert.Equal(t, target.Code, received.Code, "Expected code %s, got %s", target.Code, received.Code)
	assert.Equal(t, target.Name, received.Name, "Expected name %s, got %s", target.Name, received.Name)

	assert.Equal(t, target.Description.Story, received.Description.Story,
		"Expected story %s, got %s", target.Description.Story, received.Description.Story)
	assert.Equal(t, target.Description.Input, received.Description.Input,
		"Expected input %s, got %s", target.Description.Input, received.Description.Input)
	assert.Equal(t, target.Description.Output, received.Description.Output,
		"Expected output %s, got %s", target.Description.Output, received.Description.Output)
	assert.Equal(t, target.Description.Examples, received.Description.Examples,
		"Expected examples %v, got %v", target.Description.Examples, received.Description.Examples)
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

func createTaskVersionTarget(tx *sqlx.Tx, target objects.TaskVersion) error {

	userID, err := insertTempTestUser(tx)
	if err != nil {
		return fmt.Errorf("Failed to create temp user: %v", err)
	}

	taskID, err := insertTask(tx, userID)
	if err != nil {
		return fmt.Errorf("Failed to create task: %v", err)
	}

	taskVersionID, err := insertTaskVersion(tx, taskID)
	if err != nil {
		return fmt.Errorf("Failed to create task version: %v", err)
	}

	_, err = insertMarkdownStatement(tx, taskVersionID)
	if err != nil {
		return fmt.Errorf("Failed to create markdown statement: %v", err)
	}

	err = updateTaskRelevantAndPublishedVersionIds(tx, taskID, taskVersionID)
	if err != nil {
		return fmt.Errorf("Failed to update task relevant and published version ids: %v", err)
	}

	return nil
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
func insertTaskVersion(tx *sqlx.Tx, taskID int64) (int64, error) {
	taskVersion := model.TaskVersions{
		TaskID:          taskID,
		ShortCode:       "summa",
		FullName:        "Summa",
		TimeLimMs:       1000,
		MemLimKibibytes: 256000,
		TestingTypeID:   "simple",
	}
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

func insertMarkdownStatement(tx *sqlx.Tx, taskVersionID int64) (int64, error) {
	markdownStatement := model.MarkdownStatements{
		TaskVersionID: &taskVersionID,
		LangIso6391:   "lv",
		Story:         "Apraksts",
		Input:         "Ieeja",
		Output:        "Izeja",
	}
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

// Update the task's relevant and published version IDs to the given task version ID.
func updateTaskRelevantAndPublishedVersionIds(tx *sqlx.Tx, taskID, taskVersionID int64) error {
	updateTaskStmt := table.Tasks.UPDATE(
		table.Tasks.RelevantVersionID,
		table.Tasks.PublishedVersionID).SET(
		taskVersionID, taskVersionID).WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))
	_, err := updateTaskStmt.Exec(tx)
	return err
}
