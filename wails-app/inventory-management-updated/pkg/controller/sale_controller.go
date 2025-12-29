package controller

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"ims-intro/pkg/controller/request"
	"ims-intro/pkg/controller/response"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/service"
	"net/http"
	"strconv"
	"time"
)

type SaleController struct {
	saleService    *service.SaleService
	receiptService *service.ReceiptService
	validator      *validator.Validate
}

func NewSaleController(saleService *service.SaleService, receiptService *service.ReceiptService) *SaleController {
	return &SaleController{
		saleService:    saleService,
		receiptService: receiptService,
		validator:      validator.New(),
	}
}

// CreateSale handles POST /sales
func (c *SaleController) CreateSale(ctx echo.Context) error {
	var req request.CreateSaleRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error_message": "Invalid request body",
		})
	}

	if err := c.validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error_message": err.Error(),
		})
	}

	// For now, use userID = 1 (in production, get from JWT token)
	userID := int64(1)

	saleWithItems, err := c.saleService.CreateSale(req.ToDTO(), userID)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error_message": err.Error(),
		})
	}

	// Convert to response
	saleResp := c.toSaleWithItemsResponse(saleWithItems)

	return ctx.JSON(http.StatusCreated, saleResp)
}

// GetAllSales handles GET /sales
func (c *SaleController) GetAllSales(ctx echo.Context) error {
	// Parse query params
	var customerID *int64
	if customerIDStr := ctx.QueryParam("customer_id"); customerIDStr != "" {
		id, err := strconv.ParseInt(customerIDStr, 10, 64)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error_message": "Invalid customer_id",
			})
		}
		customerID = &id
	}

	var paymentStatus *string
	if status := ctx.QueryParam("payment_status"); status != "" {
		paymentStatus = &status
	}

	var fromDate, toDate *time.Time
	if fromDateStr := ctx.QueryParam("from_date"); fromDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, fromDateStr)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error_message": "Invalid from_date format",
			})
		}
		fromDate = &parsed
	}

	if toDateStr := ctx.QueryParam("to_date"); toDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, toDateStr)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error_message": "Invalid to_date format",
			})
		}
		toDate = &parsed
	}

	sales, err := c.saleService.GetAllSales(customerID, paymentStatus, fromDate, toDate)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": err.Error(),
		})
	}

	// Convert to response
	salesResp := make([]response.SaleResponse, len(sales))
	for i, sale := range sales {
		salesResp[i] = c.toSaleResponse(&sale)
	}

	return ctx.JSON(http.StatusOK, salesResp)
}

// GetSaleByID handles GET /sales/:id
func (c *SaleController) GetSaleByID(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error_message": "Invalid sale ID",
		})
	}

	saleWithItems, err := c.saleService.GetSaleByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error_message": err.Error(),
		})
	}

	// Convert to response
	saleResp := c.toSaleWithItemsResponse(saleWithItems)

	return ctx.JSON(http.StatusOK, saleResp)
}

// UpdatePaymentStatus handles PUT /sales/:id/payment-status
func (c *SaleController) UpdatePaymentStatus(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error_message": "Invalid sale ID",
		})
	}

	var req request.UpdatePaymentStatusRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error_message": "Invalid request body",
		})
	}

	if err := c.validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error_message": err.Error(),
		})
	}

	err = c.saleService.UpdatePaymentStatus(id, req.PaymentStatus)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error_message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Payment status updated successfully",
	})
}

// ProcessRefund handles POST /sales/:id/refund
func (c *SaleController) ProcessRefund(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error_message": "Invalid sale ID",
		})
	}

	// For now, use userID = 1 (in production, get from JWT token)
	userID := int64(1)

	err = c.saleService.ProcessRefund(id, userID)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error_message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Refund processed successfully",
	})
}

// GenerateReceipt handles GET /sales/:id/receipt
func (c *SaleController) GenerateReceipt(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error_message": "Invalid sale ID",
		})
	}

	// Get optional customer_name query parameter
	customerName := ctx.QueryParam("customer_name")

	// Fetch sale with items
	saleWithItems, err := c.saleService.GetSaleByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error_message": err.Error(),
		})
	}

	// Generate PDF receipt
	receiptData := &service.ReceiptData{
		SaleWithItems: saleWithItems,
		CustomerName:  customerName,
	}

	pdfBuffer, err := c.receiptService.GenerateReceipt(receiptData)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to generate receipt: " + err.Error(),
		})
	}

	// Set response headers for PDF display in browser
	ctx.Response().Header().Set("Content-Type", "application/pdf")
	ctx.Response().Header().Set("Content-Disposition", "inline; filename=\"receipt-"+saleWithItems.ReceiptNo+".pdf\"")

	return ctx.Blob(http.StatusOK, "application/pdf", pdfBuffer.Bytes())
}

