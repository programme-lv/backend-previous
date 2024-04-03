package graphql

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/programme-lv/backend/internal/services/submissions"
	"github.com/programme-lv/director/msg"
	amqp "github.com/rabbitmq/amqp091-go"
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
	DirectorClient msg.DirectorClient
	DirectorPasswd string
}
