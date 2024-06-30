package graphql

import (
	"context"
	"github.com/programme-lv/backend/internal/user"
)

func (r *Resolver) getUserFromContext(ctx context.Context) (*user.User, error) {
	userID := r.SessionManager.GetInt64(ctx, "user_id")
	if userID == 0 {
		r.Logger.Warn("Whoami query failed due to unauthorized user", "action", "whoami")
		return nil, newErrorUnauthorized()
	}

	user, err := r.UserSrv.GetUserByID(userID)
	if err != nil {
		r.Logger.Warn("Whoami query failed due to internal server error", "userID", userID, "error", err.Error(), "action", "whoami")
		return nil, smartError(ctx, err)
	}

	return user, nil
}
