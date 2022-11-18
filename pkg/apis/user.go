package apis

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username string `json:"username" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`

	Tasks []Task `gorm:"constraint:OnDelete:CASCADE"`

	SharedWith []Policy `gorm:"foreignKey:OwnerUsername;reference:Username;constraint:OnDelete:CASCADE"`
	SharedFrom []Policy `gorm:"foreignKey:DelegateUsername;reference:Username"`

	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
