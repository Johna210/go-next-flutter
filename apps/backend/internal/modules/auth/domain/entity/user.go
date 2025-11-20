package entity

import "github.com/johna210/go-next-flutter/internal/shared/model"

type User struct {
	model.BaseModel `gorm:"embedded"`

	Username     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	IsActive     bool   `gorm:"default:false"`

	Profile UserProfile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Roles   []UserRole  `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
