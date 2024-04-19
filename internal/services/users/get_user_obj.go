package users

import (
	"github.com/go-jet/jet/qrm"
	"github.com/programme-lv/backend/internal/services/objects"
)

func GetUserObj(db qrm.DB, userID int64) (*objects.User, error) {
	userRecord, err := FindUserByID(db, userID)
	if err != nil {
		return nil, err
	}

	return &objects.User{
		ID:        userID,
		Username:  userRecord.Username,
		Email:     userRecord.Email,
		FirstName: userRecord.FirstName,
		LastName:  userRecord.LastName,
		IsAdmin:   userRecord.IsAdmin,
	}, nil
}
