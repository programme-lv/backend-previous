package graphql

import (
	"github.com/programme-lv/backend/internal/database/dospaces"
	"github.com/programme-lv/backend/internal/eval"
	"github.com/programme-lv/backend/internal/lang"
	"github.com/programme-lv/backend/internal/task"
	"github.com/programme-lv/backend/internal/user"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	_ "github.com/lib/pq"
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
	Languages      lang.Service
	UserSrv        user.Service
	TaskSrv        task.Service
	SubmSrv        eval.Service
	SessionManager *scs.SessionManager
	Logger         *slog.Logger
	TestURLs       *dospaces.DOSpacesS3ObjStorage
	DirectorConn   *AuthDirectorConn
}
