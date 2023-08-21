package testdb

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// TestDBProvider provides a migrated database
// and a method to close it (don't forget to do it).
type TestDBProvider interface {
	GetTestDB() *sqlx.DB
	Close()
}

func NewPostgresTestcontainerProvider() (TestDBProvider, error) {
	return initPostgresContainerTestDB()
}

type postgresContainerTestDB struct {
	pgContainer *postgresContainer
	sqlxDb      *sqlx.DB
}

func initPostgresContainerTestDB() (*postgresContainerTestDB, error) {
	res := &postgresContainerTestDB{}

	pgContainer, err := startPostgresContainer("proglv", "proglv", "proglv")
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

	// strip /tcp suffix from port
	if pgContainerPort[len(pgContainerPort)-4:] == "/tcp" {
		pgContainerPort = pgContainerPort[:len(pgContainerPort)-4]
	}

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgContainerHost, pgContainerPort, pgContainer.user, pgContainer.password, pgContainer.database)

	log.Printf("connString: %s", connString)

	sqlxDb := sqlx.MustConnect("postgres", connString)

	res.sqlxDb = sqlxDb

	migrations, err := cloneDBMigrations()
	if err != nil {
		log.Printf("Failed to clone migrations: %v", err)
	}
	defer migrations.erase()

	err = execFlywayContainer(migrations.getFlywayMigrationsDir(), pgContainerHost, pgContainerPort.Port(), pgContainer.database, pgContainer.user, pgContainer.password)
	if err != nil {
		log.Printf("Failed to execute flyway container: %v", err)
	}

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

func (dbm *DBMigrations) erase() {
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

	err = ptdb.pgContainer.container.Terminate(context.Background())
	if err != nil {
		log.Printf("Failed to terminate container: %v", err)
	}
}
