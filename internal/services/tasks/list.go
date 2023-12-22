package tasks

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
)

func ListPublishedTaskVersions(db qrm.Queryable) ([]objects.TaskVersion, error) {
	tasks := make([]objects.TaskVersion, 0)

	stmt := postgres.SELECT(
		table.Tasks.ID,
		table.Tasks.CreatedAt,

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
					table.Tasks.PublishedVersionID)).
			LEFT_JOIN(table.MarkdownStatements,
				table.MarkdownStatements.TaskVersionID.EQ(
					table.Tasks.PublishedVersionID))).
		WHERE(table.MarkdownStatements.LangIso6391.EQ(
			postgres.String("lv")))

	var publishedTasks []struct {
		model.Tasks
		model.TaskVersions
		model.MarkdownStatements
	}
	err := stmt.Query(db, &publishedTasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
