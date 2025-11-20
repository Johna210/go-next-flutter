package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/johna210/go-next-flutter/internal/shared/model"
)

type UserProfile struct {
	model.BaseModel `gorm:"embedded"`

	UserID      uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	FirstName   string
	LastName    string
	PhoneNumber string
	AvatarURL   string
	Bio         string
	DateOfBirth *time.Time

	User *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
