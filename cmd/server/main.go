// backend/cmd/server/main.go
package main

import (
	"log"
	"net/http"

	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/config"
	"backend/pkg/database"
)

func main() {
	log.Println("Starting AminNCo API Server...")

	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: could not load configuration: %v", err)
	}

	// 2. Connect to Database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("FATAL: could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection successful.")

	// 3. Initialize Repositories (Database Layer)
	userRepo := repository.NewPostgresUserRepository(db)
	productRepo := repository.NewPostgresProductRepository(db)
	storeRepo := repository.NewPostgresStoreRepository(db)
	adminRepo := repository.NewPostgresAdminRepository(db)
	// logRepo := repository.NewPostgresLogRepository(db) // For later

	// 4. Initialize Services (Business Logic Layer)
	// Services can be composed of multiple repositories.
	userService := service.NewUserService(userRepo)
	catalogService := service.NewCatalogService(productRepo)
	storeService := service.NewStoreService(storeRepo)
	adminService := service.NewAdminService(adminRepo)

	// 5. Initialize Handlers (HTTP Layer)
	userHandler := handler.NewUserHandler(userService)
	catalogHandler := handler.NewCatalogHandler(catalogService)
	storeHandler := handler.NewStoreHandler(storeService)
	adminHandler := handler.NewAdminHandler(adminService)

	// 6. Setup Router and Server, injecting all handlers
	router := handler.NewRouter(
		userHandler,
		catalogHandler,
		storeHandler,
		adminHandler,
	)

	// 7. Start the server
	log.Printf("Server starting on http://localhost:%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatalf("FATAL: could not start server: %v", err)
	}
}
