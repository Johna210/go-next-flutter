package main

import (
	"fmt"
	"log"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/johna210/go-next-flutter/internal/config"
	"github.com/johna210/go-next-flutter/internal/delivery/http"
	"github.com/johna210/go-next-flutter/internal/delivery/http/handler"
	"github.com/johna210/go-next-flutter/internal/infrastructure/memory"
	"github.com/johna210/go-next-flutter/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize dependencies
	userRepo := memory.NewUserRepository()
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	// Setup router
	router, api := http.SetupRouter(userHandler)

	// Add CORS middleware to router (before any routes)
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // 12 hours
	})
	router.Use(corsMiddleware)

	// Print available endpoints
	log.Println("🚀 Server starting...")
	log.Printf("📝 API Documentation: http://%s:%s/docs", cfg.Server.Host, cfg.Server.Port)
	log.Printf("📄 OpenAPI Spec (JSON): http://%s:%s/openapi.json", cfg.Server.Host, cfg.Server.Port)
	log.Printf("📄 OpenAPI Spec (YAML): http://%s:%s/openapi.yaml", cfg.Server.Host, cfg.Server.Port)
	log.Printf("💚 Health Check: http://%s:%s/health", cfg.Server.Host, cfg.Server.Port)

	// Print registered operations
	log.Println("\n📚 Registered API Operations:")
	for path, op := range api.OpenAPI().Paths {
		// Define possible methods and their corresponding fields.
		type methodOp struct {
			name string
			op   *huma.Operation
		}

		var methodOps []methodOp
		if op.Get != nil {
			methodOps = append(methodOps, methodOp{"GET", op.Get})
		}
		if op.Put != nil {
			methodOps = append(methodOps, methodOp{"PUT", op.Put})
		}
		if op.Post != nil {
			methodOps = append(methodOps, methodOp{"POST", op.Post})
		}
		if op.Delete != nil {
			methodOps = append(methodOps, methodOp{"DELETE", op.Delete})
		}
		if op.Options != nil {
			methodOps = append(methodOps, methodOp{"OPTIONS", op.Options})
		}
		if op.Head != nil {
			methodOps = append(methodOps, methodOp{"HEAD", op.Head})
		}
		if op.Patch != nil {
			methodOps = append(methodOps, methodOp{"PATCH", op.Patch})
		}
		if op.Trace != nil {
			methodOps = append(methodOps, methodOp{"TRACE", op.Trace})
		}

		for _, mo := range methodOps {
			log.Printf("  %s %s - %s", mo.name, path, mo.op.Summary)
		}
	}

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("\n✅ Server running on http://%s\n", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
