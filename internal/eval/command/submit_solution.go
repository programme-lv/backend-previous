package command

import (
	"context"
	"github.com/programme-lv/backend/internal/common/decorator"
	"github.com/programme-lv/backend/internal/common/logs"
	"github.com/programme-lv/backend/internal/eval"
	"log"
	"log/slog"
)

type SubmitSolution struct {
	TaskCode   string
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

	log.Println("submitting solution")

	return nil
}
