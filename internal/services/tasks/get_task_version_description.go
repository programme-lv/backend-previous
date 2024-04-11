package tasks

import (
	"errors"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
	"github.com/ztrue/tracerr"
)

func GetLVTaskVersionDescription(db qrm.DB, taskVersionID int64) (*objects.Description, error) {
	stmt := postgres.SELECT(table.MarkdownStatements.AllColumns).
		FROM(table.MarkdownStatements).
		WHERE(table.MarkdownStatements.TaskVersionID.EQ(postgres.Int64(taskVersionID))).LIMIT(1)

	var description model.MarkdownStatements
	err := stmt.Query(db, &description)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, tracerr.Wrap(err)
	}

	examplesStmt := postgres.SELECT(table.StatementExamples.AllColumns).
		FROM(table.StatementExamples).
		WHERE(table.StatementExamples.TaskVersionID.EQ(postgres.Int64(taskVersionID)))

	var examples []model.StatementExamples
	err = examplesStmt.Query(db, &examples)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var examplesObj []objects.Example
	for _, example := range examples {
		exampleObj := objects.Example{
			ID:     taskVersionID,
			Input:  example.Input,
			Answer: example.Answer,
		}
		examplesObj = append(examplesObj, exampleObj)
	}

	descriptionObj := objects.Description{
		ID:       description.ID,
		Story:    description.Story,
		Input:    description.Input,
		Output:   description.Output,
		Examples: examplesObj,
		Notes:    description.Notes,
	}

	return &descriptionObj, nil
}
