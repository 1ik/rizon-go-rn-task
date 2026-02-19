package main

import (
	"fmt"
	"log"
	"net/http"

	"rizon-test-task/internal/app"
	"rizon-test-task/internal/config"
	"rizon-test-task/internal/database"
	"rizon-test-task/internal/graphql"
	"rizon-test-task/internal/graphql/generated"
	"rizon-test-task/internal/in_memory_storage"
	"rizon-test-task/internal/repository/postgres"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
)

func main() {
	cfg := config.GetServerConfig()

	// Initialize database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize in-memory storage
	store, err := in_memory_storage.NewStore()
	if err != nil {
		log.Fatal("failed to initialize in-memory storage:", err)
	}
	defer in_memory_storage.Close()

	// Create repository
	userRepo := postgres.NewUserRepository(db)

	// Initialize app with dependencies
	application := app.New(userRepo, store)
	defer application.Close()

	resolver := graphql.NewResolver(application)

	// Create GraphQL executable schema
	executableSchema := generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	})

	// Create GraphQL handler
	graphqlHandler := handler.NewDefaultServer(executableSchema)

	// Setup CORS for Apollo Explorer
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	// GraphQL endpoint with CORS
	http.Handle("/graphql", c.Handler(graphqlHandler))

	// Apollo Sandbox (local)
	http.HandleFunc("/sandbox", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/graphql/apollo-sandbox.html")
	})

	// GraphQL Playground (optional, for development)
	http.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))

	fmt.Printf("Server starting on port %s\n", cfg.Port)
	fmt.Println("Endpoints:")
	fmt.Printf("  GET http://localhost:%s/health\n", cfg.Port)
	fmt.Printf("  POST http://localhost:%s/graphql\n", cfg.Port)
	fmt.Printf("  GET http://localhost:%s/sandbox (Apollo Sandbox - Local)\n", cfg.Port)
	fmt.Printf("  GET http://localhost:%s/ (GraphQL Playground)\n", cfg.Port)

	if err := http.ListenAndServe(cfg.Addr(), nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
