package commands

import (
	"time"

	"github.com/csams/doit/pkg/apis/task"
	"github.com/csams/doit/pkg/set"
)

type Modify struct {
	Id          task.Identity `json:"id"`
	Description *string       `json:"desc"`
	Due         *time.Time    `json:"due,omitempty"`
	Priority    task.Priority `json:"priority,omitempty"`
	State       *task.State   `json:"state"`
	Status      *task.Status  `json:"status"`
	Tags        []string      `json:"tags,omitempty"`
}

func (update *Modify) Apply(orig *task.Task) {
	if update.Description != nil {
		orig.Description = *update.Description
	}

	if update.Due != nil {
		orig.Due = update.Due
	}

	if update.Priority != task.Lowest {
		orig.Priority = update.Priority
	}

	if update.State != nil {
		orig.State = *update.State
	}

	if update.Status != nil {
		orig.Status = *update.Status
	}

	if update.Tags != nil {
		orig.Tags = set.FromList(update.Tags)
	}
}
