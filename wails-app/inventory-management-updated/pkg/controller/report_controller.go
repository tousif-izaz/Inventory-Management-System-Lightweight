package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"ims-intro/pkg/controller/response"
	"ims-intro/pkg/middleware"
	"ims-intro/pkg/service"
	"ims-intro/pkg/service/dto"
)

type ReportController struct {
	salesReportService       *service.SalesReportService
	productAnalyticsService  *service.ProductAnalyticsService
	abcAnalysisService       *service.ABCAnalysisService
	inventoryTurnoverService *service.InventoryTurnoverService
}

func NewReportController(
	salesReportService *service.SalesReportService,
	productAnalyticsService *service.ProductAnalyticsService,
	abcAnalysisService *service.ABCAnalysisService,
	inventoryTurnoverService *service.InventoryTurnoverService,
) *ReportController {
	return &ReportController{
		salesReportService:       salesReportService,
		productAnalyticsService:  productAnalyticsService,
		abcAnalysisService:       abcAnalysisService,
		inventoryTurnoverService: inventoryTurnoverService,
	}
}

func (controller *ReportController) RegisterReportRoutes(e *echo.Echo) {
	reportsGroup := e.Group("/reports")
	reportsGroup.Use(middleware.AuthMiddleware)

	// Sales Reports
	reportsGroup.GET("/sales/daily", controller.GetDailySalesReport)
	reportsGroup.GET("/sales/weekly", controller.GetWeeklySalesReport)
	reportsGroup.GET("/sales/monthly", controller.GetMonthlySalesReport)
	reportsGroup.GET("/sales/quarterly", controller.GetQuarterlySalesReport)
	reportsGroup.GET("/sales/custom", controller.GetCustomSalesReport)
	reportsGroup.GET("/sales/recent", controller.GetRecentSales)

	// Product Analytics
	reportsGroup.GET("/products/best-sellers", controller.GetBestSellingProducts)
	reportsGroup.GET("/products/worst-performers", controller.GetWorstPerformingProducts)
	reportsGroup.GET("/categories/performance", controller.GetCategoryPerformance)

	// ABC Analysis
	reportsGroup.GET("/abc-analysis", controller.GetABCAnalysis)
	reportsGroup.GET("/abc-analysis/products", controller.GetABCProductsByClass)

	// Inventory Turnover
	reportsGroup.GET("/inventory/turnover", controller.GetInventoryTurnover)
	reportsGroup.GET("/inventory/movement", controller.GetProductMovement)
	reportsGroup.GET("/inventory/stock-health", controller.GetStockHealth)
}

// ========== Sales Report Endpoints ==========

func (controller *ReportController) GetDailySalesReport(c echo.Context) error {
	dateStr := c.QueryParam("date")
	if dateStr == "" {
		// Default to today
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid date format. Use YYYY-MM-DD"))
	}

	report, err := controller.salesReportService.GetDailySalesReport(date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, report)
}

func (controller *ReportController) GetWeeklySalesReport(c echo.Context) error {
	yearStr := c.QueryParam("year")
	weekStr := c.QueryParam("week")

	if yearStr == "" || weekStr == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Year and week parameters are required"))
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid year format"))
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 || week > 53 {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid week format (must be 1-53)"))
	}

	report, err := controller.salesReportService.GetWeeklySalesReport(year, week)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, report)
}

func (controller *ReportController) GetMonthlySalesReport(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	if yearStr == "" || monthStr == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Year and month parameters are required"))
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid year format"))
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid month format (must be 1-12)"))
	}

	report, err := controller.salesReportService.GetMonthlySalesReport(year, month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, report)
}

func (controller *ReportController) GetQuarterlySalesReport(c echo.Context) error {
	yearStr := c.QueryParam("year")
	quarterStr := c.QueryParam("quarter")

	if yearStr == "" || quarterStr == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Year and quarter parameters are required"))
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid year format"))
	}

	quarter, err := strconv.Atoi(quarterStr)
	if err != nil || quarter < 1 || quarter > 4 {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid quarter format (must be 1-4)"))
	}

	report, err := controller.salesReportService.GetQuarterlySalesReport(year, quarter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, report)
}

