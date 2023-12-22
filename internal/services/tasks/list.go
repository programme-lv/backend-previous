package tasks

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
	"github.com/programme-lv/backend/internal/services/objects"
)

func ListPublishedTaskVersions(db qrm.Queryable) ([]objects.TaskVersion, error) {
	res := make([]objects.TaskVersion, 0)

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

	var publishedTaskVersions []struct {
		model.Tasks
		model.TaskVersions
		model.MarkdownStatements
	}
	err := stmt.Query(db, &publishedTaskVersions)
	if err != nil {
		return nil, err
	}

	for _, version := range publishedTaskVersions {
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
				Examples: nil,
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
