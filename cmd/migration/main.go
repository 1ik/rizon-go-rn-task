package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"rizon-test-task/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	var (
		command = flag.String("command", "up", "Migration command: up or down")
		dir     = flag.String("dir", "migrations", "Directory with migration files")
	)
	flag.Parse()

	// Get database URL
	cfg := config.GetDatabaseConfig()
	dbURL := cfg.URL()

	// Get migrations directory (relative to project root)
	migrationsDir := *dir
	if !filepath.IsAbs(migrationsDir) {
		// Assume migrations directory is relative to project root
		migrationsDir = filepath.Join(".", migrationsDir)
	}

	// Execute command
	switch *command {
	case "up":
		if err := runUp(dbURL, migrationsDir); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}

	case "down":
		if err := runDown(dbURL, migrationsDir); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}

	default:
		log.Fatalf("Unknown command: %s. Use 'up' or 'down'", *command)
	}
}

func getProvider(dbURL, migrationsDir string) (*goose.Provider, *sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create filesystem from migrations directory
	fsys := os.DirFS(migrationsDir)

	// Create provider
	provider, err := goose.NewProvider(goose.DialectPostgres, db, fsys)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to create provider: %w", err)
	}

	return provider, db, nil
}

func runUp(dbURL, migrationsDir string) error {
	provider, db, err := getProvider(dbURL, migrationsDir)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx := context.Background()
	results, err := provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if len(results) > 0 {
		fmt.Printf("✅ Applied %d migration(s) successfully\n", len(results))
		for _, result := range results {
			path := ""
			if result.Source != nil {
				path = result.Source.Path
			}
			fmt.Printf("  - %s\n", path)
		}
	} else {
		fmt.Println("✅ No migrations to apply")
	}
	return nil
}

func runDown(dbURL, migrationsDir string) error {
	provider, db, err := getProvider(dbURL, migrationsDir)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx := context.Background()
	result, err := provider.Down(ctx)
	if err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	if result != nil {
		path := ""
		if result.Source != nil {
			path = result.Source.Path
		}
		fmt.Printf("✅ Rolled back migration: %s\n", path)
	} else {
		fmt.Println("✅ No migrations to rollback")
	}
	return nil
}
