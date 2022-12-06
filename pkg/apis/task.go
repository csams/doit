package apis

import (
	"net/http"
	"sort"
	"time"

	"gorm.io/gorm"

	"github.com/csams/doit/pkg/set"
)

// Task is some unit of work to do
type Task struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OwnerId uint  `json:"owner_id"`
	Owner   *User `gorm:"foreignKey:ID;references:OwnerId" json:"-"`

	AssigneeId uint  `json:"assignee_id"`
	Assignee   *User `gorm:"foreignKey:ID;references:AssigneeId" json:"-"`

	Description string       `json:"desc"`
	Due         *time.Time   `json:"due"`
	Priority    Priority     `json:"priority"`
	Private     bool         `json:"private"`
	State       State        `json:"state"`
	Status      Status       `json:"status"`
	Comments    []Comment    `json:"comments" gorm:"constraint:OnDelete:CASCADE"`
	Annotations []Annotation `json:"annotations" gorm:"constraint:OnDelete:CASCADE"`
	// Tags        *set.Set[string] `json:"tags,omitempty"`
}

func (t *Task) Bind(r *http.Request) error {
	return nil
}

// Priority is how urgent the task is. 0 is lowest priority.
type Priority uint8

// State is either open or closed
type State string

// Status is one of backlog, todo, doing, done, abandoned
type Status string

const (
	Closed State = "closed"
	Open   State = "open"

	Backlog   Status = "backlog"
	Todo      Status = "todo"
	Doing     Status = "doing"
	Done      Status = "done"
	Abandoned Status = "abandoned"
)

var (
	validStatuses = set.New(Backlog, Todo, Doing, Done, Abandoned)
)

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
