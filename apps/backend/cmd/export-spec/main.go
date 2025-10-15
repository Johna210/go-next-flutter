package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/johna210/go-next-flutter/internal/delivery/http"
	"github.com/johna210/go-next-flutter/internal/delivery/http/handler"
	"github.com/johna210/go-next-flutter/internal/infrastructure/memory"
	"github.com/johna210/go-next-flutter/internal/usecase"
)

func main() {
	// Initialize dependencies (same as main)
	userRepo := memory.NewUserRepository()
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	// Setup router to get API spec
	_, api := http.SetupRouter(userHandler)

	// Get OpenAPI spec
	spec := api.OpenAPI()

	// Write to file
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal OpenAPI spec: %v", err)
	}

	outputPath := "../../api-spec/openapi.json"
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		log.Fatalf("Failed to write OpenAPI spec: %v", err)
	}

	log.Printf("✅ OpenAPI spec exported to: %s", outputPath)
}
