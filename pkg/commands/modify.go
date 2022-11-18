package commands

import (
	"time"

	"github.com/csams/doit/pkg/apis"
	"github.com/csams/doit/pkg/set"
)

type Modify struct {
	Id          uint          `json:"id"`
	Description *string       `json:"desc"`
	Due         *time.Time    `json:"due,omitempty"`
	Priority    apis.Priority `json:"priority,omitempty"`
	State       *apis.State   `json:"state"`
	Status      *apis.Status  `json:"status"`
	Tags        []string      `json:"tags,omitempty"`
}

func (update *Modify) Apply(orig *apis.Task) {
	if update.Description != nil {
		orig.Description = *update.Description
	}

	if update.Due != nil {
		orig.Due = update.Due
	}

	if update.Priority != 0 {
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
