package postgres

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/database/postgres/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/postgres/proglv/public/table"
	"github.com/programme-lv/backend/internal/domain"
	"github.com/programme-lv/backend/internal/user"
)

type postgresUserRepoImpl struct {
	sqlxDB *sqlx.DB
	sqlxTx *sqlx.Tx
	db     qrm.DB
}

func NewPostgresUserRepo(db *sqlx.DB) user.Repo {
	return &postgresUserRepoImpl{
		sqlxDB: db,
		db:     db,
	}
}

func (p postgresUserRepoImpl) DoesUserExistByUsername(username string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresUserRepoImpl) DoesUserExistByEmail(email string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresUserRepoImpl) CreateUser(username string, hashedPassword []byte, email, firstName, lastName string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresUserRepoImpl) GetUserByID(id int64) (*domain.User, error) {
	stmt := postgres.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.ID.EQ(postgres.Int64(userID)))

	var record model.Users
	err := stmt.Query(p.db, &record)
	if err != nil {
		return nil, err
	}
	return mapUserTableRecordToDomainObject(record), nil
}

func (p postgresUserRepoImpl) GetUserByUsername(username string) (*domain.User, error) {
	selectStmt := postgres.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.Username.EQ(postgres.String(username)))

	var record model.Users
	err := selectStmt.Query(p.db, &record)
	if err != nil {
		return nil, err
	}

	return mapUserTableRecordToDomainObject(record), nil
}

func (p postgresUserRepoImpl) BeginTx() (user.RepoTx, error) {
	tx, err := p.sqlxDB.Beginx()
	if err != nil {
		return nil, err
	}
	return &postgresUserRepoImpl{
		sqlxDB: nil,
		sqlxTx: tx,
		db:     tx,
	}, nil
}

func (p postgresUserRepoImpl) Commit() error {
	return p.sqlxTx.Commit()
}

func (p postgresUserRepoImpl) Rollback() error {
	return p.sqlxTx.Rollback()
}

var _ user.Repo = &postgresUserRepoImpl{}
var _ user.RepoTx = &postgresUserRepoImpl{}

func mapUserTableRecordToDomainObject(record model.Users) *domain.User {
	return &domain.User{
		ID:        record.ID,
		Username:  record.Username,
		Email:     record.Email,
		FirstName: record.FirstName,
		LastName:  record.LastName,
	}
}
