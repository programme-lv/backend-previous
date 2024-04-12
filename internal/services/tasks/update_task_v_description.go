package tasks

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

type UpdateTaskVDescInput struct {
	Story    *string
	Input    *string
	Output   *string
	Examples []*UpdateTaskVExampleInput
	Notes    *string
}

type UpdateTaskVExampleInput struct {
	Input  *string
	Answer *string
}

func UpdateTaskVersionDescription(db qrm.DB, taskVersionID int64, input UpdateTaskVDescInput) error {
	// duplicate task version row
	stmt := table.TaskVersions.INSERT(
		table.TaskVersions.TaskID,
		table.TaskVersions.ShortCode,
		table.TaskVersions.FullName,
		table.TaskVersions.TimeLimMs,
		table.TaskVersions.MemLimKibibytes,
		table.TaskVersions.TestingTypeID,
		table.TaskVersions.Origin,
		table.TaskVersions.CreatedAt,
		table.TaskVersions.CheckerID,
		table.TaskVersions.InteractorID,
	).VALUES(
		postgres.SELECT(
			table.TaskVersions.TaskID,
			table.TaskVersions.ShortCode,
			table.TaskVersions.FullName,
			table.TaskVersions.TimeLimMs,
			table.TaskVersions.MemLimKibibytes,
			table.TaskVersions.TestingTypeID,
			table.TaskVersions.Origin,
			table.TaskVersions.CreatedAt,
			table.TaskVersions.CheckerID,
			table.TaskVersions.InteractorID,
		).FROM(table.TaskVersions).WHERE(table.TaskVersions.ID.EQ(postgres.Int64(taskVersionID))),
	).RETURNING(table.TaskVersions.ID)

	var taskVersion model.TaskVersions
	err := stmt.Query(db, &taskVersion)
	if err != nil {
		return err
	}

	// duplicate markdown statement
	// stmt = table.MarkdownStatements.INSERT(
	// 	table.MarkdownStatements.Story,
	// 	table.MarkdownStatements.Input,
	// 	table.MarkdownStatements.Output,
	// 	table.MarkdownStatements.Notes,
	// 	table.MarkdownStatements.Scoring,
	// 	table.MarkdownStatements.TaskVersionID,
	// 	table.MarkdownStatements.LangIso6391,
	// ).VALUES(
	// 	input.Story,
	// 	input.Input,
	// 	input.Output,
	// 	input.Notes,
	// 	nil,
	// 	taskVersion.ID,
	// 	"lv",
	// )
	return nil

}
