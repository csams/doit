package apis

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model

	TaskID      uint   `json:"taskid"`
	Description string `json:"description"`
}
