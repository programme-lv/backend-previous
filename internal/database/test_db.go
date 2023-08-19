package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDBProvider provides a migrated database
// and a method to close it (don't forget to do it).
type TestDBProvider interface {
	GetTestDB() *sqlx.DB
	Close()
}

func NewPostgresContainerTestDBProvider() (TestDBProvider, error) {
	return initPostgresContainerTestDB()
}

type postgresContainerTestDB struct {
	pgContainer *postgresContainer
	sqlxDb      *sqlx.DB
}

func initPostgresContainerTestDB() (*postgresContainerTestDB, error) {
	res := &postgresContainerTestDB{}

	pgContainer, err := newPostgresContainer()
	if err != nil {
		return nil, err
	}
	res.pgContainer = pgContainer

	pgContainerHost, err := pgContainer.container.Host(context.Background())
	if err != nil {
		return nil, err
	}

	pgContainerPort, err := pgContainer.container.MappedPort(context.Background(), "5432")
	if err != nil {
		return nil, err
	}

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgContainerHost, pgContainerPort, pgContainer.user, pgContainer.password, pgContainer.database)

	sqlxDb := sqlx.MustConnect("postgres", connString)

	res.sqlxDb = sqlxDb

	return res, nil
}

func (ptdb *postgresContainerTestDB) GetTestDB() *sqlx.DB {
	return ptdb.sqlxDb
}

func (ptdb *postgresContainerTestDB) Close() {
	err := ptdb.sqlxDb.Close()
	if err != nil {
		log.Printf("Failed to close sqlx db: %v", err)
	}

	err = ptdb.pgContainer.container.Terminate(context.Background())
	if err != nil {
		log.Printf("Failed to terminate container: %v", err)
	}
}

type postgresContainer struct {
	container testcontainers.Container
	user      string
	password  string
	database  string
}

func newPostgresContainer() (*postgresContainer, error) {
	DB_USER := "proglv"
	DB_PASS := "proglv"
	DB_NAME := "proglv"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     DB_USER,
			"POSTGRES_PASSWORD": DB_PASS,
			"POSTGRES_DB":       DB_NAME,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}

	res := &postgresContainer{
		container: container,
		user:      DB_USER,
		password:  DB_PASS,
		database:  DB_NAME,
	}
	return res, nil
}
