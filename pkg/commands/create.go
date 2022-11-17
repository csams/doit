package commands

import (
	"time"

	"github.com/csams/doit/pkg/apis/task"
	"github.com/csams/doit/pkg/set"
)

type Create struct {
	Description string        `json:"description"`
	Due         *time.Time    `json:"due"`
	Priority    task.Priority `json:"priority"`
	Status      task.Status   `json:"status"`
	Tags        []string      `json:"tags"`
}

func (c *Create) ToTask() *task.Task {
	t := task.Task{
		Description: c.Description,
		CreatedAt:   time.Now(),
		Due:         c.Due,
		Priority:    c.Priority,
		State:       task.Open,
		Status:      c.Status,
		Tags:        set.FromList(c.Tags),
	}
	return &t
}
