package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidName  = errors.New("invalid name")
	ErrInvalidID    = errors.New("invalid user ID")
)
