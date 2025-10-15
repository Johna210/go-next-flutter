package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/johna210/go-next-flutter/internal/domain"
	"github.com/johna210/go-next-flutter/internal/repository"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
	// Check if user exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, domain.ErrUserExists
	}

	// Create new user
	user, err := domain.NewUser(email, name)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *UserUseCase) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return uc.userRepo.List(ctx, limit, offset)
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, id uuid.UUID, name string) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := user.Update(name); err != nil {
		return nil, err
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return uc.userRepo.Delete(ctx, id)
}
