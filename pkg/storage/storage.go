package storage

import (
	"github.com/csams/doit/pkg/apis"
	"github.com/csams/doit/pkg/commands"
)

type Storage interface {
	Get(uint) (*apis.Task, error)
	Create(*commands.Create) error
	Update(*commands.Modify) error
	Delete(uint) error
	Search(*commands.Search) ([]*apis.Task, error)
}
