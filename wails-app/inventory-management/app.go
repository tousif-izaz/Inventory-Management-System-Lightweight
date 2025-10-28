package main

import (
	"context"
	"database/sql"
	"fmt"
	"inventory-management/pkg/backup"
	"inventory-management/pkg/common/app"
	"inventory-management/pkg/common/sqlite"
	"inventory-management/pkg/controller"
	"inventory-management/pkg/middleware"
	"inventory-management/pkg/repository"
	"inventory-management/pkg/service"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

// App struct
type App struct {
	ctx           context.Context
	db            *sql.DB
	backupService *backup.BackupService
	server        *echo.Echo
	serverPort    string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		serverPort: "37285", // Random port to avoid conflicts
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Get app data directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Failed to get home directory: %v", err)
		homeDir = "."
	}

	// Create app data directory
	appDataDir := filepath.Join(homeDir, ".inventory-management")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		log.Printf("Failed to create app data directory: %v", err)
	}

	dbPath := filepath.Join(appDataDir, "inventory.db")
	log.Printf("Database path: %s", dbPath)

	// Initialize configuration
	configurationManager := app.NewConfigurationManager()
	configurationManager.SqliteConfig.DatabasePath = dbPath

	// Initialize database connection
	a.db = sqlite.GetConnection(configurationManager.SqliteConfig)

	// Initialize backup service
	backupDir := filepath.Join(appDataDir, "backups")
	backupConfig := backup.BackupConfig{
		SourcePath:     dbPath,
		BackupDir:      backupDir,
		RetentionDays:  7,
		BackupInterval: 24 * time.Hour,
	}
	a.backupService = backup.NewBackupService(backupConfig)
	a.backupService.Start()

	// Initialize repositories
	userRepository := repository.NewUserRepository(a.db)
	productRepository := repository.NewProductRepository(a.db)
	categoryRepository := repository.NewCategoryRepository(a.db)
	settingRepository := repository.NewSettingRepository(a.db)
	transactionRepository := repository.NewTransactionRepository(a.db)
	saleRepository := repository.NewSaleRepository(a.db)

	// Initialize services
	userService := service.NewUserService(userRepository)
	productService := service.NewProductService(productRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	settingService := service.NewSettingService(settingRepository)
	saleService := service.NewSaleService(saleRepository, productRepository, transactionRepository)
	receiptService := service.NewReceiptService(productRepository, settingRepository)

	// Initialize controllers
	userController := controller.NewUserController(userService)
	productController := controller.NewProductController(productService)
	categoryController := controller.NewCategoryController(categoryService)
	settingController := controller.NewSettingController(settingService)
	saleController := controller.NewSaleController(saleService, receiptService)
	backupController := controller.NewBackupController(a.backupService)

	// Initialize Echo server
	a.server = echo.New()
	a.server.HideBanner = true

	// Add request logging middleware
	a.server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Printf("[%s] %s %s", c.Request().Method, c.Request().URL.Path, c.Request().RemoteAddr)
			err := next(c)
			if err != nil {
				log.Printf("Error handling request: %v", err)
			}
			return err
		}
	})

	// Add CORS middleware for local development
	a.server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get("Origin")
			if origin == "" {
				origin = "*"
			}
			c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusNoContent)
			}
			return next(c)
		}
	})

	// Add health check endpoint
	a.server.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"server": "running",
		})
	})

	// Create /api group to match frontend expectations (Vite proxy was stripping /api prefix)
	apiGroup := a.server.Group("/api")

	// Register user routes (auth routes)
	apiGroup.POST("/login", userController.Login)
	apiGroup.POST("/signup", userController.SignUp)
	apiGroup.POST("/logout", userController.Logout)
	apiGroup.GET("/profile", userController.GetProfile, middleware.AuthMiddleware)
	apiGroup.PUT("/profile/password", userController.UpdatePassword, middleware.AuthMiddleware)

	// Register product routes
	productsGroup := apiGroup.Group("/products")
	productsGroup.Use(middleware.AuthMiddleware)
	productsGroup.GET("", productController.GetAllProducts)
	productsGroup.GET("/sku/:sku", productController.GetProductBySKU)
	productsGroup.POST("", productController.AddNewProduct)
	productsGroup.PUT("/:id", productController.UpdateProductById)
	productsGroup.DELETE("/:id", productController.DeleteProductById)

	// Register category routes
	categoriesGroup := apiGroup.Group("/categories")
	categoriesGroup.Use(middleware.AuthMiddleware)
	categoriesGroup.GET("", categoryController.GetAllCategories)
	categoriesGroup.POST("", categoryController.AddNewCategory)
	categoriesGroup.PUT("/:id", categoryController.UpdateCategoryById)
	categoriesGroup.DELETE("/:id", categoryController.DeleteCategoryById)

	// Register sale routes
	salesGroup := apiGroup.Group("/sales")
	salesGroup.Use(middleware.AuthMiddleware)
	salesGroup.GET("", saleController.GetAllSales)
	salesGroup.POST("", saleController.CreateSale)
	// More specific routes must come before generic :id route
	salesGroup.GET("/:id/receipt", saleController.GenerateReceipt)
	salesGroup.PUT("/:id/payment-status", saleController.UpdatePaymentStatus)
	salesGroup.POST("/:id/refund", saleController.ProcessRefund)
	// Generic :id route comes last
	salesGroup.GET("/:id", saleController.GetSaleByID)

	// Register setting routes
	settingsGroup := apiGroup.Group("/settings")
	settingsGroup.Use(middleware.AuthMiddleware)
	settingsGroup.GET("", settingController.GetAllSettings)
	settingsGroup.PUT("/:key", settingController.UpdateSetting)

	// Register backup routes
	backupGroup := apiGroup.Group("/backups")
	backupGroup.Use(middleware.AuthMiddleware)
	backupGroup.POST("", backupController.CreateBackup)
	backupGroup.GET("", backupController.ListBackups)

	// Start HTTP server in goroutine
	go func() {
		log.Printf("Starting HTTP server on http://localhost:%s", a.serverPort)
		log.Printf("Health check endpoint: http://localhost:%s/health", a.serverPort)
		log.Printf("API endpoints available at: http://localhost:%s/api/*", a.serverPort)

		if err := a.server.Start(":" + a.serverPort); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ HTTP server error: %v", err)
			log.Printf("Port %s may be in use. Try closing other instances of the app.", a.serverPort)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test if server is accessible
	testURL := fmt.Sprintf("http://localhost:%s/health", a.serverPort)
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(testURL)
	if err != nil {
		log.Printf("⚠️  Warning: Server may not have started properly: %v", err)
	} else {
		resp.Body.Close()
		log.Printf("✅ Server health check passed - server is running")
	}

	log.Println("Application started successfully")
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	log.Println("Shutting down application...")

	// Stop HTTP server
	if a.server != nil {
		log.Println("Stopping HTTP server...")
		a.server.Shutdown(ctx)
	}

	// Stop backup service
	if a.backupService != nil {
		log.Println("Creating final backup...")
		a.backupService.Stop()
	}

	// Close database connection
	if a.db != nil {
		log.Println("Closing database connection...")
		a.db.Close()
	}

	log.Println("Application shutdown complete")
}

