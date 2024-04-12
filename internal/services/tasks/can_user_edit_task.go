package tasks

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func CanUserEditTask(db qrm.DB, userID int64, taskID int64) (bool, error) {
	// user can edit a task if they are an admin or the task was created by them

	// check if user is an admin
	stmtAdmin := table.Users.SELECT(table.Users.IsAdmin).WHERE(table.Users.ID.EQ(postgres.Int64(userID)))
	var user model.Users
	err := stmtAdmin.Query(db, &user)
	if err != nil {
		return false, err
	}

	if user.IsAdmin {
		return true, nil
	}

	// check if user created the task
	stmtTask := table.Tasks.SELECT(table.Tasks.CreatedByID).WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))
	var task model.Tasks
	err = stmtTask.Query(db, &task)
	if err != nil {
		return false, err
	}

	if task.CreatedByID == userID {
		return true, nil
	}

	return false, nil
}
