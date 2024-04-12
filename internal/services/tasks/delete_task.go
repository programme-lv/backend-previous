package tasks

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func DeleteTask(db qrm.DB, taskID int64) error {
	stmt := table.Tasks.UPDATE(table.Tasks.DeletedAt).SET(postgres.NOW()).
		WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))

	_, err := stmt.Exec(db)
	return err
}
