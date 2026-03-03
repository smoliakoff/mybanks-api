package main

import (
	"bytes"
	"context"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/formatter"
	"log"
	"mybanks-api/ent"
	"mybanks-api/graph"
	"mybanks-api/internal/config"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

// corsMiddleware sets CORS headers to allow requests from http://localhost:3000
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the desired origin or “*” if you want to allow all domains
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Process preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	_ = godotenv.Load(".env")
	cfg := config.Load()
	//dsn := "postgres://app:app@localhost:5432/postgres?sslmode=disable"

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed connecting to postgres: %v", err)
	}
	defer func(db *sql.Driver) {
		err := db.Close()
		if err != nil {
			fmt.Printf("error closing db: %v", err)
		}
	}(db)

	client := ent.NewClient(ent.Driver(db))
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Client: client}})

	srv := handler.New(schema)

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
	// schema through HTTP:
	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf)
	f.FormatSchema(schema.Schema())

	http.HandleFunc("/schema.graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write(buf.Bytes())
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// Wrap DefaultServeMux in corsMiddleware
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux)))
}
