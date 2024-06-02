package graphql

import (
	"fmt"
	"github.com/programme-lv/backend/internal/domain"
)

func mapDomainUserObjToGQLUserObj(user *domain.User) *User {
	if user == nil {
		return nil
	}

	return &User{
		ID:        fmt.Sprintf("%d", user.ID),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsAdmin:   user.IsAdmin,
	}
}
