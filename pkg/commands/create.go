package commands

import (
	"time"

	"github.com/csams/doit/pkg/apis"
	"github.com/csams/doit/pkg/set"
)

type Create struct {
	Description string        `json:"description"`
	Due         *time.Time    `json:"due"`
	Priority    apis.Priority `json:"priority"`
	Status      apis.Status   `json:"status"`
	Tags        []string      `json:"tags"`
}

func (c *Create) ToTask() *apis.Task {
	t := apis.Task{
		Description: c.Description,
		Due:         c.Due,
		Priority:    c.Priority,
		State:       apis.Open,
		Status:      c.Status,
		Tags:        set.FromList(c.Tags),
	}
	return &t
}
