package tasks

import (
	"fmt"

	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/ztrue/tracerr"
)

type UpdateTaskVStatementInput struct {
	Story  *string
	Input  *string
	Output *string
	Notes  *string
}

// duplicates task version and updates markdown statement, returns new task version id
func UpdateCurrentTaskVersionStatement(db qrm.DB, taskID int64, input UpdateTaskVStatementInput) error {
	// select markdown statement
	// compare markdown statement
	// if doesn't require update, return
	// insert new markdown statement
	// duplicate task version with the new markdown statement
	// assign new task version to task (current version)

	currTaskVersID, err := selectCurrentTaskVersionID(db, taskID)
	if err != nil {
		return err
	}

	statement, err := selectTaskVersionMdStatement(db, currTaskVersID)
	if err != nil {
		return err
	}

	found := checkIfDiffUpdateIfIs(statement, input)
	if !found {
		return nil
	}

	newStmtID, err := insertNewMdStatement(db, statement)
	if err != nil {
		return err
	}

	newTaskVersID, err := duplicateTaskVersionWithNewStmtID(db, currTaskVersID, newStmtID)
	if err != nil {
		return err
	}

	err = assignNewTaskVersionToTask(db, taskID, newTaskVersID)
	if err != nil {
		return err
	}

	return nil
}

func selectCurrentTaskVersionID(db qrm.DB, taskID int64) (int64, error) {
	selectStatement := postgres.SELECT(table.Tasks.CurrentVersionID).
		FROM(table.Tasks).
		WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))

	var task model.Tasks
	err := selectStatement.Query(db, &task)
	if err != nil {
		tracerr.Wrap(err)
		return 0, err
	}

	if task.CurrentVersionID == nil {
		return 0, fmt.Errorf("task %d has no current version", taskID)
	}

	return *task.CurrentVersionID, nil
}

func assignNewTaskVersionToTask(db qrm.DB, taskID int64, newTaskVersionID int64) error {
	updateStatement := table.Tasks.UPDATE(table.Tasks.CurrentVersionID).
		SET(postgres.Int64(newTaskVersionID)).
		WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))

	_, err := updateStatement.Exec(db)
	if err != nil {
		tracerr.Wrap(err)
		return err
	}

	return nil
}

func selectTaskVersionMdStatement(db qrm.DB, taskVersionID int64) (*model.MarkdownStatements, error) {
	selectStatement := postgres.SELECT(table.MarkdownStatements.AllColumns).
		FROM(table.TaskVersions.INNER_JOIN(table.MarkdownStatements, table.MarkdownStatements.ID.EQ(table.TaskVersions.MdStatementID))).
		WHERE(table.TaskVersions.ID.EQ(postgres.Int64(taskVersionID)))

	var statement model.MarkdownStatements
	err := selectStatement.Query(db, &statement)
	if err != nil {
		tracerr.Wrap(err)
		return nil, err
	}

	return &statement, nil
}

func insertNewMdStatement(db qrm.DB, stmt *model.MarkdownStatements) (int64, error) {
	insertStatement := table.MarkdownStatements.INSERT(table.MarkdownStatements.MutableColumns).
		MODEL(stmt).
		RETURNING(table.MarkdownStatements.ID)

	var mdStatementsRecord model.MarkdownStatements
	err := insertStatement.Query(db, &mdStatementsRecord)
	if err != nil {
		tracerr.Wrap(err)
		return 0, err
	}

	return mdStatementsRecord.ID, nil
}

func checkIfDiffUpdateIfIs(statement *model.MarkdownStatements, input UpdateTaskVStatementInput) bool {
	found := false

	if input.Story != nil {
		if *input.Story != statement.Story {
			found = true
			statement.Story = *input.Story
		}
	}

	if input.Input != nil {
		if *input.Input != statement.Input {
			found = true
			statement.Input = *input.Input
		}
	}

	if input.Output != nil {
		if *input.Output != statement.Output {
			found = true
			statement.Output = *input.Output
		}
	}

	if input.Notes != nil {
		if statement.Notes == nil || *input.Notes != *statement.Notes {
			found = true
			statement.Notes = input.Notes
		}
	}

	return found
}

func duplicateTaskVersionWithNewStmtID(db qrm.DB, taskVersionID int64, newStmtID int64) (int64, error) {
	// insert new task version
	insertStatement := table.TaskVersions.INSERT(
		table.TaskVersions.MutableColumns.
			Except(table.TaskVersions.MdStatementID).
			Except(table.TaskVersions.CreatedAt),
		table.TaskVersions.MdStatementID,
		table.TaskVersions.CreatedAt,
	).QUERY(
		postgres.SELECT(
			table.TaskVersions.MutableColumns.
				Except(table.TaskVersions.MdStatementID).
				Except(table.TaskVersions.CreatedAt),
			postgres.Int64(newStmtID).AS(table.TaskVersions.MdStatementID.Name()),
			postgres.NOW().AS(table.TaskVersions.CreatedAt.Name()),
		).FROM(table.TaskVersions).
			WHERE(table.TaskVersions.ID.EQ(postgres.Int64(taskVersionID))),
	).
		RETURNING(table.TaskVersions.ID)

	var taskVersion model.TaskVersions
	err := insertStatement.Query(db, &taskVersion)
	if err != nil {
		tracerr.Wrap(err)
		return 0, err
	}

	return taskVersion.ID, nil
}
