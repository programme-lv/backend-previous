package graphql

import (
	"github.com/programme-lv/backend/internal/domain"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/programme-lv/backend/internal/dospaces"
	"github.com/programme-lv/director/msg"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type AuthDirectorConn struct {
	GRPCClient msg.DirectorClient
	Password   string
}

type Resolver struct {
	UserQuerySrv   domain.UserService
	PostgresDB     *sqlx.DB
	SessionManager *scs.SessionManager
	Logger         *slog.Logger
	// SubmissionRMQ  *amqp.Connection
	TestURLs     *dospaces.DOSpacesS3ObjStorage
	DirectorConn *AuthDirectorConn
}
