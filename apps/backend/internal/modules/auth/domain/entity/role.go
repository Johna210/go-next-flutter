package entity

import "github.com/johna210/go-next-flutter/internal/shared/model"

type Role struct {
	model.BaseModel `gorm:"embedded"`

	Name        string `gorm:"uniqueIndex;not null"`
	Description string

	Permissions []RolePermission `gorm:"foreignKey:RoleID"`
	Users       []UserRole       `gorm:"foreignKey:RoleID"`
}

func (Role) TableName() string {
	return "roles"
}
