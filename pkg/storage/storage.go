package storage

import (
	"github.com/csams/doit/pkg/apis/task"
	"github.com/csams/doit/pkg/commands"
)

type Storage interface {
	Get(task.Identity) (*task.Task, error)
	Create(*commands.Create) error
	Update(*commands.Modify) error
	Delete(task.Identity) error
	Search(*commands.Search) ([]*task.Task, error)
}
