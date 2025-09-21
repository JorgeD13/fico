package db

import (
	"fmt"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// OpenSQLiteAt opens (or creates) a SQLite database file at a path relative to the gol module.
func OpenSQLiteAt(relPath string) (*gorm.DB, error) {
	dbPath := filepath.FromSlash(relPath)
	dsn := fmt.Sprintf("file:%s?_fk=1", dbPath)
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}
