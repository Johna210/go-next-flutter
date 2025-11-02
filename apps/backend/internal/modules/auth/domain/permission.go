package domain

import "github.com/johna210/go-next-flutter/internal/shared/model"

type Permission struct {
	model.BaseModel

	Name        string `gorm:"uniqueIndex;not null"`
	Description string

	Roles []RolePermission `gorm:"foreignKey:PermissionID"`
}

func (Permission) TableName() string {
	return "permissions"
}
