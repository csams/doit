package apis

import (
	"gorm.io/gorm"
)

type Annotation struct {
	gorm.Model

	TaskID      uint   `json:"taskid"`
	Description string `json:"description"`
}
