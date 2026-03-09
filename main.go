package main

import (
	"bytes"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/formatter"
	"log"
	"mybanks-api/ent"
	"strings"

	"mybanks-api/graph"
	"mybanks-api/internal/config"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

var BuildSHA = "dev"

func corsMiddleware(allowedOrigins map[string]struct{}, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, ok := allowedOrigins[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Internal-API-Key")
		}

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Printf("BuildSHA=%s", BuildSHA)

	_ = godotenv.Load(".env")
	cfg := config.Load()

	allowedOrigins := parseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"))
	internalAPIKey := os.Getenv("INTERNAL_API_KEY")

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

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Client: client}})

	srv := handler.New(schema)

	srv.AddTransport(transport.Options{})
	if cfg.Environment == "development" {
		srv.AddTransport(transport.GET{})
	}
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	if cfg.Environment == "development" {
		srv.Use(extension.Introspection{})
	}
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/query", internalAuthMiddleware(internalAPIKey, srv))
	if cfg.Environment == "development" {
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
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
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// Wrap DefaultServeMux in corsMiddleware
	rootHandler := corsMiddleware(allowedOrigins, http.DefaultServeMux)
	log.Fatal(http.ListenAndServe(":"+port, rootHandler))
}

func internalAuthMiddleware(internalAPIKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isInternal := false

		if internalAPIKey != "" {
			provided := r.Header.Get("X-Internal-API-Key")
			if provided == internalAPIKey {
				isInternal = true
			}
		}

		ctx := graph.WithInternalRequest(r.Context(), isInternal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseAllowedOrigins(raw string) map[string]struct{} {
	result := make(map[string]struct{})

	for _, item := range strings.Split(raw, ",") {
		origin := strings.TrimSpace(item)
		if origin == "" {
			continue
		}
		result[origin] = struct{}{}
	}

	return result
}
