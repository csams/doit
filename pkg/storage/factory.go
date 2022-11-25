package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(c CompletedConfig) (*gorm.DB, error) {
	if c.DSN != "" {
		return gorm.Open(postgres.Open(c.DSN), &gorm.Config{})
	}
	return gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
}
