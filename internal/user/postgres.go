package user

import (
	"github.com/go-jet/jet/qrm"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/programme-lv/backend/internal/database/postgres/proglv/public/model"
	"github.com/programme-lv/backend/internal/database/postgres/proglv/public/table"
	"github.com/programme-lv/backend/internal/domain"
	"time"
)

type userRepoPostgresImpl struct {
	db qrm.DB
}

func (u userRepoPostgresImpl) DoesUserExistByID(id int64) (bool, error) {
	stmt := postgres.SELECT(table.Users.ID).
		FROM(table.Users).
		WHERE(table.Users.ID.EQ(postgres.Int64(id))).
		LIMIT(1)

	var record model.Users
	err := stmt.Query(u.db, &record)
	if err != nil {
		return false, err
	}

	return record.ID != 0, nil
}

func (u userRepoPostgresImpl) DoesUserExistByUsername(username string) (bool, error) {
	stmt := postgres.SELECT(table.Users.ID).
		FROM(table.Users).
		WHERE(table.Users.Username.EQ(postgres.String(username))).
		LIMIT(1)

	var record model.Users
	err := stmt.Query(u.db, &record)
	if err != nil {
		return false, err
	}

	return record.ID != 0, nil
}

func (u userRepoPostgresImpl) DoesUserExistByEmail(email string) (bool, error) {
	stmt := postgres.SELECT(table.Users.ID).
		FROM(table.Users).
		WHERE(table.Users.Email.EQ(postgres.String(email))).
		LIMIT(1)

	var record model.Users
	err := stmt.Query(u.db, &record)
	if err != nil {
		return false, err
	}

	return record.ID != 0, nil
}

func (u userRepoPostgresImpl) CreateUser(username string, hashedPassword []byte, email, firstName, lastName string) (int64, error) {
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
	err := stmt.Query(u.db, &record)
	if err != nil {
		return 0, err
	}

	return record.ID, nil
}

func (u userRepoPostgresImpl) GetUserByID(id int64) (*domain.User, error) {
	stmt := postgres.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.ID.EQ(postgres.Int64(id)))

	var record model.Users
	err := stmt.Query(u.db, &record)
	if err != nil {
		return nil, err
	}

	return mapUserTableRowToUserDomainObject(record), nil
}

func (u userRepoPostgresImpl) GetUserByUsername(username string) (*domain.User, error) {
	stmt := postgres.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.Username.EQ(postgres.String(username)))

	var record model.Users
	err := stmt.Query(u.db, &record)
	if err != nil {
		return nil, err
	}

	return mapUserTableRowToUserDomainObject(record), nil
}

func newUserRepoPostgresImpl(db qrm.DB) *userRepoPostgresImpl {
	return &userRepoPostgresImpl{db: db}
}

var _ userRepo = &userRepoPostgresImpl{}

func mapUserTableRowToUserDomainObject(record model.Users) *domain.User {
	return &domain.User{
		ID:        record.ID,
		Username:  record.Username,
		Email:     record.Email,
		FirstName: record.FirstName,
		LastName:  record.LastName,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		IsAdmin:   record.IsAdmin,
	}
}