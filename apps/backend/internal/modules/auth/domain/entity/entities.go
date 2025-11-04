package entity

import (
	"github.com/johna210/go-next-flutter/internal/core"
)

// EntityProvider implements core.EntityProvider for auth module
type EntityProvider struct{}

// NewEntityProvider creates the entity provider
func NewEntityProvider() core.EntityProvider {
	return &EntityProvider{}
}

// Entities returns all domain entities for auth module
func (p *EntityProvider) Entities() []interface{} {
	return []interface{}{
		&User{},
		&UserProfile{},
		&Session{},
		&Role{},
		&Permission{},
		&RolePermission{},
		&UserRole{},
	}
}

// ModuleName returns the module identifier
func (p *EntityProvider) ModuleName() string {
	return "auth"
}
