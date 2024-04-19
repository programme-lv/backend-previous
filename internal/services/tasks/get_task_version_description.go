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

func GetTaskVersionDescriptionObj(db qrm.DB, taskVersionID int64) (*objects.Description, error) {
	selectDescriptionStmt := postgres.SELECT(table.MarkdownStatements.AllColumns).
		FROM(table.TaskVersions.INNER_JOIN(table.MarkdownStatements, table.TaskVersions.MdStatementID.EQ(table.MarkdownStatements.ID))).
		WHERE(table.TaskVersions.ID.EQ(postgres.Int64(taskVersionID)))

	var description model.MarkdownStatements
	err := selectDescriptionStmt.Query(db, &description)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, tracerr.Wrap(err)
	}

	selectExampleSetIDStmt := postgres.SELECT(table.TaskVersions.ExampleSetID).
		FROM(table.TaskVersions).
		WHERE(table.TaskVersions.ID.EQ(postgres.Int64(taskVersionID)))

	var taskVersionRecord model.TaskVersions
	err = selectExampleSetIDStmt.Query(db, &taskVersionRecord)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var examplesObjList []objects.Example
	if taskVersionRecord.ExampleSetID != nil {

		exampleSetID := *taskVersionRecord.ExampleSetID

		examplesStmt := postgres.SELECT(table.StatementExamples.AllColumns).
			FROM(table.StatementExamples).
			WHERE(table.StatementExamples.ExampleSetID.EQ(postgres.Int64(exampleSetID)))

		var examples []model.StatementExamples
		err = examplesStmt.Query(db, &examples)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		for _, example := range examples {
			exampleObj := objects.Example{
				ID:     example.ID,
				Input:  example.Input,
				Answer: example.Answer,
			}
			examplesObjList = append(examplesObjList, exampleObj)
		}
	}

	descriptionObj := objects.Description{
		ID:       description.ID,
		Story:    description.Story,
		Input:    description.Input,
		Output:   description.Output,
		Examples: examplesObjList,
		Notes:    description.Notes,
	}

	return &descriptionObj, nil
}
