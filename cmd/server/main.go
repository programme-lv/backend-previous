package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/programme-lv/backend/internal/environment"
	"github.com/programme-lv/backend/internal/graphql"
)

const defaultPort = "3001"

func main() {
	conf := environment.ReadEnvConfig()
	conf.Print()

	sqlxDb := sqlx.MustConnect("postgres", conf.SqlxConnString)
	defer sqlxDb.Close()
	log.Println("Connected to database")

	// Initialize session manager
	sessions := scs.New()
	sessions.Lifetime = 24 * time.Hour

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	resolver := &graphql.Resolver{
		DB:             sqlxDb,
		SessionManager: sessions,
		Logger:         logger,
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", sessions.LoadAndSave(srv))

	log.Printf("http://localhost:%s/ = GraphQL playground", defaultPort)
	log.Printf("http://localhost:%s/query = GraphQL query endpoint", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}
