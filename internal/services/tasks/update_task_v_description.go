package tasks

import (
	"fmt"

	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

type UpdateTaskVStatementInput struct {
	Story  *string
	Input  *string
	Output *string
	Notes  *string
}

// duplicates task version and updates markdown statement, returns new task version id
func UpdateTaskVersionStatement(db *sqlx.DB, taskVersionID int64, input UpdateTaskVStatementInput) (int64, error) {
	curMarkdownStmts, err := selectMarkdownStatements(db, taskVersionID)
	if err != nil {
		return 0, err
	}

	if len(curMarkdownStmts) != 1 {
		return 0, fmt.Errorf("expected 1 markdown statement, got %d", len(curMarkdownStmts))
	}

	curMarkdownStmt := curMarkdownStmts[0]

	if !compareMarkdownStatementToInput(&curMarkdownStmt, input) {
		return 0, nil
	}

	tx, err := db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	newTaskVersID, err := DuplicateTaskVersion(tx, taskVersionID)
	if err != nil {
		return 0, err
	}

	oldMarkdowns, err := selectMarkdownStatements(tx, newTaskVersID)
	if err != nil {
		return 0, err
	}

	if len(oldMarkdowns) != 1 {
		return 0, fmt.Errorf("expected 1 markdown statement, got %d", len(oldMarkdowns))
	}

	oldMarkdown := oldMarkdowns[0]

	if input.Story != nil {
		oldMarkdown.Story = *input.Story
	}
	if input.Input != nil {
		oldMarkdown.Input = *input.Input
	}
	if input.Output != nil {
		oldMarkdown.Output = *input.Output
	}
	if input.Notes != nil {
		oldMarkdown.Notes = input.Notes
	}

	updStmt := table.MarkdownStatements.
		UPDATE(table.MarkdownStatements.MutableColumns).
		MODEL(oldMarkdown).
		WHERE(table.MarkdownStatements.ID.EQ(postgres.Int64(oldMarkdown.ID)))

	_, err = updStmt.Exec(tx)
	if err != nil {
		return 0, err
	}

	return newTaskVersID, tx.Commit()
}

// returns true if the markdown statement should be updated
func compareMarkdownStatementToInput(markdown *model.MarkdownStatements, input UpdateTaskVStatementInput) bool {
	if input.Story != nil && *input.Story != markdown.Story {
		return true
	}

	if input.Input != nil && *input.Input != markdown.Input {
		return true
	}

	if input.Output != nil && *input.Output != markdown.Output {
		return true
	}

	if input.Notes != nil && (markdown.Notes == nil || *input.Notes != *markdown.Notes) {
		return true
	}

	return false
}

func selectMarkdownStatements(db qrm.DB, taskVersionID int64) ([]model.MarkdownStatements, error) {
	stmt := postgres.SELECT(table.MarkdownStatements.AllColumns).
		FROM(table.MarkdownStatements).
		WHERE(table.MarkdownStatements.TaskVersionID.EQ(postgres.Int64(taskVersionID)))

	var markdowns []model.MarkdownStatements
	err := stmt.Query(db, &markdowns)
	if err != nil {
		return nil, err
	}

	return markdowns, nil
}
