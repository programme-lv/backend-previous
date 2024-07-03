package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/common/metrics"
	"github.com/programme-lv/backend/internal/eval/adapters"
	"github.com/programme-lv/backend/internal/eval/process"
	"github.com/programme-lv/backend/internal/eval/query"
	"log/slog"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	SubmitSolution process.SubmitSolutionHandler
}

type Queries struct {
	AllSubmissions    query.AllSubmissionsHandler
	GetSubmissionByID query.GetSubmissionByIDHandler
}

func NewApplication(pgDB *sqlx.DB) Application {
	logger := slog.Default()

	postgresRepo := adapters.NewEvaluationPostgresRepo(pgDB)

	metricsClient := metrics.NoOp{}

	return Application{
		Commands: Commands{
			SubmitSolution: process.NewSubmitSolutionHandler(postgresRepo, logger, metricsClient),
		},
		Queries: Queries{
			AllSubmissions:    query.NewAllSubmissionsHandler(postgresRepo, logger, metricsClient),
			GetSubmissionByID: query.NewGetSubmissionByIDHandler(postgresRepo, logger, metricsClient),
		},
	}
}
