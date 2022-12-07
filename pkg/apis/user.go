package apis

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username string `gorm:"unique" json:"username"`
	Name     string `gorm:"not null" json:"name"`

	OwnedTasks    []Task `gorm:"foreignKey:OwnerId;constraint:OnDelete:CASCADE"`
	AssignedTasks []Task `gorm:"foreignKey:AssigneeId"`

	SharedWith []Policy `gorm:"foreignKey:OwnerUserId;constraint:OnDelete:CASCADE"`
	SharedFrom []Policy `gorm:"foreignKey:DelegateUserId;constraint:OnDelete:CASCADE"`
}
