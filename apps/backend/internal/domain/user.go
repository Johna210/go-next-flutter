package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, name string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}

	if name == "" {
		return nil, ErrInvalidName
	}

	return &User{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (u *User) Update(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}
