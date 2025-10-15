package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/johna210/go-next-flutter/internal/domain"
)

type IdUUIDPathParam struct {
	ID uuid.UUID `json:"id" path:"id" format:"uuid" doc:"User unique identifier" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type CreateUserRequest struct {
	Body struct {
		Email string `json:"email" format:"email" doc:"User email address" example:"user@example.com"`
		Name  string `json:"name" minLength:"1" maxLength:"100" doc:"User full name" example:"John Doe"`
	}
}

type UpdateUserRequest struct {
	IdUUIDPathParam
	Body struct {
		Name string `json:"name" minLength:"1" maxLength:"100" doc:"User full name" example:"John Doe"`
	}
}

type GetUserRequest struct {
	IdUUIDPathParam
}

type DeleteUserRequest struct {
	IdUUIDPathParam
}

type ListUsersRequest struct {
	Limit  int `query:"limit" minimum:"1" maximum:"100" default:"10" doc:"Number of users to return"`
	Offset int `query:"offset" minimum:"0" default:"0" doc:"Number of users to skip"`
}

type UserResponse struct {
	Body struct {
		ID        string    `json:"id" doc:"User ID" example:"123e4567-e89b-12d3-a456-426614174000"`
		Email     string    `json:"email" doc:"User email" example:"user@example.com"`
		Name      string    `json:"name" doc:"User name" example:"John Doe"`
		CreatedAt time.Time `json:"created_at" doc:"Creation timestamp"`
		UpdatedAt time.Time `json:"updated_at" doc:"Last update timestamp"`
	}
}

// ListUsersResponse represents a list of users response
type ListUsersResponse struct {
	Body struct {
		Users  []UserData `json:"users" doc:"List of users"`
		Total  int        `json:"total" doc:"Total number of users"`
		Limit  int        `json:"limit" doc:"Limit used"`
		Offset int        `json:"offset" doc:"Offset used"`
	}
}

type UserData struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Body struct {
		Error   string `json:"error" doc:"Error message"`
		Code    string `json:"code" doc:"Error code"`
		Details string `json:"details,omitempty" doc:"Additional error details"`
	}
}

type MessageResponse struct {
	Body struct {
		Message string `json:"message" doc:"Response message"`
	}
}

func ToUserResponse(user *domain.User) *UserResponse {
	resp := &UserResponse{}
	resp.Body.ID = user.ID.String()
	resp.Body.Email = user.Email
	resp.Body.Name = user.Name
	resp.Body.CreatedAt = user.CreatedAt
	resp.Body.UpdatedAt = user.UpdatedAt
	return resp
}

// ToListUsersResponse converts domain users to response DTO
func ToListUsersResponse(users []*domain.User, limit, offset int) *ListUsersResponse {
	resp := &ListUsersResponse{}
	resp.Body.Users = make([]UserData, len(users))
	for i, user := range users {
		resp.Body.Users[i] = UserData{
			ID:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}
	resp.Body.Total = len(users)
	resp.Body.Limit = limit
	resp.Body.Offset = offset
	return resp
}
