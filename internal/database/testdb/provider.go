package testdb

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"log"
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

type migratedPostgresTestcontainer struct {
	container *postgresContainer
	network   testcontainers.Network
	sqlxDb    *sqlx.DB
}

func initPostgresContainerTestDB() (x *migratedPostgresTestcontainer, err error) {
	x = &migratedPostgresTestcontainer{}

	x.network, err = createNetwork("proglv-test-network")
	if err != nil {
		return nil, err
	}

	pgUsername := "proglv"
	pgPassword := "proglv"
	pgDatabase := "proglv"

	x.container, err = startPostgresContainer(pgUsername, pgPassword, pgDatabase)
	if err != nil {
		return nil, err
	}

	host, port, err := extractTestcontainerHostPort(x.container.container)
	if err != nil {
		return nil, err
	}

	sqlxConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, x.container.user, x.container.password, x.container.database)
	log.Printf("connString: %s", sqlxConnString)

	x.sqlxDb = sqlx.MustConnect("postgres", sqlxConnString)

	migrations, err := cloneDBMigrations()
	if err != nil {
		log.Printf("Failed to clone migrations: %v", err)
	}
	defer migrations.erase()

	err = execFlywayContainer(migrations.getFlywayMigrationsDir(),
		host, port, x.container.database, x.container.user, x.container.password)
	if err != nil {
		log.Printf("Failed to execute flyway container: %v", err)
	}

	return x, nil
}

func (ptdb *migratedPostgresTestcontainer) GetTestDB() *sqlx.DB {
	return ptdb.sqlxDb
}

func (x *migratedPostgresTestcontainer) Close() {
	err := x.sqlxDb.Close()
	if err != nil {
		log.Printf("Failed to close sqlx db: %v", err)
	}

	err = x.network.Remove(context.Background())
	if err != nil {
		log.Printf("Failed to remove network: %v", err)
	}

	err = x.container.container.Terminate(context.Background())
	if err != nil {
		log.Printf("Failed to terminate container: %v", err)
	}
}
