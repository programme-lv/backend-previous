package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	pgContainerHost, err := pgContainer.dockerContainer.Host(context.Background())
	if err != nil {
		return nil, err
	}

	pgContainerPort, err := pgContainer.dockerContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		return nil, err
	}

	// strip /tcp suffix from port
	if pgContainerPort[len(pgContainerPort)-4:] == "/tcp" {
		pgContainerPort = pgContainerPort[:len(pgContainerPort)-4]
	}

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgContainerHost, pgContainerPort, pgContainer.user, pgContainer.password, pgContainer.database)

	log.Printf("connString: %s", connString)

	sqlxDb := sqlx.MustConnect("postgres", connString)

	res.sqlxDb = sqlxDb

	// TODO: apply the migrations here
	// start flyway container passing
	// host as the network
	// checkout git database to tmp dir
	// or pull git submodule

	return res, nil
}

type DBMigrations struct {
	rootDir string
}

func cloneDBMigrations() (*DBMigrations, error) {
	tmpDir, err := os.MkdirTemp("", "proglv-db-migrations")
	if err != nil {
		return nil, err
	}
	repoUrl := "https://github.com/programme-lv/database"

	_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	res := &DBMigrations{
		rootDir: tmpDir,
	}
	return res, nil
}

func (dbm *DBMigrations) getFlywayMigrationsDir() string {
	return filepath.Join(dbm.rootDir, "flyway-migrations")
}

func (dbm *DBMigrations) Close() {
	err := os.RemoveAll(dbm.rootDir)
	if err != nil {
		log.Printf("Failed to remove tmp dir: %v", err)
	}
}

func (ptdb *postgresContainerTestDB) GetTestDB() *sqlx.DB {
	return ptdb.sqlxDb
}

func (ptdb *postgresContainerTestDB) Close() {
	err := ptdb.sqlxDb.Close()
	if err != nil {
		log.Printf("Failed to close sqlx db: %v", err)
	}

	err = ptdb.pgContainer.dockerContainer.Terminate(context.Background())
	if err != nil {
		log.Printf("Failed to terminate container: %v", err)
	}
}

type postgresContainer struct {
	dockerContainer testcontainers.Container
	user            string
	password        string
	database        string
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
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForExposedPort(),
		),
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
		dockerContainer: container,
		user:            DB_USER,
		password:        DB_PASS,
		database:        DB_NAME,
	}
	return res, nil
}
