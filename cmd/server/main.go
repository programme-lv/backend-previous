package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/programme-lv/backend/internal/environment"
	"github.com/programme-lv/backend/internal/graphql"
	amqp "github.com/rabbitmq/amqp091-go"
)

const defaultPort = "3001"

func main() {
	conf := environment.ReadEnvConfig()
	conf.Print()

	log.Println("Connecting to database...")
	sqlxDb := sqlx.MustConnect("postgres", conf.SqlxConnString)
	defer sqlxDb.Close()
	log.Println("Connected to database")

	log.Println("Connecting to RabbitMQ...")
	rmqConn, err := amqp.Dial(conf.AMQPConnString)
	if err != nil {
		panic(err)
	}
	defer rmqConn.Close()
	log.Println("Connected to RabbitMQ")

	sessions := scs.New()
	sessions.Lifetime = 24 * time.Hour

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	resolver := &graphql.Resolver{
		PostgresDB:     sqlxDb,
		SessionManager: sessions,
		Logger:         logger,
		SubmissionRMQ:  rmqConn,
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world :)"))
	})
	http.Handle("/query", sessions.LoadAndSave(srv))

	log.Println("Listening on port " + defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}
