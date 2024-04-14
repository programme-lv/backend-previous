package tasks

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func CreateTask(db qrm.DB, name string, code string, userID int64) (int64, error) {
	createTaskStmt := table.Tasks.INSERT(
		table.Tasks.CreatedAt,
		table.Tasks.CreatedByID,
	).VALUES(postgres.NOW(), userID).RETURNING(table.Tasks.ID)

	var task model.Tasks
	err := createTaskStmt.Query(db, &task)
	if err != nil {
		return 0, err
	}

	taskVersID, err := createTaskVersion(db, task.ID, code, name)
	if err != nil {
		return 0, err
	}

	err = assignCurrentVersionToTask(db, task.ID, taskVersID)
	if err != nil {
		return 0, err
	}

	mdStmtID, err := createDefaultMarkdownStatement(db)
	if err != nil {
		return 0, err
	}

	err = assignMdStatementToTaskVersion(db, mdStmtID, taskVersID)
	if err != nil {
		return 0, err
	}

	return task.ID, nil
}

func createTaskVersion(db qrm.Queryable, taskID int64, code, name string) (int64, error) {
	insertStmt := table.TaskVersions.INSERT(
		table.TaskVersions.TaskID,
		table.TaskVersions.ShortCode,
		table.TaskVersions.FullName,
		table.TaskVersions.TimeLimMs,
		table.TaskVersions.MemLimKibibytes,
		table.TaskVersions.TestingTypeID,
		table.TaskVersions.CreatedAt,
	).VALUES(
		taskID,
		code,
		name,
		1000,
		65536,
		"simple",
		postgres.NOW(),
	).RETURNING(table.TaskVersions.ID)

	var taskVersion model.TaskVersions
	err := insertStmt.Query(db, &taskVersion)
	if err != nil {
		return 0, err
	}

	return taskVersion.ID, nil
}

func assignCurrentVersionToTask(db qrm.Executable, taskID, versionID int64) error {
	updateStmt := table.Tasks.UPDATE(table.Tasks.CurrentVersionID).
		SET(versionID).
		WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))

	_, err := updateStmt.Exec(db)
	return err
}

func createDefaultMarkdownStatement(db qrm.Queryable) (int64, error) {
	insertStmt := table.MarkdownStatements.INSERT(
		table.MarkdownStatements.Story,
		table.MarkdownStatements.Input,
		table.MarkdownStatements.Output,
		table.MarkdownStatements.LangIso6391,
	).VALUES("", "", "", "lv").RETURNING(table.MarkdownStatements.ID)

	var stmtRecord model.MarkdownStatements
	err := insertStmt.Query(db, &stmtRecord)
	if err != nil {
		return 0, err
	}

	return stmtRecord.ID, nil
}

func assignMdStatementToTaskVersion(db qrm.Executable, statementID, versionID int64) error {
	updateStmt := table.TaskVersions.UPDATE(table.TaskVersions.MdStatementID).
		SET(statementID).
		WHERE(table.TaskVersions.ID.EQ(postgres.Int64(int64(versionID))))

	_, err := updateStmt.Exec(db)
	return err
}
