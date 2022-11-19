package apis

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username string `json:"username" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`

	OwnedTasks    []Task `gorm:"foreignKey:OwnerName;references:Username;constraint:OnDelete:CASCADE"`
	AssignedTasks []Task `gorm:"foreignKey:AssigneeName;references:Username"`

	SharedWith []Policy `gorm:"foreignKey:OwnerUsername;references:Username;constraint:OnDelete:CASCADE"`
	SharedFrom []Policy `gorm:"foreignKey:DelegateUsername;references:Username;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
