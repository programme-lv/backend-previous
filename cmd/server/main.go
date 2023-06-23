package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/programme-lv/backend/internal/environment"
	"github.com/programme-lv/backend/internal/graph"
)

const defaultPort = "3001"

func main() {
	conf := environment.ReadEnvConfig()
	conf.Print()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("http://localhost:%s/ = GraphQL playground", defaultPort)
	log.Printf("http://localhost:%s/query = GraphQL query endpoint", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}
