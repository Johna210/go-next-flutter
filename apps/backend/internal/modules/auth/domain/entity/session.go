package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/johna210/go-next-flutter/internal/shared/model"
)

type Session struct {
	model.BaseModel

	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	JWTID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	RefreshToken string    `gorm:"not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	Revoked      bool      `gorm:"default:false"`
	IPAddress    string
	UserAgent    string
}

func (Session) TableName() string {
	return "sessions"
}
