package storage

import (
	"gorm.io/gorm"

	"github.com/csams/doit/pkg/apis"
)

// Migrate the schemas
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&apis.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&apis.Task{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&apis.Policy{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&apis.Comment{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&apis.Annotation{}); err != nil {
		return err
	}
	return nil
}
