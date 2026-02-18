package database

import (
	"fmt"
	"log"

	"rizon-test-task/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DB is the global database connection
	DB *gorm.DB
)

// Connect initializes the database connection
func Connect() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}

	cfg := config.GetDatabaseConfig()
	dsn := cfg.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to database: %s@%s:%s/%s",
		cfg.User, cfg.Host, cfg.Port, cfg.DBName)

	DB = db
	return DB, nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// GetDB returns the global database connection
func GetDB() *gorm.DB {
	return DB
}