// GetServerPort returns the port the HTTP server is running on
func (a *App) GetServerPort() string {
	return a.serverPort
}

// GetAppVersion returns the application version
func (a *App) GetAppVersion() string {
	return "1.0.0"
}

// GetDatabasePath returns the database file path
func (a *App) GetDatabasePath() string {
	homeDir, _ := os.UserHomeDir()
	appDataDir := filepath.Join(homeDir, ".inventory-management")
	return filepath.Join(appDataDir, "inventory.db")
}

// Greet returns a greeting (for testing)
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// TestDatabaseConnection tests the database connection and returns diagnostic info
func (a *App) TestDatabaseConnection() string {
	if a.db == nil {
		return "❌ Database connection is nil"
	}

	// Test ping
	if err := a.db.Ping(); err != nil {
		return fmt.Sprintf("❌ Database ping failed: %v", err)
	}

	// Get table counts
	var result string
	result += "✅ Database connection OK\n"
	result += fmt.Sprintf("📁 Database path: %s\n\n", a.GetDatabasePath())

	tables := []string{"Users", "Products", "Categories", "Sales", "Purchases", "Customers", "Suppliers", "Settings"}
	result += "Table Counts:\n"
	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		err := a.db.QueryRow(query).Scan(&count)
		if err != nil {
			result += fmt.Sprintf("  ❌ %s: error - %v\n", table, err)
		} else {
			result += fmt.Sprintf("  ✅ %s: %d rows\n", table, count)
		}
	}

	return result
}

