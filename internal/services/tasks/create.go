package tasks

import (
	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/services/objects"
)

func CreateTaskVersion(db *sqlx.DB, task objects.TaskVersion) (int64, error) {
	return 0, nil
}
