/*
 */
package sql

import (
	"github.com/csams/doit/pkg/apis"
	"github.com/csams/doit/pkg/commands"
	"github.com/csams/doit/pkg/storage"
	"github.com/csams/doit/pkg/storage/factory"
	"github.com/spf13/viper"
)

func init() {
	factory.Register("sql", New)
}

func New(v *viper.Viper) (storage.Storage, error) {
	return &sqlStorage{}, nil
}

func (db *sqlStorage) Get(_ task.Identity) (*task.Task, error) {
	panic("not implemented") // TODO: Implement
}

func (db *sqlStorage) Create(_ *commands.Create) error {
	panic("not implemented") // TODO: Implement
}

func (db *sqlStorage) Update(_ *commands.Modify) error {
	panic("not implemented") // TODO: Implement
}

func (db *sqlStorage) Delete(_ task.Identity) error {
	panic("not implemented") // TODO: Implement
}

func (db *sqlStorage) Search(_ *commands.Search) ([]*task.Task, error) {
	panic("not implemented") // TODO: Implement
}

type sqlStorage struct {
}
