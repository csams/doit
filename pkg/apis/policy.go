package apis

import (
	"time"

	"gorm.io/gorm"
)

type Policy struct {
	OwnerUserId    uint `json:"owner_user_id" gorm:"primaryKey"`
	DelegateUserId uint `json:"delegate_user_id" gorm:"primaryKey"`

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
