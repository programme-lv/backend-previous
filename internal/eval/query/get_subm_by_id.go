package query

import (
	"context"
	"github.com/google/uuid"
	"github.com/programme-lv/backend/internal/common/decorator"
	"log/slog"
)

type GetSubmissionByUUID struct {
	UUID uuid.UUID
}

type GetSubmissionByIDHandler decorator.QueryHandler[GetSubmissionByUUID, *Submission]

type getSubmissionByIDHandler struct {
	readModel GetSubmissionByIDReadModel
}

func NewGetSubmissionByIDHandler(
	readModel GetSubmissionByIDReadModel,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetSubmissionByIDHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[GetSubmissionByUUID, *Submission](
		getSubmissionByIDHandler{readModel: readModel},
		logger,
		metricsClient,
	)
}

type GetSubmissionByIDReadModel interface {
	GetSubmissionByID(ctx context.Context, uuid uuid.UUID) (*Submission, error)
}

func (h getSubmissionByIDHandler) Handle(ctx context.Context, query GetSubmissionByUUID) (*Submission, error) {
	return h.readModel.GetSubmissionByID(ctx, query.UUID)
}
