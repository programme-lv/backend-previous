package graphql

import (
	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/programme-lv/backend/internal/services/submissions"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/exp/slog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostgresDB     *sqlx.DB
	SessionManager *scs.SessionManager
	Logger         *slog.Logger
	SubmissionRMQ  *amqp.Connection
	TestURLs       *submissions.S3TestURLs
}
