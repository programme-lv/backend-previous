package repository

import (
	"context"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal"
	"github.com/programme-lv/backend/internal/database/proglv/public/model"
	"github.com/programme-lv/backend/internal/domain"
	"time"
)

type UserRepoPostgreSQLImpl struct {
	pg *sqlx.DB
}

type UserRepoPostgreSQLTxImpl struct {
	UserRepoPostgreSQLImpl
	tx *sqlx.Tx
}

func (u UserRepoPostgreSQLTxImpl) Commit() error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepoPostgreSQLTxImpl) Rollback() error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepoPostgreSQLImpl) BeginTx(ctx context.Context) (internal.UserRepoTx, error) {
	tx, err := u.pg.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &UserRepoPostgreSQLTxImpl{tx: tx}, nil
}

func NewUserRepoPostgreSQLImpl(pg *sqlx.DB) *UserRepoPostgreSQLImpl {
	return &UserRepoPostgreSQLImpl{pg: pg}
}

var _ internal.UserRepo = &UserRepoPostgreSQLImpl{}

func (u UserRepoPostgreSQLImpl) DoesUserExistByUsername(username string) (bool, error) {
	stmt := postgres.SELECT(table.Users.ID).
		FROM(table.Users).
		WHERE(table.Users.Username.EQ(postgres.String(username))).
		LIMIT(1)

	var record model.Users
	err := stmt.Query(u.pg, &record)
	if err != nil {
		return false, err
	}

	return record.ID != 0, nil
}

func (u UserRepoPostgreSQLImpl) DoesUserExistByEmail(email string) (bool, error) {
	stmt := postgres.SELECT(table.Users.ID).
		FROM(table.Users).
		WHERE(table.Users.Email.EQ(postgres.String(email))).
		LIMIT(1)

	var record model.Users
	err := stmt.Query(u.pg, &record)
	if err != nil {
		return false, err
	}

	return record.ID != 0, nil
}

func (u UserRepoPostgreSQLImpl) CreateUser(username string, hashedPassword []byte, email, firstName, lastName string) (int64, error) {
	user := model.Users{
		ID:             0, // auto-assigned by postgres
		Username:       username,
		Email:          email,
		HashedPassword: string(hashedPassword),
		FirstName:      firstName,
		LastName:       lastName,
		CreatedAt:      time.Now(),
		UpdatedAt:      nil,   // hasn't been updated yet
		IsAdmin:        false, // default to false
	}
	stmt := table.Users.INSERT(table.Users.MutableColumns).MODEL(user).RETURNING(table.Users.ID)

	var record model.Users
	err := stmt.Query(u.pg, &record)
	if err != nil {
		return 0, err
	}

	return record.ID, nil
}

func (u UserRepoPostgreSQLImpl) GetUserByID(id int64) (*domain.User, error) {
	stmt := postgres.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.ID.EQ(postgres.Int64(id)))

	var record model.Users
	err := stmt.Query(u.pg, &record)
	if err != nil {
		return nil, err
	}

	return mapUserRecordToDomainObject(record), nil
}

func (u UserRepoPostgreSQLImpl) GetUserByUsername(username string) (*domain.User, error) {
	stmt := postgres.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.Username.EQ(postgres.String(username)))

	var record model.Users
	err := stmt.Query(u.pg, &record)
	if err != nil {
		return nil, err
	}

	return mapUserRecordToDomainObject(record), nil
}
