package tasks

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func GetTaskIDByPublishedTaskCode(db qrm.DB, code string) (int64, error) {
	stmt := postgres.SELECT(table.PublishedTaskCodes.TaskID).
		FROM(table.PublishedTaskCodes).
		WHERE(table.PublishedTaskCodes.TaskCode.EQ(postgres.String(code)))
	var record model.PublishedTaskCodes
	err := stmt.Query(db, &record)
	if err != nil {
		return 0, err
	}
	return record.TaskID, nil
}
