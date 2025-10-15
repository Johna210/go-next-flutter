package handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/johna210/go-next-flutter/internal/delivery/http/dto"
	"github.com/johna210/go-next-flutter/internal/domain"
	"github.com/johna210/go-next-flutter/internal/usecase"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, input *dto.CreateUserRequest) (*dto.UserResponse, error) {
	user, err := h.userUseCase.CreateUser(ctx, input.Body.Email, input.Body.Name)
	if err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			return nil, huma.Error409Conflict("User with this email already exists")
		}

		if errors.Is(err, domain.ErrInvalidEmail) || errors.Is(err, domain.ErrInvalidName) {
			return nil, huma.Error400BadRequest(err.Error())
		}
		return nil, huma.Error500InternalServerError("Failed to create user")
	}

	return dto.ToUserResponse(user), nil
}

func (h *UserHandler) GetUser(ctx context.Context, input *dto.GetUserRequest) (*dto.UserResponse, error) {
	id, err := uuid.Parse(input.ID.String())
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid user ID format")
	}

	user, err := h.userUseCase.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, huma.Error404NotFound("User not found")
		}
		return nil, huma.Error500InternalServerError("Failed to get user")
	}

	return dto.ToUserResponse(user), nil
}

func (h *UserHandler) ListUsers(ctx context.Context, input *dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	users, err := h.userUseCase.ListUsers(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to list users")
	}

	return dto.ToListUsersResponse(users, input.Limit, input.Offset), nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, input *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	id, err := uuid.Parse(input.ID.String())
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid user ID format")
	}

	user, err := h.userUseCase.UpdateUser(ctx, id, input.Body.Name)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, huma.Error404NotFound("User not found")
		}
		if errors.Is(err, domain.ErrInvalidName) {
			return nil, huma.Error400BadRequest(err.Error())
		}
		return nil, huma.Error500InternalServerError("Failed to update user")
	}

	return dto.ToUserResponse(user), nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, input *dto.DeleteUserRequest) (*dto.MessageResponse, error) {
	id, err := uuid.Parse(input.ID.String())
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid user ID format")
	}

	err = h.userUseCase.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, huma.Error404NotFound("User not found")
		}
		return nil, huma.Error500InternalServerError("Failed to delete user")
	}

	resp := &dto.MessageResponse{}
	resp.Body.Message = "User deleted successfully"
	return resp, nil
}
