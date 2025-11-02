package domain

import "github.com/google/uuid"

type UserRole struct {
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	RoleID uuid.UUID `gorm:"type:uuid;not null;index"`
}
