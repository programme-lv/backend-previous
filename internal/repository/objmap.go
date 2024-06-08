package repository

import (
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/domain"
)

func mapUserRecordToDomainObject(record model.Users) *domain.User {
	return &domain.User{
		ID:             record.ID,
		Username:       record.Username,
		Email:          record.Email,
		FirstName:      record.FirstName,
		LastName:       record.LastName,
		IsAdmin:        record.IsAdmin,
		HashedPassword: []byte(record.HashedPassword),
	}
}
