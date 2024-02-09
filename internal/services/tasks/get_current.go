package tasks

import (
	"fmt"
	"strconv"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
)

func GetCurrentTaskVersionByID(db qrm.Queryable, id string, userID int64) (*objects.TaskVersion, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	selectCurrentTaskVersionByCodeStmt := postgres.SELECT(
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
		WHERE(table.MarkdownStatements.LangIso6391.EQ(postgres.String("lv")).AND(table.Tasks.ID.EQ(postgres.Int64(idInt))).AND(table.Tasks.CreatedByID.EQ(postgres.Int64(userID))))
		

	var currentTaskVersion []struct {
		model.Tasks
		model.TaskVersions
		model.MarkdownStatements
	}
	err = selectCurrentTaskVersionByCodeStmt.Query(db, &currentTaskVersion)
	if err != nil {
		return nil, err
	}

	if len(currentTaskVersion) != 1 {
		return nil, fmt.Errorf("expected 1 task version, got %d", len(currentTaskVersion))
	}

	var taskVersionIDs []postgres.Expression
	for _, version := range currentTaskVersion {
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

	taskVersion := currentTaskVersion[0]

	res := &objects.TaskVersion{
		ID:     taskVersion.TaskVersions.ID,
		TaskID: taskVersion.Tasks.ID,
		Code:   taskVersion.TaskVersions.ShortCode,
		Name:   taskVersion.TaskVersions.FullName,
		Description: &objects.Description{
			ID:       taskVersion.MarkdownStatements.ID,
			Story:    taskVersion.MarkdownStatements.Story,
			Input:    taskVersion.MarkdownStatements.Input,
			Output:   taskVersion.MarkdownStatements.Output,
			Examples: examplesMap[taskVersion.TaskVersions.ID],
			Notes:    taskVersion.MarkdownStatements.Notes,
		},
		TimeLimitMs:   taskVersion.TaskVersions.TimeLimMs,
		MemoryLimitKb: taskVersion.TaskVersions.MemLimKibibytes,
		CreatedAt:     taskVersion.TaskVersions.CreatedAt,
		UpdatedAt:     nil,
	}

	return res, nil
}