func (controller *ReportController) GetCustomSalesReport(c echo.Context) error {
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	log.Printf("📊 [GetCustomSalesReport] Received request - start_date: %s, end_date: %s", startDateStr, endDateStr)

	if startDateStr == "" || endDateStr == "" {
		log.Printf("❌ [GetCustomSalesReport] Missing date parameters")
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("start_date and end_date parameters are required"))
	}

	// Try parsing as ISO 8601 first, then fall back to YYYY-MM-DD
	var startDate, endDate time.Time
	var err error

	// Try ISO 8601 format (with time)
	startDate, err = time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		// Fall back to date-only format
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			log.Printf("❌ [GetCustomSalesReport] Failed to parse start_date: %v", err)
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid start_date format. Use ISO 8601 or YYYY-MM-DD"))
		}
	}

	endDate, err = time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		// Fall back to date-only format
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			log.Printf("❌ [GetCustomSalesReport] Failed to parse end_date: %v", err)
			return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid end_date format. Use ISO 8601 or YYYY-MM-DD"))
		}
		// If only date provided, set to end of day
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
	}

	log.Printf("📅 [GetCustomSalesReport] Parsed dates - Start: %s, End: %s", startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))

	req := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Parse optional filters
	if customerIDStr := c.QueryParam("customer_id"); customerIDStr != "" {
		if customerID, err := strconv.ParseInt(customerIDStr, 10, 64); err == nil {
			req.CustomerID = &customerID
		}
	}

	if paymentStatus := c.QueryParam("payment_status"); paymentStatus != "" {
		req.PaymentStatus = &paymentStatus
	}

	if paymentMethod := c.QueryParam("payment_method"); paymentMethod != "" {
		req.PaymentMethod = &paymentMethod
	}

	log.Printf("🔍 [GetCustomSalesReport] Calling sales report service...")
	report, err := controller.salesReportService.GetSalesReport(req)
	if err != nil {
		log.Printf("❌ [GetCustomSalesReport] Service error: %v", err)
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	log.Printf("✅ [GetCustomSalesReport] Report generated successfully - Transactions: %d, Revenue: %.2f", report.TransactionCount, report.TotalRevenue)
	return c.JSON(http.StatusOK, report)
}

func (controller *ReportController) GetRecentSales(c echo.Context) error {
	// Parse limit parameter
	limitStr := c.QueryParam("limit")
	limit := 10 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	log.Printf("📋 [GetRecentSales] Received request - limit: %d", limit)

	// Parse optional filters
	var filters *dto.SalesReportRequest
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	log.Printf("📋 [GetRecentSales] Date params - start_date: %s, end_date: %s", startDateStr, endDateStr)

	if startDateStr != "" && endDateStr != "" {
		// Try ISO 8601 format first, then fall back to YYYY-MM-DD
		var startDate, endDate time.Time
		var err error

		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				log.Printf("❌ [GetRecentSales] Failed to parse start_date: %v", err)
				return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid start_date format"))
			}
		}

		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				log.Printf("❌ [GetRecentSales] Failed to parse end_date: %v", err)
				return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid end_date format"))
			}
			// If only date provided, set to end of day
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
		}

		log.Printf("📅 [GetRecentSales] Parsed dates - Start: %s, End: %s", startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))

		filters = &dto.SalesReportRequest{
			StartDate: startDate,
			EndDate:   endDate,
		}

		// Apply other optional filters
		if paymentMethod := c.QueryParam("payment_method"); paymentMethod != "" {
			filters.PaymentMethod = &paymentMethod
			log.Printf("🔍 [GetRecentSales] Filter by payment_method: %s", paymentMethod)
		}

		if paymentStatus := c.QueryParam("payment_status"); paymentStatus != "" {
			filters.PaymentStatus = &paymentStatus
			log.Printf("🔍 [GetRecentSales] Filter by payment_status: %s", paymentStatus)
		}
	}

	log.Printf("🔍 [GetRecentSales] Calling GetRecentSales service...")
	sales, err := controller.salesReportService.GetRecentSales(limit, filters)
	if err != nil {
		log.Printf("❌ [GetRecentSales] Service error: %v", err)
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	log.Printf("✅ [GetRecentSales] Retrieved %d sales", len(sales))
	return c.JSON(http.StatusOK, sales)
}

