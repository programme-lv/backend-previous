package users

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/proglv/public/table"
)

func GetUserIDByUsername(db qrm.DB, username string) (int64, error) {
	selectStmt := postgres.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.Username.EQ(postgres.String(username)))

	var record model.Users
	err := selectStmt.Query(db, &record)
	if err != nil {
		return 0, err
	}

	return record.ID, nil
}
