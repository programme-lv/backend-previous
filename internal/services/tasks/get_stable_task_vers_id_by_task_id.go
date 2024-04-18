package tasks

import (
	"fmt"

	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func GetStableTaskVerssionIDByTaskID(db qrm.DB, taskID int64) (int64, error) {
	stmt := postgres.SELECT(table.Tasks.StableVersionID).
		FROM(table.Tasks).WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))

	var record model.Tasks
	err := stmt.Query(db, &record)
	if err != nil {
		return 0, err
	}

	if record.StableVersionID == nil {
		return 0, fmt.Errorf("task with ID %d has no stable version", taskID)
	}

	return *record.StableVersionID, nil
}
