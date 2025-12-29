package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ims-intro/pkg/backup"
	"ims-intro/pkg/common/app"
	"ims-intro/pkg/common/sqlite"
	"ims-intro/pkg/controller"
	"ims-intro/pkg/middleware"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service"

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
	squareRepository := repository.NewSquareRepository(a.db)
	terminalRepository := repository.NewTerminalRepository(a.db)

	// Initialize services
	userService := service.NewUserService(userRepository)
	productService := service.NewProductService(productRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	settingService := service.NewSettingService(settingRepository)
	salesReportService := service.NewSalesReportService(a.db)
	productAnalyticsService := service.NewProductAnalyticsService(a.db)
	abcAnalysisService := service.NewABCAnalysisService(a.db)
	inventoryTurnoverService := service.NewInventoryTurnoverService(a.db)
	squareService := service.NewSquareService(squareRepository)
	terminalService := service.NewTerminalService(terminalRepository, squareRepository, saleRepository)
	saleService := service.NewSaleService(saleRepository, productRepository, transactionRepository, terminalService)
	receiptService := service.NewReceiptService(productRepository, settingRepository)

	// Initialize controllers
	userController := controller.NewUserController(userService)
	productController := controller.NewProductController(productService)
	categoryController := controller.NewCategoryController(categoryService)
	settingController := controller.NewSettingController(settingService)
	saleController := controller.NewSaleController(saleService, receiptService)
	backupController := controller.NewBackupController(a.backupService)
	reportController := controller.NewReportController(
		salesReportService,
		productAnalyticsService,
		abcAnalysisService,
		inventoryTurnoverService,
	)
	squareController := controller.NewSquareController(squareService)
	terminalController := controller.NewTerminalController(terminalService)
	webhookController := controller.NewWebhookController(terminalService, saleService)

	// Initialize Echo server
	a.server = echo.New()
	a.server.HideBanner = false

	// Add detailed request logging middleware
	a.server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Log incoming request
			log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			log.Printf("→ [%s] %s", c.Request().Method, c.Request().URL.Path)
			log.Printf("  Headers: %v", c.Request().Header)
			log.Printf("  Remote: %s", c.Request().RemoteAddr)

			// Check for auth cookie
			cookie, err := c.Cookie("token")
			if err == nil && cookie != nil {
				log.Printf("  Auth Cookie: present (value length: %d)", len(cookie.Value))
			} else {
				log.Printf("  Auth Cookie: missing or invalid")
			}

			// Check for Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				log.Printf("  Authorization Header: present")
			}

			// Execute the handler
			err = next(c)

			// Log response status
			if err != nil {
				log.Printf("✗ Error: %v", err)
			} else {
				log.Printf("✓ Status: %d", c.Response().Status)
			}
			log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

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

	// Create /api subrouter by using a middleware to strip /api prefix
	apiEcho := echo.New()
	apiEcho.HideBanner = true

	// Register routes on apiEcho (without /api prefix)
	userController.RegisterUserRoutes(apiEcho, middleware.AuthMiddleware)
	productController.RegisterProductRoutes(apiEcho)
	categoryController.RegisterCategoryRoutes(apiEcho)
	saleController.RegisterSaleRoutes(apiEcho, middleware.AuthMiddleware, middleware.AdminOnlyMiddleware)
	settingController.RegisterSettingRoutes(apiEcho)
	backupController.RegisterBackupRoutes(apiEcho, middleware.AuthMiddleware)
	reportController.RegisterReportRoutes(apiEcho)
	squareController.RegisterSquareRoutes(apiEcho)
	terminalController.RegisterTerminalRoutes(apiEcho)
	webhookController.RegisterWebhookRoutes(apiEcho)

	// Mount apiEcho under /api/* on main server
	a.server.Any("/api/*", func(c echo.Context) error {
		req := c.Request()
		// Strip /api prefix from path
		req.URL.Path = req.URL.Path[4:] // Remove "/api"
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}
		apiEcho.ServeHTTP(c.Response(), req)
		return nil
	})

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