// Helper methods to convert domain models to responses
func (c *SaleController) toSaleResponse(sale *domain.Sale) response.SaleResponse {
	return response.SaleResponse{
		SaleID:         sale.SaleID,
		SaleDate:       sale.SaleDate,
		ReceiptNo:      sale.ReceiptNo,
		CustomerID:     sale.CustomerID,
		TotalAmount:    sale.TotalAmount,
		TaxAmount:      sale.TaxAmount,
		DiscountAmount: sale.DiscountAmount,
		NetAmount:      sale.NetAmount,
		PaymentStatus:  sale.PaymentStatus,
		PaymentMethod:  sale.PaymentMethod,
		Notes:          sale.Notes,
		SoldBy:         sale.SoldBy,
		CreatedAt:      sale.CreatedAt,
	}
}

func (c *SaleController) toSaleItemResponse(item *domain.SaleItem) response.SaleItemResponse {
	return response.SaleItemResponse{
		SaleItemID:      item.SaleItemID,
		SaleID:          item.SaleID,
		ProductID:       item.ProductID,
		Quantity:        item.Quantity,
		UnitPrice:       item.UnitPrice,
		TaxRate:         item.TaxRate,
		DiscountPercent: item.DiscountPercent,
		LineTotal:       item.LineTotal,
	}
}

func (c *SaleController) toSaleWithItemsResponse(saleWithItems *domain.SaleWithItems) response.SaleWithItemsResponse {
	items := make([]response.SaleItemResponse, len(saleWithItems.Items))
	for i, item := range saleWithItems.Items {
		items[i] = c.toSaleItemResponse(&item)
	}

	return response.SaleWithItemsResponse{
		SaleResponse: c.toSaleResponse(&saleWithItems.Sale),
		Items:        items,
	}
}

