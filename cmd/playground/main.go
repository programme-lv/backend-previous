package main

import (
	"fmt"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/config"

	_ "github.com/lib/pq"
	"log/slog"
)

var logger *slog.Logger

func main() {
	logger = slog.Default()

	conf := loadConfigFromEnvFile()

	pgDB := connectToPostgres(conf.Postgres.Host, conf.Postgres.User, conf.Postgres.Password,
		conf.Postgres.DBName, conf.Postgres.SSLMode, conf.Postgres.Port)

	var id struct {
		ID int64
	}
	err := postgres.SELECT(postgres.Raw("nextval('task_submissions_id_seq'::regclass)").AS("id")).Query(pgDB, &id)
	if err != nil {
		panic(err)
	}

	fmt.Println(id.ID)
}

func loadConfigFromEnvFile() *config.Config {
	conf, err := config.LoadConfig(".env.toml")
	if err != nil {
		logger.Error("could not load config", "error", err)
		panic(err)
	}
	conf.Print()
	return conf
}

func connectToPostgres(host, user, password, dbname, sslmode string, port int) *sqlx.DB {
	connStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	logger.Info("connecting to PostgreSQL")
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		logger.Error("could not connect to PostgreSQL", "error", err)
		panic(err)
	}
	logger.Info("successfully connected to PostgreSQL")

	return db
}
