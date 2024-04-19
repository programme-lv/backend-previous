package tasks

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
	"github.com/ztrue/tracerr"
)

func GetTaskVersionObjByTaskVersionID(db qrm.DB, taskVersionID int64,
	fillDescription bool) (*objects.TaskVersion, error) {

	stmt := postgres.SELECT(table.TaskVersions.AllColumns).FROM(
		table.TaskVersions).
		WHERE(table.TaskVersions.ID.EQ(postgres.Int64(taskVersionID)))

	var tv model.TaskVersions
	err := stmt.Query(db, &tv)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var descriptionObj *objects.Description = nil

	if fillDescription {
		descriptionObj, err = GetTaskVersionDescriptionObj(db, tv.ID)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
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
	}

	return &taskVersionObj, nil
}