// CreateSampleData creates sample categories and products for testing
func (a *App) CreateSampleData() string {
	if a.db == nil {
		return "❌ Database connection is nil"
	}

	var result string
	result += "Creating sample data...\n\n"

	// Create sample categories
	categories := []struct {
		name        string
		description string
	}{
		{"Electronics", "Electronic devices and accessories"},
		{"Furniture", "Office and home furniture"},
		{"Stationery", "Office supplies and stationery"},
	}

	result += "Creating categories:\n"
	for _, cat := range categories {
		_, err := a.db.Exec(`
			INSERT INTO Categories (Name, Description, CreatedAt)
			VALUES (?, ?, datetime('now'))
		`, cat.name, cat.description)
		if err != nil {
			result += fmt.Sprintf("  ❌ Failed to create category '%s': %v\n", cat.name, err)
		} else {
			result += fmt.Sprintf("  ✅ Created category: %s\n", cat.name)
		}
	}

	// Get category IDs
	var electronicsID, furnitureID, stationeryID int64
	a.db.QueryRow("SELECT CategoryID FROM Categories WHERE Name = 'Electronics'").Scan(&electronicsID)
	a.db.QueryRow("SELECT CategoryID FROM Categories WHERE Name = 'Furniture'").Scan(&furnitureID)
	a.db.QueryRow("SELECT CategoryID FROM Categories WHERE Name = 'Stationery'").Scan(&stationeryID)

	// Create sample products
	products := []struct {
		name            string
		description     string
		sku             string
		categoryID      int64
		costPrice       float64
		sellingPrice    float64
		currentQuantity int
		minStockLevel   int
		reorderPoint    int
		unit            string
	}{
		{"Laptop HP ProBook", "15-inch business laptop", "LAP-HP-001", electronicsID, 500.00, 750.00, 10, 2, 5, "pcs"},
		{"Wireless Mouse", "Ergonomic wireless mouse", "MOU-WIR-001", electronicsID, 15.00, 25.00, 50, 10, 15, "pcs"},
		{"Office Desk", "Adjustable standing desk", "DSK-OFF-001", furnitureID, 200.00, 350.00, 5, 1, 2, "pcs"},
		{"Office Chair", "Ergonomic office chair", "CHR-OFF-001", furnitureID, 100.00, 180.00, 8, 2, 3, "pcs"},
		{"Notebook A4", "Ruled notebook 200 pages", "NOT-A4-001", stationeryID, 2.00, 5.00, 100, 20, 30, "pcs"},
		{"Pen Blue", "Blue ballpoint pen", "PEN-BLU-001", stationeryID, 0.50, 1.50, 200, 50, 75, "pcs"},
	}

	result += "\nCreating products:\n"
	for _, prod := range products {
		_, err := a.db.Exec(`
			INSERT INTO Products (
				Name, Description, SKU, CategoryID, CostPrice, SellingPrice,
				CurrentQuantity, MinStockLevel, ReorderPoint, Unit, IsActive,
				CreatedAt, UpdatedAt
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, datetime('now'), datetime('now'))
		`, prod.name, prod.description, prod.sku, prod.categoryID, prod.costPrice,
			prod.sellingPrice, prod.currentQuantity, prod.minStockLevel, prod.reorderPoint, prod.unit)

		if err != nil {
			result += fmt.Sprintf("  ❌ Failed to create product '%s': %v\n", prod.name, err)
		} else {
			result += fmt.Sprintf("  ✅ Created product: %s (SKU: %s)\n", prod.name, prod.sku)
		}
	}

	result += "\n✅ Sample data creation complete!"
	return result
}

// TestServerConnection tests if the HTTP server is responding
func (a *App) TestServerConnection() string {
	url := fmt.Sprintf("http://localhost:%s/api/products", a.serverPort)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Sprintf("❌ Server connection failed: %v", err)
	}
	defer resp.Body.Close()

	return fmt.Sprintf("✅ Server is responding\nStatus: %d %s\nURL: %s", resp.StatusCode, resp.Status, url)
}
