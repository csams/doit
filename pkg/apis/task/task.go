package task

import (
	"sort"
	"time"

	"github.com/csams/doit/pkg/set"
)

type Identity uint64

// Task is some unit of work to do
type Task struct {
	ID          Identity         `json:"id" gorm:"primaryKey"`
	Description string           `json:"desc"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
	Due         *time.Time       `json:"due,omitempty"`
	Priority    Priority         `json:"priority,omitempty"`
	State       State            `json:"state"`
	Status      Status           `json:"status"`
	Tags        *set.Set[string] `json:"tags,omitempty"`
}

type State string
type Status string
type Priority int32

const (
	Undefined Priority = -1
	Lowest    Priority = 0
	Low       Priority = 1
	Medium    Priority = 2
	High      Priority = 3
)

const (
	Closed State = "closed"
	Open   State = "open"
)

const (
	Todo      Status = "todo"
	Started   Status = "started"
	Stopped   Status = "stopped"
	Done      Status = "done"
	Abandoned Status = "abandoned"
)

var validStatuses = set.New(Todo, Started, Done, Abandoned)

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

func (p Priority) String() string {
	switch p {
	case -1:
		return "Undefined"
	case 0:
		return "Lowest"
	case 1:
		return "Low"
	case 2:
		return "Medium"
	case 3:
		return "High"
	}
	return "Lowest"
}
