package tasks

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
)

func GetStableTaskVersionByTaskID(db qrm.DB, taskID int64) (*objects.TaskVersion, error) {
	stmt := postgres.SELECT(table.TaskVersions.AllColumns).FROM(
		table.Tasks.INNER_JOIN(table.TaskVersions, table.TaskVersions.ID.EQ(table.Tasks.StableVersionID)),
	).WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))

	var stableTaskVersion model.TaskVersions
	err := stmt.Query(db, &stableTaskVersion)
	if err != nil {
		return nil, err
	}

	descriptionObj, err := GetLVTaskVersionDescription(db, stableTaskVersion.ID)
	if err != nil {
		return nil, err
	}

	taskVersionObj := objects.TaskVersion{
		ID:            stableTaskVersion.ID,
		TaskID:        stableTaskVersion.TaskID,
		Code:          stableTaskVersion.ShortCode,
		Name:          stableTaskVersion.FullName,
		Description:   descriptionObj,
		TimeLimitMs:   stableTaskVersion.TimeLimMs,
		MemoryLimitKb: stableTaskVersion.MemLimKibibytes,
		CreatedAt:     stableTaskVersion.CreatedAt,
		UpdatedAt:     stableTaskVersion.UpdatedAt,
	}

	return &taskVersionObj, nil
}
