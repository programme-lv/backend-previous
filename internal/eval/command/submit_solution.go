package command

import (
	"context"
	"github.com/google/uuid"
	"github.com/programme-lv/backend/internal/common/decorator"
	"github.com/programme-lv/backend/internal/common/logs"
	"github.com/programme-lv/backend/internal/eval"
	"log"
	"log/slog"
)

type SubmitSolution struct {
	SubmissionUUID uuid.UUID
	//TaskCode   string
	TaskID     int64
	AuthorID   int64
	ProgLangID string
	Submission string
}

type SubmitSolutionHandler decorator.CommandHandler[SubmitSolution]

type submitSolutionHandler struct {
	repo eval.Repository
}

func NewSubmitSolutionHandler(
	repo eval.Repository,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) SubmitSolutionHandler {
	return decorator.ApplyCommandDecorators[SubmitSolution](
		submitSolutionHandler{
			repo: repo,
		},
		logger,
		metricsClient,
	)
}

func (h submitSolutionHandler) Handle(ctx context.Context, cmd SubmitSolution) (err error) {
	defer func() {
		logs.LogCommandExecution("SubmitSolution", cmd, err)
	}()

	subm, err := eval.NewSubmission(cmd.SubmissionUUID, cmd.TaskID, cmd.AuthorID, cmd.ProgLangID, cmd.Submission)
	if err != nil {
		return err
	}

	err = h.repo.AddSubmission(ctx, *subm)
	if err != nil {
		return err
	}

	log.Println("submitting solution")

	return nil
}
