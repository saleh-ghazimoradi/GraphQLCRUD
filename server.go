package main

import (
	"context"
	"github.com/saleh-ghazimoradi/GraphQLCRUD/database"
	"github.com/saleh-ghazimoradi/GraphQLCRUD/repository"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/saleh-ghazimoradi/GraphQLCRUD/graph"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	mongoDB := database.NewMongoDB(
		database.WithHost("localhost"),
		database.WithPort(27018),
		database.WithUser("graphql"),
		database.WithPass("graphql"),
		database.WithDBName("job_listings"),
		database.WithAuthSource("admin"),
		database.WithMaxPoolSize(10),
		database.WithMinPoolSize(2),
		database.WithTimeout(10*time.Second),
	)

	client, db, err := mongoDB.Connect()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer func() {
		if err := mongoDB.Disconnect(context.Background(), client); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	jobRepo := repository.NewJobRepository(db, "jobs")
	graph.SetJobRepository(jobRepo)
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	port := defaultPort

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
