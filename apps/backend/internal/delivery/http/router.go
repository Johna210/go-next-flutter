package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"github.com/johna210/go-next-flutter/internal/delivery/http/handler"
)

func SetupRouter(userHandler *handler.UserHandler) (*gin.Engine, huma.API) {
	// Create Gin router
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "user-api",
		})
	})

	config := huma.DefaultConfig("Simple User API", "1.0.0")
	config.Info.Description = "A clean architecture REST API for user management with OpenAPI 3.1 documentation"
	config.Info.Contact = &huma.Contact{
		Name:  "test",
		Email: "test@test.com",
	}

	config.Servers = []*huma.Server{
		{URL: "http://localhost:8080", Description: "Local development server"},
	}

	// Initialize Huma with Gin adapter
	api := humagin.New(router, config)

	// Register user routes
	registerUserRoutes(api, userHandler)

	return router, api
}

func registerUserRoutes(api huma.API, handler *handler.UserHandler) {
	// Create user
	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/api/v1/users",
		Summary:     "Create a new user",
		Description: "Creates a new user account with the provided email and name",
		Tags:        []string{"Users"},
	}, handler.CreateUser)

	// Get user by ID
	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/api/v1/users/{id}",
		Summary:     "Get user by ID",
		Description: "Retrieves a user by their unique identifier",
		Tags:        []string{"Users"},
	}, handler.GetUser)

	// List users
	huma.Register(api, huma.Operation{
		OperationID: "list-users",
		Method:      http.MethodGet,
		Path:        "/api/v1/users",
		Summary:     "List users",
		Description: "Retrieves a paginated list of users",
		Tags:        []string{"Users"},
	}, handler.ListUsers)

	// Update user
	huma.Register(api, huma.Operation{
		OperationID: "update-user",
		Method:      http.MethodPut,
		Path:        "/api/v1/users/{id}",
		Summary:     "Update user",
		Description: "Updates user information",
		Tags:        []string{"Users"},
	}, handler.UpdateUser)

	// Delete user
	huma.Register(api, huma.Operation{
		OperationID: "delete-user",
		Method:      http.MethodDelete,
		Path:        "/api/v1/users/{id}",
		Summary:     "Delete user",
		Description: "Deletes a user by their ID",
		Tags:        []string{"Users"},
	}, handler.DeleteUser)
}
