package tasks

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/ztrue/tracerr"
)

func DuplicateTaskVersion(db qrm.DB, taskVersionID int64) (int64, error) {
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
		return 0, tracerr.Wrap(err)
	}

	// duplicate markdown statements which point to this task version
	stmt = table.MarkdownStatements.INSERT(
		table.MarkdownStatements.Story,
		table.MarkdownStatements.Input,
		table.MarkdownStatements.Output,
		table.MarkdownStatements.Notes,
		table.MarkdownStatements.Scoring,
		table.MarkdownStatements.TaskVersionID,
		table.MarkdownStatements.LangIso6391,
	).VALUES(
		postgres.SELECT(
			table.MarkdownStatements.Story,
			table.MarkdownStatements.Input,
			table.MarkdownStatements.Output,
			table.MarkdownStatements.Notes,
			table.MarkdownStatements.Scoring,
			postgres.Int64(taskVersion.ID),
			table.MarkdownStatements.LangIso6391,
		).FROM(table.MarkdownStatements).WHERE(table.MarkdownStatements.TaskVersionID.EQ(postgres.Int64(taskVersionID))),
	).RETURNING(table.MarkdownStatements.ID)

	var markdownStatements []model.MarkdownStatements
	err = stmt.Query(db, &markdownStatements)
	if err != nil {
		return 0, tracerr.Wrap(err)
	}

	// duplicate statement_examples which point to this task version
	stmt = table.StatementExamples.
		INSERT(table.StatementExamples.MutableColumns).
		VALUES(
			postgres.SELECT(
				table.StatementExamples.Input,
				table.StatementExamples.Answer,
				postgres.Int64(taskVersion.ID),
			).FROM(table.StatementExamples).WHERE(table.StatementExamples.TaskVersionID.EQ(postgres.Int64(taskVersionID))),
		).RETURNING(table.StatementExamples.ID)

	var statementExamples []model.StatementExamples
	err = stmt.Query(db, &statementExamples)
	if err != nil {
		return 0, tracerr.Wrap(err)
	}

	// duplicate task_version_tests which point to this task version
	stmt = table.TaskVersionTests.
		INSERT(table.TaskVersionTests.MutableColumns).
		VALUES(
			postgres.SELECT(
				table.TaskVersionTests.TestFilename,
				// table.TaskVersionTests.TaskVersionID,
				postgres.Int64(taskVersion.ID),
				table.TaskVersionTests.InputTextFileID,
				table.TaskVersionTests.AnswerTextFileID,
			).FROM(table.TaskVersionTests).WHERE(table.TaskVersionTests.TaskVersionID.EQ(postgres.Int64(taskVersionID))),
		).RETURNING(table.TaskVersionTests.ID)

	var taskVersionTests []model.TaskVersionTests
	err = stmt.Query(db, &taskVersionTests)
	if err != nil {
		return 0, tracerr.Wrap(err)
	}

	return taskVersion.ID, nil
}
