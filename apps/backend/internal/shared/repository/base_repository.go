package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/johna210/go-next-flutter/internal/core"
)

var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
)

type BaseRepository[T any] struct {
	db     *core.Database
	logger core.Logger
}

func NewBaseRepository[T any](db *core.Database, logger core.Logger) GenericRepository[T] {
	return &BaseRepository[T]{
		db:     db,
		logger: logger,
	}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create entity", core.Error(err))
		return err
	}
	return nil
}

func (r *BaseRepository[T]) BulkCreate(ctx context.Context, entities []*T) error {
	if err := r.db.WithContext(ctx).CreateInBatches(entities, 100).Error; err != nil {
		r.logger.Error("Failed to bulk create entities", core.Error(err))
		return err
	}
	return nil
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		r.logger.Error("Failed to get entity by ID", core.Error(err))
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update entity", core.Error(err))
		return err
	}
	return nil
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(new(T), "id = ?", id).Error; err != nil {
		r.logger.Error("Failed to delete entity", core.Error(err))
		return err
	}
	return nil
}

func (r *BaseRepository[T]) HardDelete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Unscoped().Delete(new(T), "id = ?", id).Error; err != nil {
		r.logger.Error("Failed to hard delete entity", core.Error(err))
		return err
	}
	return nil
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, page, pageSize int) ([]*T, int64, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	var entities []*T
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		r.logger.Error("Failed to count entities", core.Error(err))
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		r.logger.Error("Failed to find entities", core.Error(err))
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *BaseRepository[T]) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*T, error) {
	var entities []*T
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&entities).Error; err != nil {
		r.logger.Error("Failed to find entities by IDs", core.Error(err))
		return nil, err
	}
	return entities, nil
}

func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(new(T)).Count(&count).Error; err != nil {
		r.logger.Error("Failed to count entities", core.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *BaseRepository[T]) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Count(&count).Error; err != nil {
		r.logger.Error("Failed to check entity existence", core.Error(err))
		return false, err
	}

	return count > 0, nil
}

func (r *BaseRepository[T]) Transaction(ctx context.Context, fn func(repo GenericRepository[T]) error) error {
	return r.db.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := &BaseRepository[T]{
			db:     &core.Database{DB: tx},
			logger: r.logger,
		}
		return fn(txRepo)
	})
}

func (r *BaseRepository[T]) GetDB() *gorm.DB {
	return r.db.DB
}
