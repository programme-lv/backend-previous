package authsessions

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"github.com/programme-lv/backend/internal"
	"github.com/programme-lv/backend/internal/domain"
)

type AuthSessionManagerImpl struct {
	sessions scs.SessionManager
}

func (a AuthSessionManagerImpl) PutUserIDIntoCtx(ctx context.Context, userID int64) {
	a.sessions.Put(ctx, "user_id", userID)
}

func (a AuthSessionManagerImpl) GetUserIDFromCtx(ctx context.Context) (int64, error) {
	userID, ok := a.sessions.Get(ctx, "user_id").(int64)
	if !ok {
		return 0, domain.NewErrorNotLoggedIn()
	}
	return userID, nil
}

func (a AuthSessionManagerImpl) PopUserIDFromCtx(ctx context.Context) (int64, error) {
	userID, err := a.GetUserIDFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	a.sessions.Remove(ctx, "user_id")
	return userID, nil
}

var _ internal.AuthSessionManager = &AuthSessionManagerImpl{}
