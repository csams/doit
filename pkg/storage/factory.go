package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(c CompletedConfig) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
}
