package graphql

import (
	"context"
	"fmt"

	"github.com/programme-lv/backend/internal/database"
	"golang.org/x/crypto/bcrypt"
)

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

func (r *Resolver) HashPassword(password string) (string, error) {
	var hashedPassword []byte
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
