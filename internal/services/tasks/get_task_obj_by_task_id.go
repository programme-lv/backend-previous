package tasks

import (
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
	"github.com/ztrue/tracerr"
)

func GetTaskObjByTaskID(db qrm.DB, taskID int64) (*objects.Task, error) {
	stmt := postgres.SELECT(table.Tasks.AllColumns).FROM(table.Tasks).WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))

	var task model.Tasks
	err := stmt.Query(db, &task)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var currTaskVers *objects.TaskVersion = nil
	if task.CurrentVersionID != nil {
		currTaskVers, err = GetTaskVersionByTaskVersionID(db, *task.CurrentVersionID)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
	}

	var stableTaskVers *objects.TaskVersion = nil
	if task.StableVersionID != nil {
		stableTaskVers, err = GetTaskVersionByTaskVersionID(db, *task.StableVersionID)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
	}

	stmt = postgres.SELECT(postgres.MAX(table.TaskVersions.CreatedAt).AS("abc")).
		FROM(table.TaskVersions).
		WHERE(table.TaskVersions.TaskID.EQ(postgres.Int64(taskID)))

	var updatedAt struct {
		Abc *time.Time `alias:"abc"`
	}
	err = stmt.Query(db, &updatedAt)
	if err != nil {
		return nil, tracerr.Wrap(err)
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
