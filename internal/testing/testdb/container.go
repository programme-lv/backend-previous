package testdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
)

type DBTestcontainer interface {
	GetConn() *sqlx.DB
	Close()
}

const (
	networkName   = "proglv-test-network"
	defaultUser   = "proglv"
	defaultPass   = "proglv"
	defaultDBName = "proglv"
)

func NewPostgresTestcontainer() (DBTestcontainer, error) {
	return initPostgresContainerTestDB()
}

type migratedPostgresTestcontainer struct {
	postgres *postgresContainer
	network  testcontainers.Network
	sqlxDb   *sqlx.DB
}

func initPostgresContainerTestDB() (x *migratedPostgresTestcontainer, err error) {
	x = &migratedPostgresTestcontainer{}

	x.network, err = createNewNetwork()
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	postgresAlias := randomLowercaseLetterString(10)
	x.postgres, err = startPostgresContainer(networkName, postgresAlias, defaultUser, defaultPass, defaultDBName)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	pgHost, pgPort, err := extractTestcontainerHostAndPort(x.postgres.container)
	if err != nil {
		return nil, fmt.Errorf("failed to extract testcontainer host and port: %w", err)
	}

	sqlxConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, x.postgres.user, x.postgres.password, x.postgres.database)
	log.Println("sqlxConnString: ", sqlxConnString)

	for i := 0; i < 10; i++ {
		db, err := sqlx.Connect("postgres", sqlxConnString)
		if err == nil {
			x.sqlxDb = db
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	migrations, err := cloneGitDBMigrations()
	if err != nil {
		return nil, fmt.Errorf("failed to clone git db migrations: %w", err)
	}
	defer migrations.erase()

	flywayAlias := randomLowercaseLetterString(10)
	err = execFlywayContainer(networkName, flywayAlias,
		migrations.getFlywayMigrationsDir(),
		postgresAlias, "5432", x.postgres.database, x.postgres.user, x.postgres.password)
	if err != nil {
		return nil, fmt.Errorf("failed to exec flyway container: %w", err)
	}

	return x, nil
}

func (ptdb *migratedPostgresTestcontainer) GetConn() *sqlx.DB {
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
