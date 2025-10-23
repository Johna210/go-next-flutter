package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GenericRepository defines common CRUD operations
type GenericRepository[T any] interface {
	// Create creates a new entity
	Create(ctx context.Context, entity *T) error

	// BulkCreate creates multiple entities
	BulkCreate(ctx context.Context, entities []*T) error

	// GetByID retrieves an entity by ID
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)

	// Update updates an entity
	Update(ctx context.Context, entity *T) error

	// Delete soft deletes an entity
	Delete(ctx context.Context, id uuid.UUID) error

	// HardDelete permanently deletes an entity
	HardDelete(ctx context.Context, id uuid.UUID) error

	// FindAll retrieves all entities with pagination
	FindAll(ctx context.Context, page, pageSize int) ([]*T, int64, error)

	// FindByIDs retrieves multiple entities by IDs
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*T, error)

	// Count counts total entities
	Count(ctx context.Context) (int64, error)

	// Exists checks if an entity exists
	Exists(ctx context.Context, id uuid.UUID) (bool, error)

	// Transaction executes a function within a transaction
	Transaction(ctx context.Context, fn func(repo GenericRepository[T]) error) error

	// GetDB returns the underlying GORM DB instance for custom queries
	GetDB() *gorm.DB
}

// QueryOptions defines options for queries
type QueryOptions struct {
	Preload        []string               // Relations to preload
	Order          string                 // Order by clause
	Conditions     map[string]interface{} // Where conditions
	TenantID       string                 // Tenant isolation
	IncludeDeleted bool                   // Include soft deleted records
}

type PaginatedResult[T any] struct {
	Data       []*T  `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}
