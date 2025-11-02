package repository

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/johna210/go-next-flutter/internal/core"
	collectionquery "github.com/johna210/go-next-flutter/pkg/collection_query"
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

func (r *BaseRepository[T]) FindAll(ctx context.Context, query collectionquery.CollectionQuery) PaginatedResult[T] {
	const defaultPageSize = 10
	const defaultSkip = 0

	pageSize := defaultPageSize
	if query.Take != nil && *query.Take > 0 {
		pageSize = *query.Take
	}

	skip := defaultSkip
	if query.Skip != nil && *query.Skip >= 0 {
		skip = *query.Skip
	}

	qc := collectionquery.QueryConstructor[T]{}

	result, err := qc.Find(r.db.WithContext(ctx), query, false)
	if err != nil {
		r.logger.Error("Failed to find all entities", core.Error(err))
		return PaginatedResult[T]{}
	}

	if result.Total == 0 {
		return PaginatedResult[T]{
			Total:      0,
			Data:       result.Items,
			Page:       1,
			PageSize:   pageSize,
			TotalPages: 0,
		}
	}

	var totalPages int
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(result.Total) / float64(pageSize)))
	}

	currentPage := (skip / pageSize) + 1

	return PaginatedResult[T]{
		Total:      result.Total,
		Data:       result.Items,
		Page:       currentPage,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

func (r *BaseRepository[T]) FindAllArchived(
	ctx context.Context,
	query collectionquery.CollectionQuery,
) PaginatedResult[T] {
	query.Where = append(query.Where, []collectionquery.Where{
		{
			Column:   "deletedAt",
			Value:    "",
			Operator: collectionquery.IsNotNull,
		},
	})

	const defaultPageSize = 10
	const defaultSkip = 0

	pageSize := defaultPageSize
	if query.Take != nil && *query.Take > 0 {
		pageSize = *query.Take
	}

	skip := defaultSkip
	if query.Skip != nil && *query.Skip >= 0 {
		skip = *query.Skip
	}

	qc := collectionquery.QueryConstructor[T]{}
	result, err := qc.Find(r.db.WithContext(ctx), query, true)
	if err != nil {
		r.logger.Error("Failed to find all entities", core.Error(err))
		return PaginatedResult[T]{}
	}

	if result.Total == 0 {
		return PaginatedResult[T]{
			Total:      0,
			Data:       result.Items,
			Page:       1,
			PageSize:   pageSize,
			TotalPages: 0,
		}
	}

	var totalPages int
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(result.Total) / float64(pageSize)))
	}

	currentPage := (skip / pageSize) + 1

	return PaginatedResult[T]{
		Total:      result.Total,
		Data:       result.Items,
		Page:       currentPage,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
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
