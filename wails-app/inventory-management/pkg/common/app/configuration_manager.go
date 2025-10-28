package app

import (
	"inventory-management/pkg/common/sqlite"
	"os"
)

type ConfigurationManager struct {
	SqliteConfig sqlite.Config
}

func NewConfigurationManager() *ConfigurationManager {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/inventory.db"
	}

	sqliteConfig := sqlite.Config{
		DatabasePath: dbPath,
	}
	return &ConfigurationManager{sqliteConfig}
}
