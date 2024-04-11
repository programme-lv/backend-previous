package tasks

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
)

func GetTaskByTaskID(db qrm.DB, taskID int64) (*objects.Task, error) {
	stmt := postgres.SELECT(table.Tasks.AllColumns).FROM(table.Tasks).WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))

	var task model.Tasks
	err := stmt.Query(db, &task)
	if err != nil {
		return nil, err
	}

	var currTaskVers *objects.TaskVersion
	if task.CurrentVersionID != nil {
		currTaskVers, err = GetTaskVersionByTaskVersionID(db, *task.CurrentVersionID)
		if err != nil {
			return nil, err
		}
	}

	var stableTaskVers *objects.TaskVersion
	if task.StableVersionID != nil {
		stableTaskVers, err = GetTaskVersionByTaskVersionID(db, *task.StableVersionID)
		if err != nil {
			return nil, err
		}
	}

	taskObj := objects.Task{
		ID:          task.ID,
		CreatedByID: task.CreatedByID,
		Current:     currTaskVers,
		Stable:      stableTaskVers,
		CreatedAt:   task.CreatedAt,
	}

	return &taskObj, nil
}
