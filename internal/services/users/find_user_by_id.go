package users

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func FindUserByID(db qrm.DB, userID int64) (*model.Users, error) {
	stmt := postgres.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.ID.EQ(postgres.Int64(userID)))

	var record model.Users
	err := stmt.Query(db, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
