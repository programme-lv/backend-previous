package tasks

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
)

func GetCurrentTaskVersionByTaskID(db qrm.DB, taskID int64) (*objects.TaskVersion, error) {
	stmt := postgres.SELECT(table.TaskVersions.AllColumns).FROM(
		table.Tasks.INNER_JOIN(table.TaskVersions, table.TaskVersions.ID.EQ(table.Tasks.CurrentVersionID)),
	).WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))

	var taskVersion model.TaskVersions
	err := stmt.Query(db, &taskVersion)
	if err != nil {
		return nil, err
	}

	descriptionObj, err := GetLVTaskVersionDescription(db, taskVersion.ID)
	if err != nil {
		return nil, err
	}

	taskVersionObj := objects.TaskVersion{
		ID:            taskVersion.ID,
		TaskID:        taskVersion.TaskID,
		Code:          taskVersion.ShortCode,
		Name:          taskVersion.FullName,
		Description:   descriptionObj,
		TimeLimitMs:   taskVersion.TimeLimMs,
		MemoryLimitKb: taskVersion.MemLimKibibytes,
		CreatedAt:     taskVersion.CreatedAt,
		UpdatedAt:     taskVersion.UpdatedAt,
	}

	return &taskVersionObj, nil
}
