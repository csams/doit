package tui

import (
	"time"

	"github.com/araddon/dateparse"
	"github.com/csams/doit/pkg/apis"
)

type taskFormData struct {
	Description string
	Due         string
	State       apis.State
	Status      apis.Status
	Priority    apis.Priority
	Private     bool
}

func (d *taskFormData) ApplyTo(t *apis.Task) error {
	var due *time.Time
	if d.Due != "" {
		if dueDate, err := dateparse.ParseLocal(d.Due); err != nil {
			return err
		} else {
			due = &dueDate
		}
	}
	t.Description = d.Description
	t.Due = due
	t.State = d.State
	t.Status = d.Status
	t.Priority = d.Priority
	t.Private = d.Private

	return nil
}

func formDataFromTask(task *apis.Task) *taskFormData {
	fmtString := "2006-01-02 15:04:05 MST"
	var due string
	if task.Due != nil {
		due = task.Due.Format(fmtString)
	}
	return &taskFormData{
		Description: task.Description,
		Due:         due,
		State:       task.State,
		Status:      task.Status,
		Priority:    task.Priority,
		Private:     task.Private,
	}
}
