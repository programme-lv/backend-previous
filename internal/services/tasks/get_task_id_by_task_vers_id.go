package tasks

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func GetTaskIDByTaskVersionID(db qrm.DB, taskVersionID int64) (int64, error) {
	stmt := postgres.SELECT(table.TaskVersions.TaskID).FROM(table.TaskVersions).
		WHERE(table.TaskVersions.ID.EQ(postgres.Int64(taskVersionID)))

	var taskID model.TaskVersions
	err := stmt.Query(db, &taskID)
	if err != nil {
		return 0, err
	}

	return taskID.TaskID, nil
}
