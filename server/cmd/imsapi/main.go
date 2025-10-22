package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"ims-intro/pkg/backup"
	"ims-intro/pkg/common/app"
	"ims-intro/pkg/common/sqlite"
	"ims-intro/pkg/controller"
	"ims-intro/pkg/middleware"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	// Try to load .env file, but don't fail if it doesn't exist
	// (environment variables may be set directly, e.g., in Docker)
	envPath := filepath.Join("..", "..", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	configurationManager := app.NewConfigurationManager()
	db := sqlite.GetConnection(configurationManager.SqliteConfig)

	// Initialize backup service
	backupConfig := backup.BackupConfig{
		SourcePath:     configurationManager.SqliteConfig.DatabasePath,
		BackupDir:      filepath.Join("data", "backups"),
		RetentionDays:  7,                   // Keep backups for 7 days
		BackupInterval: 24 * time.Hour,      // Backup every 24 hours
	}
	backupService := backup.NewBackupService(backupConfig)
	backupService.Start()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productController := controller.NewProductController(productService)

	categoryRepository := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryController := controller.NewCategoryController(categoryService)

	settingRepository := repository.NewSettingRepository(db)
	settingService := service.NewSettingService(settingRepository)
	settingController := controller.NewSettingController(settingService)

	transactionRepository := repository.NewTransactionRepository(db)
	saleRepository := repository.NewSaleRepository(db)
	saleService := service.NewSaleService(saleRepository, productRepository, transactionRepository)
	receiptService := service.NewReceiptService(productRepository, settingRepository)
	saleController := controller.NewSaleController(saleService, receiptService)

	backupController := controller.NewBackupController(backupService)

	e := echo.New()

	userController.RegisterUserRoutes(e, middleware.AuthMiddleware)
	productController.RegisterProductRoutes(e)
	categoryController.RegisterCategoryRoutes(e)
	saleController.RegisterSaleRoutes(e)
	settingController.RegisterSettingRoutes(e)
	backupController.RegisterBackupRoutes(e, middleware.AuthMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Server is running on port", port)
		if err := e.Start(":" + port); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server gracefully...")

	// Stop backup service (creates final backup)
	log.Println("Creating final backup before shutdown...")
	backupService.Stop()
	log.Println("Backup completed")

	// Close database connection
	log.Println("Closing database connection...")
	db.Close()
	log.Println("Server shutdown complete")
}
