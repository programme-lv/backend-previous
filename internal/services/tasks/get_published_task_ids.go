package tasks

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func GetPublishedTaskIDs(db qrm.DB) ([]int64, error) {
	// published_task_codes (task_code, task_id)
	selectAllPublishedTaskIDs := postgres.SELECT(table.PublishedTaskCodes.TaskID).
		FROM(table.PublishedTaskCodes)

	var publishedTaskCodeRecords []model.PublishedTaskCodes
	err := selectAllPublishedTaskIDs.Query(db, &publishedTaskCodeRecords)
	if err != nil {
		return nil, err
	}

	publishedTaskIDs := make([]int64, 0, len(publishedTaskCodeRecords))
	for _, publishedTaskCodeRecord := range publishedTaskCodeRecords {
		publishedTaskIDs = append(publishedTaskIDs, publishedTaskCodeRecord.TaskID)
	}

	return publishedTaskIDs, nil
}
