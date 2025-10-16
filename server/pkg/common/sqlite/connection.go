package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func GetConnection(config Config) *sql.DB {
	// Ensure the directory exists
	dbDir := filepath.Dir(config.DatabasePath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Errorf("Failed to create database directory: %v", err)
		panic(err)
	}

	// Check if database file exists
	dbExists := fileExists(config.DatabasePath)

	// Open SQLite database
	db, err := sql.Open("sqlite", config.DatabasePath)
	if err != nil {
		log.Errorf("Unable to open database: %v", err)
		panic(err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		log.Errorf("Unable to ping database: %v", err)
		panic(err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		log.Errorf("Failed to enable WAL mode: %v", err)
		panic(err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		log.Errorf("Failed to enable foreign keys: %v", err)
		panic(err)
	}

	// Initialize schema if database is new
	if !dbExists {
		log.Info("Database file not found, initializing schema...")
		if err := initializeSchema(db); err != nil {
			log.Errorf("Failed to initialize schema: %v", err)
			panic(err)
		}
		log.Info("Database schema initialized successfully")
	}

	log.Info(fmt.Sprintf("Connected to SQLite database at: %s", config.DatabasePath))
	return db
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func initializeSchema(db *sql.DB) error {
	// Try to read the extended schema file
	schemaPath := "./schema_sqlite_extended.sql"
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		// Fall back to basic schema if extended file not found
		log.Warn("Extended schema file not found, using basic schema")
		return initializeBasicSchema(db)
	}

	schema := string(schemaBytes)
	_, err = db.Exec(schema)
	if err != nil {
		log.Errorf("Failed to execute extended schema: %v", err)
		return err
	}

	log.Info("Extended database schema initialized successfully")
	return nil
}

// initializeBasicSchema creates a minimal schema for backward compatibility
func initializeBasicSchema(db *sql.DB) error {
	schema := `
	-- Basic Users table
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Basic Products table
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		quantity INTEGER NOT NULL,
		category TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Create indexes
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
	CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
	`

	_, err := db.Exec(schema)
	return err
}
