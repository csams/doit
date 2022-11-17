package memory

import (
	"fmt"
	"sort"

	"github.com/csams/doit/pkg/apis/task"
	"github.com/csams/doit/pkg/commands"
	"github.com/csams/doit/pkg/storage"
	"github.com/csams/doit/pkg/storage/factory"
	"github.com/spf13/viper"
)

func init() {
	factory.Register("memory", New)
}

func New(v *viper.Viper) (storage.Storage, error) {
	return &memoryStorage{
		tasks: make(map[task.Identity]*task.Task),
	}, nil
}

type memoryStorage struct {
	tasks map[task.Identity]*task.Task
}

func (m *memoryStorage) Get(id task.Identity) (*task.Task, error) {
	t, exists := m.tasks[id]
	if exists {
		return t, nil
	}
	return nil, fmt.Errorf("no task for id [%d]", id)
}

func (m *memoryStorage) Create(c *commands.Create) error {
	id := task.Identity(1)

	if len(m.tasks) > 0 {
		for k := range m.tasks {
			if k > id {
				id = k
			}
		}
		id = id + 1
	}

	t := c.ToTask()
	t.ID = id

	m.tasks[id] = t
	return nil
}

func (m *memoryStorage) Update(c *commands.Modify) error {
	t, exists := m.tasks[c.Id]
	if !exists {
		return fmt.Errorf("task [%d] can't be found for update", c.Id)
	}
	c.Apply(t)
	return nil
}

func (m *memoryStorage) Delete(id task.Identity) error {
	delete(m.tasks, id)
	return nil
}

func (m *memoryStorage) Search(c *commands.Search) ([]*task.Task, error) {
	res := make([]*task.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		res = append(res, t)
	}
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})
	return res, nil
}
