package graphql

import (
	"context"
	"fmt"

	"github.com/alexedwards/scs/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/programme-lv/backend/internal/database"
	"golang.org/x/exp/slog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB             *sqlx.DB
	SessionManager *scs.SessionManager
	Logger         *slog.Logger
}

func (r *Resolver) GetUserFromContext(ctx context.Context) (*database.User, error) {
	userId, ok := r.SessionManager.Get(ctx, "user_id").(int64)
	if !ok {
		return nil, fmt.Errorf("user is not logged in")
	}

	var user database.User
	err := r.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
