package testdb

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
)

// TestDBProvider provides a migrated database
// and a method to close it (don't forget to do it).
type TestDBProvider interface {
	GetTestDB() *sqlx.DB
	Close()
}

const (
	networkName   = "proglv-test-network"
	defaultUser   = "proglv"
	defaultPass   = "proglv"
	defaultDBName = "proglv"
)

func NewPostgresTestcontainerProvider() (TestDBProvider, error) {
	return initPostgresContainerTestDB()
}

type migratedPostgresTestcontainer struct {
	postgres *postgresContainer
	network  testcontainers.Network
	sqlxDb   *sqlx.DB
}

func initPostgresContainerTestDB() (x *migratedPostgresTestcontainer, err error) {
	x = &migratedPostgresTestcontainer{}

	x.network, err = createNetwork(networkName)
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	postgresAlias := randomLowercaseLetterString(10)
	x.postgres, err = startPostgresContainer(networkName, postgresAlias, defaultUser, defaultPass, defaultDBName)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	host, port, err := extractTestcontainerHostAndPort(x.postgres.container)
	if err != nil {
		return nil, fmt.Errorf("failed to extract testcontainer host and port: %w", err)
	}

	sqlxConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, x.postgres.user, x.postgres.password, x.postgres.database)
	log.Println("sqlxConnString: ", sqlxConnString)

	x.sqlxDb = sqlx.MustConnect("postgres", sqlxConnString)

	migrations, err := cloneGitDBMigrations()
	if err != nil {
		return nil, fmt.Errorf("failed to clone git db migrations: %w", err)
	}
	defer migrations.erase()

	err = execFlywayContainer(migrations.getFlywayMigrationsDir(),
		host, port, x.postgres.database, x.postgres.user, x.postgres.password)
	if err != nil {
		return nil, fmt.Errorf("failed to exec flyway container: %w", err)
	}

	return x, nil
}

func (ptdb *migratedPostgresTestcontainer) GetTestDB() *sqlx.DB {
	return ptdb.sqlxDb
}

func (x *migratedPostgresTestcontainer) Close() {
	err := x.sqlxDb.Close()
	if err != nil {
		log.Printf("failed to close sqlx db: %v", err)
	}

	err = x.postgres.container.Terminate(context.Background())
	if err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	err = x.network.Remove(context.Background())
	if err != nil {
		log.Printf("failed to remove network: %v", err)
	}
}