// ========== Product Analytics Endpoints ==========

func (controller *ReportController) GetBestSellingProducts(c echo.Context) error {
	req, err := controller.parseProductPerformanceRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	products, err := controller.productAnalyticsService.GetBestSellingProducts(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, products)
}

func (controller *ReportController) GetWorstPerformingProducts(c echo.Context) error {
	req, err := controller.parseProductPerformanceRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	products, err := controller.productAnalyticsService.GetWorstPerformingProducts(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, products)
}

func (controller *ReportController) GetCategoryPerformance(c echo.Context) error {
	req, err := controller.parseProductPerformanceRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	categories, err := controller.productAnalyticsService.GetCategoryPerformance(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, categories)
}

// ========== ABC Analysis Endpoints ==========

func (controller *ReportController) GetABCAnalysis(c echo.Context) error {
	req, err := controller.parseABCAnalysisRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	analysis, err := controller.abcAnalysisService.PerformABCAnalysis(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, analysis)
}

func (controller *ReportController) GetABCProductsByClass(c echo.Context) error {
	class := c.QueryParam("class")
	if class == "" || (class != "A" && class != "B" && class != "C") {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("class parameter is required (A, B, or C)"))
	}

	req, err := controller.parseABCAnalysisRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	products, err := controller.abcAnalysisService.GetProductsByClass(req, class)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, products)
}

// ========== Inventory Turnover Endpoints ==========

func (controller *ReportController) GetInventoryTurnover(c echo.Context) error {
	metrics, err := controller.inventoryTurnoverService.GetInventoryTurnoverMetrics()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, metrics)
}

func (controller *ReportController) GetProductMovement(c echo.Context) error {
	movements, err := controller.inventoryTurnoverService.GetProductMovementAnalysis()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, movements)
}

func (controller *ReportController) GetStockHealth(c echo.Context) error {
	health, err := controller.inventoryTurnoverService.GetStockHealthReport()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, health)
}

// ========== Helper Methods ==========

func (controller *ReportController) parseProductPerformanceRequest(c echo.Context) (*dto.ProductPerformanceRequest, error) {
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	// Default to last 30 days if not specified
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	if startDateStr != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, err
		}
	}

	if endDateStr != "" {
		var err error
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, err
		}
	}

	req := &dto.ProductPerformanceRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     10, // Default limit
	}

	// Parse optional filters
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}

	if categoryIDStr := c.QueryParam("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			req.CategoryID = &categoryID
		}
	}

	if supplierIDStr := c.QueryParam("supplier_id"); supplierIDStr != "" {
		if supplierID, err := strconv.ParseInt(supplierIDStr, 10, 64); err == nil {
			req.SupplierID = &supplierID
		}
	}

	return req, nil
}

func (controller *ReportController) parseABCAnalysisRequest(c echo.Context) (*dto.ABCAnalysisRequest, error) {
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	// Default to last 90 days if not specified
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -90)

	if startDateStr != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, err
		}
	}

	if endDateStr != "" {
		var err error
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, err
		}
	}

	req := &dto.ABCAnalysisRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Parse optional filters
	if categoryIDStr := c.QueryParam("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			req.CategoryID = &categoryID
		}
	}

	if supplierIDStr := c.QueryParam("supplier_id"); supplierIDStr != "" {
		if supplierID, err := strconv.ParseInt(supplierIDStr, 10, 64); err == nil {
			req.SupplierID = &supplierID
		}
	}

	return req, nil
}
