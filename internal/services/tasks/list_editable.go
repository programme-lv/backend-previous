package tasks

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
)

func ListEditableTaskVersions(db qrm.Queryable, userID int64) ([]objects.TaskVersion, error) {
	res := make([]objects.TaskVersion, 0)

	selectEditableTaskVersionsStmt := postgres.SELECT(
		table.Tasks.ID,
		table.Tasks.CreatedAt,

		table.TaskVersions.ID,
		table.TaskVersions.TaskID,
		table.TaskVersions.ShortCode,
		table.TaskVersions.FullName,
		table.TaskVersions.TimeLimMs,
		table.TaskVersions.MemLimKibibytes,
		table.TaskVersions.CreatedAt,
		

		table.MarkdownStatements.ID,
		table.MarkdownStatements.Story,
		table.MarkdownStatements.Input,
		table.MarkdownStatements.Output,
		table.MarkdownStatements.Notes).
		FROM(table.Tasks.
			INNER_JOIN(table.TaskVersions,
				table.TaskVersions.ID.EQ(
					table.Tasks.RelevantVersionID)).
			LEFT_JOIN(table.MarkdownStatements,
				table.MarkdownStatements.TaskVersionID.EQ(
					table.Tasks.RelevantVersionID))).
		WHERE(table.Tasks.CreatedByID.EQ(
				postgres.Int64(userID))).
		ORDER_BY(table.TaskVersions.CreatedAt.DESC());

	var editableTaskVersions []struct {
		model.Tasks
		model.TaskVersions
		model.MarkdownStatements
	}
	err := selectEditableTaskVersionsStmt.Query(db, &editableTaskVersions)
	if err != nil {
		return nil, err
	}

	if len(editableTaskVersions) == 0 {
		return res, nil
	}

	var taskVersionIDs []postgres.Expression
	for _, version := range editableTaskVersions {
		taskVersionIDs = append(taskVersionIDs, postgres.Int64(version.TaskVersions.ID))
	}

	selectExamplesStmt := postgres.SELECT(
		table.StatementExamples.ID,
		table.StatementExamples.Input,
		table.StatementExamples.Answer,
		table.StatementExamples.TaskVersionID).
		FROM(table.StatementExamples).
		WHERE(table.StatementExamples.TaskVersionID.IN(taskVersionIDs...))

	var examples []struct {
		model.StatementExamples
	}

	err = selectExamplesStmt.Query(db, &examples)
	if err != nil {
		return nil, err
	}

	examplesMap := make(map[int64][]objects.Example)
	for _, example := range examples {
		if _, ok := examplesMap[example.StatementExamples.TaskVersionID]; !ok {
			examplesMap[example.StatementExamples.TaskVersionID] = make([]objects.Example, 0)
		}
		examplesMap[example.StatementExamples.TaskVersionID] = append(
			examplesMap[example.StatementExamples.TaskVersionID],
			objects.Example{
				ID:     example.StatementExamples.ID,
				Input:  example.StatementExamples.Input,
				Answer: example.StatementExamples.Answer,
			})
	}

	for _, version := range editableTaskVersions {
		res = append(res, objects.TaskVersion{
			ID:     version.TaskVersions.ID,
			TaskID: version.Tasks.ID,
			Code:   version.TaskVersions.ShortCode,
			Name:   version.TaskVersions.FullName,
			Description: &objects.Description{
				ID:       version.MarkdownStatements.ID,
				Story:    version.MarkdownStatements.Story,
				Input:    version.MarkdownStatements.Input,
				Output:   version.MarkdownStatements.Output,
				Examples: examplesMap[version.TaskVersions.ID],
				Notes:    version.MarkdownStatements.Notes,
			},
			TimeLimitMs:   version.TaskVersions.TimeLimMs,
			MemoryLimitKb: version.TaskVersions.MemLimKibibytes,
			CreatedAt:     version.TaskVersions.CreatedAt,
			UpdatedAt:     nil,
		})
	}

	return res, nil
}
