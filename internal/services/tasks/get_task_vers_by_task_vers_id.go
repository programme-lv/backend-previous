package tasks

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
)

func GetTaskVersionByTaskVersionID(db qrm.DB, taskVersionID int64) (*objects.TaskVersion, error) {
	stmt := postgres.SELECT(table.TaskVersions.AllColumns).FROM(
		table.TaskVersions).
		WHERE(table.TaskVersions.ID.EQ(postgres.Int64(taskVersionID)))

	var tv model.TaskVersions
	err := stmt.Query(db, &tv)
	if err != nil {
		return nil, err
	}

	descriptionObj, err := GetLVTaskVersionDescription(db, tv.ID)
	if err != nil {
		return nil, err
	}

	taskVersionObj := objects.TaskVersion{
		ID:            tv.ID,
		TaskID:        tv.TaskID,
		Code:          tv.ShortCode,
		Name:          tv.FullName,
		Description:   descriptionObj,
		TimeLimitMs:   tv.TimeLimMs,
		MemoryLimitKb: tv.MemLimKibibytes,
		CreatedAt:     tv.CreatedAt,
		UpdatedAt:     tv.UpdatedAt,
	}

	return &taskVersionObj, nil
}
