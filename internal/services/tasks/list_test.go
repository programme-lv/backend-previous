package tasks

import (
	"context"
	"os"
	"testing"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/database/testdb"
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
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	tasks, err := ListPublishedTaskVersions(tx)
	if err != nil {
		t.Fatalf("Failed to list published task versions: %v", err)
	}

	if len(tasks) != 0 {
		t.Fatalf("Expected 0 tasks, got %d", len(tasks))
	}

	userID, err := createTempTestUser(tx)
	if err != nil {
		t.Fatalf("Failed to create temp user: %v", err)
	}

	taskID, err := createTask(tx, userID)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	taskVersionID, err := createTaskVersion(tx, taskID)
	if err != nil {
		t.Fatalf("Failed to create task version: %v", err)
	}

	err = updateTaskRelevantAndPublishedVersionIds(tx, taskID, taskVersionID)
	if err != nil {
		t.Fatalf("Failed to update task relevant and published version ids: %v", err)
	}

	tasks, err = ListPublishedTaskVersions(tx)
	if err != nil {
		t.Fatalf("Failed to list published task versions: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}
}

// Create a temporary test user and return its ID.
func createTempTestUser(tx *sqlx.Tx) (int64, error) {
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
func createTask(tx *sqlx.Tx, userID int64) (int64, error) {
	createTaskStmt := table.Tasks.INSERT(
		table.Tasks.CreatedByID).VALUES(userID).RETURNING(table.Tasks.ID)
	taskDest := &model.Tasks{}
	err := createTaskStmt.Query(tx, taskDest)
	return taskDest.ID, err
}

// Create a task version "summa" for the given task and return its ID.
func createTaskVersion(tx *sqlx.Tx, taskID int64) (int64, error) {
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

// Update the task's relevant and published version IDs to the given task version ID.
func updateTaskRelevantAndPublishedVersionIds(tx *sqlx.Tx, taskID, taskVersionID int64) error {
	updateTaskStmt := table.Tasks.UPDATE(
		table.Tasks.RelevantVersionID,
		table.Tasks.PublishedVersionID).SET(
		taskVersionID, taskVersionID).WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))
	_, err := updateTaskStmt.Exec(tx)
	return err
}
