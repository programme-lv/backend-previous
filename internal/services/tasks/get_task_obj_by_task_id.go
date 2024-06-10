package tasks

import (
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/domain"
	"github.com/ztrue/tracerr"
)

// 0 - no version
// 1 - version without description
// 2 - full version
func GetTaskObjByTaskID(db qrm.DB, taskID int64, currVersDepth, stableVersDepth int) (*domain.Task, error) {
	stmt := postgres.SELECT(table.Tasks.AllColumns).FROM(table.Tasks).WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))

	var task model.Tasks
	err := stmt.Query(db, &task)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var currTaskVers *domain.TaskVersion = nil
	if task.CurrentVersionID != nil {
		if currVersDepth == 2 {
			currTaskVers, err = GetTaskVersionObjByTaskVersionID(db, *task.CurrentVersionID, true)
			if err != nil {
				return nil, tracerr.Wrap(err)
			}
		} else if currVersDepth == 1 {
			currTaskVers, err = GetTaskVersionObjByTaskVersionID(db, *task.CurrentVersionID, false)
			if err != nil {
				return nil, tracerr.Wrap(err)
			}
		}
	}

	var stableTaskVers *objects.TaskVersion = nil
	if task.StableVersionID != nil {
		if stableVersDepth == 2 {
			stableTaskVers, err = GetTaskVersionObjByTaskVersionID(db, *task.StableVersionID, true)
			if err != nil {
				return nil, tracerr.Wrap(err)
			}
		} else if stableVersDepth == 1 {
			stableTaskVers, err = GetTaskVersionObjByTaskVersionID(db, *task.StableVersionID, false)
			if err != nil {
				return nil, tracerr.Wrap(err)
			}
		}
	}

	taskObj := domain.Task{
		ID:        task.ID,
		OwnerID:   task.CreatedByID,
		Current:   currTaskVers,
		Stable:    stableTaskVers,
		CreatedAt: task.CreatedAt,
	}

	return &taskObj, nil
}

func getMaxCreatedAtOfTaskVersions(db qrm.DB, taskID int64) (*time.Time, error) {
	stmt := postgres.SELECT(postgres.MAX(table.TaskVersions.CreatedAt).AS("abc")).
		FROM(table.TaskVersions).
		WHERE(table.TaskVersions.TaskID.EQ(postgres.Int64(taskID)))

	var updatedAt struct {
		Abc *time.Time `alias:"abc"`
	}
	err := stmt.Query(db, &updatedAt)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return updatedAt.Abc, nil
}