// BackupSales handles GET /sales/backup - exports sales to CSV
func (c *SaleController) BackupSales(ctx echo.Context) error {
	// Parse query params for date filtering
	var fromDate, toDate *time.Time
	if fromDateStr := ctx.QueryParam("from_date"); fromDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, fromDateStr)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error_message": "Invalid from_date format",
			})
		}
		fromDate = &parsed
	}

	if toDateStr := ctx.QueryParam("to_date"); toDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, toDateStr)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error_message": "Invalid to_date format",
			})
		}
		toDate = &parsed
	}

	// Get all sales with filters
	sales, err := c.saleService.GetAllSales(nil, nil, fromDate, toDate)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to fetch sales: " + err.Error(),
		})
	}

	// Fetch all sale items for each sale
	type SaleWithItemsExport struct {
		Sale  domain.Sale
		Items []domain.SaleItem
	}
	salesWithItems := make([]SaleWithItemsExport, 0)

	for _, sale := range sales {
		saleWithItems, err := c.saleService.GetSaleByID(sale.SaleID)
		if err != nil {
			continue // Skip sales that can't be fetched
		}
		salesWithItems = append(salesWithItems, SaleWithItemsExport{
			Sale:  saleWithItems.Sale,
			Items: saleWithItems.Items,
		})
	}

	// Create CSV buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{
		"Sale ID", "Sale Date", "Receipt No", "Customer ID",
		"Total Amount", "Tax Amount", "Discount Amount", "Net Amount",
		"Payment Status", "Payment Method", "Notes", "Sold By",
		"Item ID", "Product ID", "Quantity", "Unit Price",
		"Tax Rate", "Discount Percent", "Line Total",
	}
	if err := writer.Write(header); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to write CSV header",
		})
	}

	// Write data rows
	for _, saleData := range salesWithItems {
		sale := saleData.Sale
		if len(saleData.Items) == 0 {
			// Write sale without items
			row := []string{
				strconv.FormatInt(sale.SaleID, 10),
				sale.SaleDate.Format(time.RFC3339),
				sale.ReceiptNo,
				formatInt64Ptr(sale.CustomerID),
				fmt.Sprintf("%.2f", sale.TotalAmount),
				fmt.Sprintf("%.2f", sale.TaxAmount),
				fmt.Sprintf("%.2f", sale.DiscountAmount),
				fmt.Sprintf("%.2f", sale.NetAmount),
				sale.PaymentStatus,
				formatStringPtr(sale.PaymentMethod),
				formatStringPtr(sale.Notes),
				formatInt64Ptr(sale.SoldBy),
				"", "", "", "", "", "", "",
			}
			if err := writer.Write(row); err != nil {
				continue
			}
		} else {
			// Write sale with each item
			for _, item := range saleData.Items {
				row := []string{
					strconv.FormatInt(sale.SaleID, 10),
					sale.SaleDate.Format(time.RFC3339),
					sale.ReceiptNo,
					formatInt64Ptr(sale.CustomerID),
					fmt.Sprintf("%.2f", sale.TotalAmount),
					fmt.Sprintf("%.2f", sale.TaxAmount),
					fmt.Sprintf("%.2f", sale.DiscountAmount),
					fmt.Sprintf("%.2f", sale.NetAmount),
					sale.PaymentStatus,
					formatStringPtr(sale.PaymentMethod),
					formatStringPtr(sale.Notes),
					formatInt64Ptr(sale.SoldBy),
					strconv.FormatInt(item.SaleItemID, 10),
					strconv.FormatInt(item.ProductID, 10),
					strconv.FormatInt(item.Quantity, 10),
					fmt.Sprintf("%.2f", item.UnitPrice),
					fmt.Sprintf("%.2f", item.TaxRate),
					fmt.Sprintf("%.2f", item.DiscountPercent),
					fmt.Sprintf("%.2f", item.LineTotal),
				}
				if err := writer.Write(row); err != nil {
					continue
				}
			}
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to generate CSV",
		})
	}

	// Generate filename with timestamp
	filename := fmt.Sprintf("sales_backup_%s.csv", time.Now().Format("20060102_150405"))

	// Set response headers
	ctx.Response().Header().Set("Content-Type", "text/csv")
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return ctx.Blob(http.StatusOK, "text/csv", buf.Bytes())
}

