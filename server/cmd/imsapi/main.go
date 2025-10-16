package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"ims-intro/pkg/common/app"
	"ims-intro/pkg/common/sqlite"
	"ims-intro/pkg/controller"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service"
	"log"
	"os"
	"path/filepath"
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
	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productController := controller.NewProductController(productService)

	e := echo.New()
	userController.RegisterUserRoutes(e)
	productController.RegisterProductRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server is running on port", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
