package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jmoiron/sqlx"

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

	resolver := &graphql.Resolver{
		DB: sqlxDb,
	}

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("http://localhost:%s/ = GraphQL playground", defaultPort)
	log.Printf("http://localhost:%s/query = GraphQL query endpoint", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}