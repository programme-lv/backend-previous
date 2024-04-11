package tasks

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func ListEditableTaskIDs(db qrm.DB, userID int64) ([]int64, error) {
	stmt := postgres.SELECT(table.Tasks.ID).FROM(table.Tasks).
		WHERE(table.Tasks.CreatedByID.EQ(postgres.Int64(userID)))

	var tasks []struct {
		ID int64
	}
	err := stmt.Query(db, &tasks)
	if err != nil {
		return nil, err
	}

	var taskIDs []int64
	for _, task := range tasks {
		taskIDs = append(taskIDs, task.ID)
	}

	return taskIDs, nil
}
