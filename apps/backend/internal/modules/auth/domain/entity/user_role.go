package entity

import "github.com/google/uuid"

type UserRole struct {
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	RoleID uuid.UUID `gorm:"type:uuid;not null;index"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
