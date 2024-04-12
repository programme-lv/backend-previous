package tasks

import (
	"context"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func CreateTask(db *sqlx.DB, name string, code string, userID int64) (int64, error) {
	t, err := db.BeginTxx(context.TODO(), nil)
	if err != nil {
		return 0, err
	}
	defer t.Rollback()

	createTaskStmt := table.Tasks.INSERT(
		table.Tasks.CreatedAt,
		table.Tasks.CreatedByID,
	).VALUES(postgres.NOW(), userID).RETURNING(table.Tasks.ID)

	var task model.Tasks
	err = createTaskStmt.Query(t, &task)
	if err != nil {
		return 0, err
	}

	createTaskVersionStmt := table.TaskVersions.INSERT(
		table.TaskVersions.TaskID,
		table.TaskVersions.ShortCode,
		table.TaskVersions.FullName,
		table.TaskVersions.TimeLimMs,
		table.TaskVersions.MemLimKibibytes,
		table.TaskVersions.TestingTypeID,
		table.TaskVersions.CreatedAt,
	).VALUES(
		task.ID, code, name, 1000, 65536, "simple", postgres.NOW()).
		RETURNING(table.TaskVersions.ID)

	var taskVersion model.TaskVersions
	err = createTaskVersionStmt.Query(t, &taskVersion)
	if err != nil {
		return 0, err
	}

	updateTaskStmt := table.Tasks.UPDATE(
		table.Tasks.CurrentVersionID,
	).SET(
		taskVersion.ID,
	).WHERE(
		table.Tasks.ID.EQ(postgres.Int64(task.ID)),
	)

	_, err = updateTaskStmt.Exec(t)
	if err != nil {
		return 0, err
	}

	createMarkdownStmt := table.MarkdownStatements.INSERT(
		table.MarkdownStatements.Story,
		table.MarkdownStatements.Input,
		table.MarkdownStatements.Output,
		table.MarkdownStatements.TaskVersionID,
		table.MarkdownStatements.LangIso6391,
	).VALUES("", "", "", taskVersion.ID, "lv")

	_, err = createMarkdownStmt.Exec(t)
	if err != nil {
		return 0, err
	}

	err = t.Commit()
	if err != nil {
		return 0, err
	}

	return task.ID, nil
}
