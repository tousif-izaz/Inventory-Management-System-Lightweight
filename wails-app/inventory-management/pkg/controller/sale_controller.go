package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"inventory-management/pkg/controller/request"
	"inventory-management/pkg/controller/response"
	"inventory-management/pkg/domain"
	"inventory-management/pkg/service"
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

// RegisterSaleRoutes registers all sale routes
func (c *SaleController) RegisterSaleRoutes(e *echo.Echo) {
	sales := e.Group("/sales")

	sales.POST("", c.CreateSale)
	sales.GET("", c.GetAllSales)
	sales.GET("/:id", c.GetSaleByID)
	sales.GET("/:id/receipt", c.GenerateReceipt)
	sales.PUT("/:id/payment-status", c.UpdatePaymentStatus)
	sales.POST("/:id/refund", c.ProcessRefund)
}
