package domain

import "github.com/google/uuid"

type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;not null;index"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null;index"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
