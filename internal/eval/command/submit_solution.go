package command

import (
	"context"
	"github.com/programme-lv/backend/internal/common/decorator"
	"github.com/programme-lv/backend/internal/common/logs"
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
}

func NewSubmitSolutionHandler(
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) SubmitSolutionHandler {
	return decorator.ApplyCommandDecorators[SubmitSolution](
		submitSolutionHandler{},
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
