// api/database.go
package api

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&User{},
		&CalendarConnection{},
		&Link{},
		&Slot{},
		&Booking{},
		&Vote{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
