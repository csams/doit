package commands

import "github.com/csams/doit/pkg/apis/task"

type Remove struct {
	Id task.Identity `json:"id"`
}