// DeleteSalesInRange handles DELETE /sales/range - deletes sales in a date range and returns backup
func (c *SaleController) DeleteSalesInRange(ctx echo.Context) error {
	// Parse query params for date filtering
	var fromDate, toDate *time.Time
	fromDateStr := ctx.QueryParam("from_date")
	toDateStr := ctx.QueryParam("to_date")

	ctx.Logger().Infof("DeleteSalesInRange called with from_date=%s, to_date=%s", fromDateStr, toDateStr)

	if fromDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, fromDateStr)
		if err != nil {
			ctx.Logger().Errorf("Failed to parse from_date: %v", err)
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error_message": "Invalid from_date format: " + err.Error(),
			})
		}
		fromDate = &parsed
	}

	if toDateStr != "" {
		parsed, err := time.Parse(time.RFC3339, toDateStr)
		if err != nil {
			ctx.Logger().Errorf("Failed to parse to_date: %v", err)
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error_message": "Invalid to_date format: " + err.Error(),
			})
		}
		toDate = &parsed
	}

	// Get all sales to be deleted
	ctx.Logger().Info("Fetching sales to be deleted...")
	sales, err := c.saleService.GetAllSales(nil, nil, fromDate, toDate)
	if err != nil {
		ctx.Logger().Errorf("Failed to fetch sales: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to fetch sales: " + err.Error(),
		})
	}

	ctx.Logger().Infof("Found %d sales to delete", len(sales))

	if len(sales) == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error_message": "No sales found in the specified date range",
		})
	}

	// Fetch all sale items for backup
	type SaleWithItemsExport struct {
		Sale  domain.Sale
		Items []domain.SaleItem
	}
	salesWithItems := make([]SaleWithItemsExport, 0)

	for _, sale := range sales {
		saleWithItems, err := c.saleService.GetSaleByID(sale.SaleID)
		if err != nil {
			continue
		}
		salesWithItems = append(salesWithItems, SaleWithItemsExport{
			Sale:  saleWithItems.Sale,
			Items: saleWithItems.Items,
		})
	}

	// Create CSV backup
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{
		"Sale ID", "Sale Date", "Receipt No", "Customer ID",
		"Total Amount", "Tax Amount", "Discount Amount", "Net Amount",
		"Payment Status", "Payment Method", "Notes", "Sold By",
		"Item ID", "Product ID", "Quantity", "Unit Price",
		"Tax Rate", "Discount Percent", "Line Total",
	}
	if err := writer.Write(header); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to write CSV header",
		})
	}

	// Write data rows
	for _, saleData := range salesWithItems {
		sale := saleData.Sale
		if len(saleData.Items) == 0 {
			row := []string{
				strconv.FormatInt(sale.SaleID, 10),
				sale.SaleDate.Format(time.RFC3339),
				sale.ReceiptNo,
				formatInt64Ptr(sale.CustomerID),
				fmt.Sprintf("%.2f", sale.TotalAmount),
				fmt.Sprintf("%.2f", sale.TaxAmount),
				fmt.Sprintf("%.2f", sale.DiscountAmount),
				fmt.Sprintf("%.2f", sale.NetAmount),
				sale.PaymentStatus,
				formatStringPtr(sale.PaymentMethod),
				formatStringPtr(sale.Notes),
				formatInt64Ptr(sale.SoldBy),
				"", "", "", "", "", "", "",
			}
			if err := writer.Write(row); err != nil {
				continue
			}
		} else {
			for _, item := range saleData.Items {
				row := []string{
					strconv.FormatInt(sale.SaleID, 10),
					sale.SaleDate.Format(time.RFC3339),
					sale.ReceiptNo,
					formatInt64Ptr(sale.CustomerID),
					fmt.Sprintf("%.2f", sale.TotalAmount),
					fmt.Sprintf("%.2f", sale.TaxAmount),
					fmt.Sprintf("%.2f", sale.DiscountAmount),
					fmt.Sprintf("%.2f", sale.NetAmount),
					sale.PaymentStatus,
					formatStringPtr(sale.PaymentMethod),
					formatStringPtr(sale.Notes),
					formatInt64Ptr(sale.SoldBy),
					strconv.FormatInt(item.SaleItemID, 10),
					strconv.FormatInt(item.ProductID, 10),
					strconv.FormatInt(item.Quantity, 10),
					fmt.Sprintf("%.2f", item.UnitPrice),
					fmt.Sprintf("%.2f", item.TaxRate),
					fmt.Sprintf("%.2f", item.DiscountPercent),
					fmt.Sprintf("%.2f", item.LineTotal),
				}
				if err := writer.Write(row); err != nil {
					continue
				}
			}
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to generate backup CSV",
		})
	}

	// Delete the sales
	ctx.Logger().Info("Deleting sales...")
	deletedCount, err := c.saleService.DeleteSalesInRange(fromDate, toDate)
	if err != nil {
		ctx.Logger().Errorf("Failed to delete sales: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to delete sales: " + err.Error(),
		})
	}

	ctx.Logger().Infof("Successfully deleted %d sales", deletedCount)

	// Generate filename with timestamp
	filename := fmt.Sprintf("sales_deleted_backup_%s.csv", time.Now().Format("20060102_150405"))

	// Set response headers
	ctx.Response().Header().Set("Content-Type", "text/csv")
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Response().Header().Set("X-Deleted-Count", strconv.Itoa(deletedCount))

	return ctx.Blob(http.StatusOK, "text/csv", buf.Bytes())
}

// Helper function to format int64 pointer
func formatInt64Ptr(val *int64) string {
	if val == nil {
		return ""
	}
	return strconv.FormatInt(*val, 10)
}

// Helper function to format string pointer
func formatStringPtr(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

// RegisterSaleRoutes registers all sale routes
func (c *SaleController) RegisterSaleRoutes(e *echo.Echo, authMiddleware, adminOnly echo.MiddlewareFunc) {
	sales := e.Group("/sales")

	sales.POST("", c.CreateSale)
	sales.GET("", c.GetAllSales)
	sales.GET("/:id", c.GetSaleByID)
	sales.GET("/:id/receipt", c.GenerateReceipt)
	sales.PUT("/:id/payment-status", c.UpdatePaymentStatus)
	sales.POST("/:id/refund", c.ProcessRefund)

	// Admin-only routes - requires both auth and admin check
	sales.GET("/backup", c.BackupSales, authMiddleware, adminOnly)
	sales.DELETE("/range", c.DeleteSalesInRange, authMiddleware, adminOnly)
}
