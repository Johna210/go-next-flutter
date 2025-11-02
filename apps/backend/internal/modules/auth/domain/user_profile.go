package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/johna210/go-next-flutter/internal/shared/model"
)

type UserProfile struct {
	model.BaseModel

	UserID      uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	FirstName   string
	LastName    string
	PhoneNumber string
	AvatarURL   string
	Bio         string
	DateOfBirth *time.Time

	User *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
