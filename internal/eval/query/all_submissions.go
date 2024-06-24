package query

import (
	"context"
	"github.com/programme-lv/backend/internal/common/decorator"
	"log/slog"
)

type AllSubmissions struct{}

type AllSubmissionsHandler decorator.QueryHandler[AllSubmissions, []Submission]

type allSubmissionsHandler struct {
	readModel AllSubmissionsReadModel
}

func NewAllSubmissionsHandler(
	readModel AllSubmissionsReadModel,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) AllSubmissionsHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[AllSubmissions, []Submission](
		allSubmissionsHandler{readModel: readModel},
		logger,
		metricsClient,
	)
}

type AllSubmissionsReadModel interface {
	AllSubmissions(ctx context.Context) ([]Submission, error)
}

func (h allSubmissionsHandler) Handle(ctx context.Context, _ AllSubmissions) ([]Submission, error) {
	return h.readModel.AllSubmissions(ctx)
}
