package tasks

import "github.com/go-jet/jet/qrm"

type UpdateTaskVStatementInput struct {
	Story  *string
	Input  *string
	Output *string
	Notes  *string
}

// duplicates task version and updates markdown statement, returns new task version id
func UpdateTaskVersionStatement(db qrm.DB, taskVersionID int64, input UpdateTaskVStatementInput) (int64, error) {
	panic("not implemented")
}
