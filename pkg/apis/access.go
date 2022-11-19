package apis

import (
	"time"

	"gorm.io/gorm"
)

type Policy struct {
	OwnerUsername    string `json:"owner_user_name" gorm:"primaryKey"`
	DelegateUsername string `json:"delegate_user_name" gorm:"primaryKey"`

	Mode PolicyMode `json:"mode"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type PolicyMode string

const (
	View          = "view"
	ViewAndUpdate = "view_and_update"
)
