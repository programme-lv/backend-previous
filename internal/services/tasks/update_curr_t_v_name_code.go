package tasks

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/ztrue/tracerr"
)

func UpdateCurrentTaskVersionNameAndCode(db qrm.DB, taskID int64, name, code string) error {
	currTaskVersID, err := selectCurrentTaskVersionID(db, taskID)
	if err != nil {
		return err
	}

	cName, cCode, err := selectTaskVersionNameAndCode(db, currTaskVersID)
	if err != nil {
		return err
	}

	if name == cName && code == cCode {
		return nil // no need to update
	}

	newTaskVersID, err := duplicateTaskVersionWithNewNameAndCode(db, currTaskVersID, name, code)
	if err != nil {
		return err
	}

	err = assignNewTaskVersionToTask(db, taskID, newTaskVersID)
	if err != nil {
		return err
	}

	return nil
}

func duplicateTaskVersionWithNewNameAndCode(db qrm.DB, taskVersionID int64, name, code string) (int64, error) {
	// insert new task version
	insertStatement := table.TaskVersions.INSERT(
		table.TaskVersions.MutableColumns.
			Except(table.TaskVersions.FullName).
			Except(table.TaskVersions.ShortCode),
		table.TaskVersions.FullName,
		table.TaskVersions.ShortCode,
	).QUERY(
		postgres.SELECT(
			table.TaskVersions.MutableColumns.
				Except(table.TaskVersions.FullName).
				Except(table.TaskVersions.ShortCode),
			postgres.String(name).AS(table.TaskVersions.FullName.Name()),
			postgres.String(code).AS(table.TaskVersions.ShortCode.Name()),
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

/*
	func duplicateTaskVersionWithNewStmtID(db qrm.DB, taskVersionID int64, newStmtID int64) (int64, error) {
		// insert new task version
		insertStatement := table.TaskVersions.INSERT(
			table.TaskVersions.MutableColumns.Except(table.TaskVersions.MdStatementID),
			table.TaskVersions.MdStatementID,
		).QUERY(
			postgres.SELECT(
				table.TaskVersions.MutableColumns.Except(table.TaskVersions.MdStatementID),
				postgres.Int64(newStmtID).AS(table.TaskVersions.MdStatementID.Name()),
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
*/
func selectTaskVersionNameAndCode(db qrm.DB, taskVersID int64) (string, string, error) {
	selectStmt := postgres.SELECT(table.TaskVersions.FullName, table.TaskVersions.ShortCode).
		FROM(table.TaskVersions).
		WHERE(table.TaskVersions.ID.EQ(postgres.Int64(taskVersID)))

	var taskVersRecord model.TaskVersions
	err := selectStmt.Query(db, &taskVersRecord)
	if err != nil {
		return "", "", err
	}

	return taskVersRecord.FullName, taskVersRecord.ShortCode, nil
}
