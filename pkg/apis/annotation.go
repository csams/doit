package apis

import (
	"time"

	"gorm.io/gorm"
)

type Annotation struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	TaskID      uint   `json:"taskid"`
	Description string `json:"description"`
}
