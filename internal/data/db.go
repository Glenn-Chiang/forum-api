package data

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func MustOpenDB(dsn string) *gorm.DB {
	db, err := OpenDB(dsn)
	if err != nil {
		panic("failed to connect to database")
	}
	return db
}
