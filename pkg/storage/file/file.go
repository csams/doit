package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/csams/doit/pkg/apis/task"
	"github.com/csams/doit/pkg/commands"
	"github.com/csams/doit/pkg/storage"
	"github.com/csams/doit/pkg/storage/factory"
	"github.com/spf13/viper"
)

func init() {
	factory.Register("file", New)
}

func New(v *viper.Viper) (storage.Storage, error) {
	path := os.ExpandEnv(v.GetString("path"))
	if err := ensureExists(path); err != nil {
		return nil, err
	}
	return &fileStorage{
		Path: path,
	}, nil
}

type fileStorage struct {
	Path string
}

func (f *fileStorage) Get(id task.Identity) (*task.Task, error) {
	tasks, err := f.load()
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, fmt.Errorf("no task for id [%d]", id)
}

func (f *fileStorage) Create(c *commands.Create) error {
	tasks, err := f.load()
	if err != nil {
		return err
	}

	t := c.ToTask()
	t.ID = generateId(tasks)

	tasks = append(tasks, t)
	return f.save(tasks)
}

func (f *fileStorage) Update(c *commands.Modify) error {
	tasks, err := f.load()
	if err != nil {
		return err
	}

	for _, t := range tasks {
		if t.ID == c.Id {
			c.Apply(t)
			return f.save(tasks)
		}
	}
	return fmt.Errorf("no task to update for id [%d]", c.Id)
}

func (f *fileStorage) Delete(c task.Identity) error {
	tasks, err := f.load()
	if err != nil {
		return err
	}

	var res []*task.Task

	for _, t := range tasks {
		if t.ID != c {
			res = append(res, t)
		}
	}

	return f.save(res)
}

func (f *fileStorage) Search(c *commands.Search) ([]*task.Task, error) {
	return f.load()
}

func ensureExists(p string) error {
	if _, err := os.Stat(p); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		parent, _ := path.Split(p)
		if err = os.MkdirAll(parent, 0700); err != nil {
			return err
		}

		f, err := os.Create(p)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func (f *fileStorage) load() ([]*task.Task, error) {
	var tasks []*task.Task

	file, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		t := task.Task{}
		err := json.Unmarshal(scanner.Bytes(), &t)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil
}

func (f *fileStorage) save(tasks []*task.Task) error {
	file, err := os.OpenFile(f.Path, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	w := bufio.NewWriter(file)
	defer w.Flush()
	for _, t := range tasks {
		bytes, err := json.Marshal(&t)
		if err != nil {
			return err
		}
		_, err = w.Write(bytes)
		if err != nil {
			return err
		}
		w.WriteByte('\n')
	}

	return nil
}

func generateId(tasks []*task.Task) task.Identity {
	if len(tasks) == 0 {
		return task.Identity(1)
	}

	prev := tasks[0].ID
	for _, t := range tasks[1:] {
		if t.ID > prev+1 {
			return task.Identity(prev + 1)
		}
		prev = t.ID
	}
	return prev
}
