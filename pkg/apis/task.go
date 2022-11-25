package apis

import (
	"sort"
	"time"

	"gorm.io/gorm"

	"github.com/csams/doit/pkg/set"
)

// Task is some unit of work to do
type Task struct {
	gorm.Model

	OwnerId uint `json:"owner"`
	Owner   User `gorm:"foreignKey:ID;references:OwnerId"`

	AssigneeId uint `json:"assignee"`
	Assignee   User `gorm:"foreignKey:ID;references:AssigneeId"`

	Description string       `json:"desc"`
	Due         *time.Time   `json:"due,omitempty"`
	Priority    Priority     `json:"priority,omitempty"`
	Private     bool         `json:"private,omitempty"`
	State       State        `json:"state"`
	Status      Status       `json:"status"`
	Comments    []Comment    `json:"comments,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Annotations []Annotation `json:"annotations,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	// Tags        *set.Set[string] `json:"tags,omitempty"`
}

// Priority is how urgent the task is. 0 is lowest priority.
type Priority uint8

type State string

const (
	Closed State = "closed"
	Open   State = "open"
)

type Status string

const (
	Backlog   Status = "backlog"
	Todo      Status = "todo"
	Doing     Status = "doing"
	Done      Status = "done"
	Abandoned Status = "abandoned"
)

var validStatuses = set.New(Todo, Todo, Doing, Done, Abandoned)

func IsValidStatus(s Status) bool {
	return validStatuses.Has(s)
}

func Statuses() []Status {
	l := validStatuses.ToList()
	sort.SliceStable(l, func(i, j int) bool {
		return l[i] < l[j]
	})
	return l
}

func StatusStrings() []string {
	valid := Statuses()
	statuses := make([]string, 0, len(valid))
	for _, s := range valid {
		statuses = append(statuses, string(s))
	}
	return statuses
}
