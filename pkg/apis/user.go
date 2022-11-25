package apis

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Name     string `json:"name" gorm:"not null"`

	OwnedTasks    []Task `gorm:"foreignKey:OwnerId;constraint:OnDelete:CASCADE"`
	AssignedTasks []Task `gorm:"foreignKey:AssigneeId"`

	SharedWith []Policy `gorm:"foreignKey:OwnerUserId;constraint:OnDelete:CASCADE"`
	SharedFrom []Policy `gorm:"foreignKey:DelegateUserId;constraint:OnDelete:CASCADE"`
}
