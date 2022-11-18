package apis

import (
	"gorm.io/gorm"
)

type Policy struct {
	gorm.Model
	OwnerUsername    string `json:"owner_user_name"`
	DelegateUsername string `json:"delegate_user_name"`

	Mode PolicyMode `json:"mode"`
}

type PolicyMode string

const (
	View          = "view"
	ViewAndUpdate = "view_and_update"
)
