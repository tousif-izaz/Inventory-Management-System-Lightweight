package controller

import (
	"ims-intro/pkg/service"
	"ims-intro/pkg/service/dto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type TerminalController struct {
	terminalService service.ITerminalService
}

func NewTerminalController(terminalService service.ITerminalService) *TerminalController {
	return &TerminalController{terminalService}
}

// RegisterTerminalRoutes registers all Terminal-related routes
func (c *TerminalController) RegisterTerminalRoutes(e *echo.Echo) {
	terminalGroup := e.Group("/api/terminal")

	// Device Management
	terminalGroup.POST("/devices", c.RegisterDevice)
	terminalGroup.GET("/devices", c.GetDevices)
	terminalGroup.GET("/devices/:id", c.GetDevice)
	terminalGroup.PUT("/devices/:id", c.UpdateDevice)
	terminalGroup.DELETE("/devices/:id", c.DeleteDevice)
	terminalGroup.POST("/devices/:id/default", c.SetDefaultDevice)

	// Checkout Management
	terminalGroup.POST("/checkout", c.CreateCheckout)
	terminalGroup.GET("/checkout/:id", c.GetCheckout)
	terminalGroup.POST("/checkout/:id/cancel", c.CancelCheckout)
	terminalGroup.GET("/checkout/:id/poll", c.PollCheckoutStatus)

	// Refund Management
	terminalGroup.POST("/refund", c.ProcessRefund)
	terminalGroup.GET("/refund/sale/:saleId", c.GetRefundForSale)

	// Sale-related checkout
	e.GET("/api/sales/:saleId/checkout", c.GetCheckoutBySaleID)
}

// RegisterDevice registers a new Terminal device
// POST /api/terminal/devices
func (c *TerminalController) RegisterDevice(ctx echo.Context) error {
	var deviceDTO dto.TerminalDeviceCreate
	if err := ctx.Bind(&deviceDTO); err != nil {
		log.Errorf("error binding device data: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// TODO: Get merchant ID from auth context/session
	// For now, using query parameter
	merchantID := ctx.QueryParam("merchant_id")
	if merchantID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Merchant ID required",
		})
	}

	device, err := c.terminalService.RegisterDevice(merchantID, deviceDTO)
	if err != nil {
		log.Errorf("error registering device: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, device)
}

// GetDevices retrieves all devices for a merchant
// GET /api/terminal/devices?merchant_id=XXX
func (c *TerminalController) GetDevices(ctx echo.Context) error {
	// TODO: Get merchant ID from auth context/session
	merchantID := ctx.QueryParam("merchant_id")
	if merchantID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Merchant ID required",
		})
	}

	devices, err := c.terminalService.GetDevices(merchantID)
	if err != nil {
		log.Errorf("error getting devices: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve devices",
		})
	}

	return ctx.JSON(http.StatusOK, devices)
}

// GetDevice retrieves a specific device
// GET /api/terminal/devices/:id
func (c *TerminalController) GetDevice(ctx echo.Context) error {
	deviceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid device ID",
		})
	}

	device, err := c.terminalService.GetDevice(deviceID)
	if err != nil {
		log.Errorf("error getting device: %v", err)
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "Device not found",
		})
	}

	return ctx.JSON(http.StatusOK, device)
}

// UpdateDevice updates a device
// PUT /api/terminal/devices/:id
func (c *TerminalController) UpdateDevice(ctx echo.Context) error {
	deviceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid device ID",
		})
	}

	var updateDTO dto.TerminalDeviceUpdate
	if err := ctx.Bind(&updateDTO); err != nil {
		log.Errorf("error binding update data: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := c.terminalService.UpdateDevice(deviceID, updateDTO); err != nil {
		log.Errorf("error updating device: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Device updated successfully",
	})
}

// DeleteDevice deletes a device
// DELETE /api/terminal/devices/:id
func (c *TerminalController) DeleteDevice(ctx echo.Context) error {
	deviceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid device ID",
		})
	}

	if err := c.terminalService.DeleteDevice(deviceID); err != nil {
		log.Errorf("error deleting device: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Device deleted successfully",
	})
}

// SetDefaultDevice sets a device as the default
// POST /api/terminal/devices/:id/default
func (c *TerminalController) SetDefaultDevice(ctx echo.Context) error {
	deviceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid device ID",
		})
	}

	// TODO: Get merchant ID from auth context
	merchantID := ctx.QueryParam("merchant_id")
	if merchantID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Merchant ID required",
		})
	}

	if err := c.terminalService.SetDefaultDevice(deviceID, merchantID); err != nil {
		log.Errorf("error setting default device: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Device set as default successfully",
	})
}

// CreateCheckout creates a new Terminal checkout
// POST /api/terminal/checkout
func (c *TerminalController) CreateCheckout(ctx echo.Context) error {
	var checkoutDTO dto.TerminalCheckoutCreate
	if err := ctx.Bind(&checkoutDTO); err != nil {
		log.Errorf("error binding checkout data: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	checkout, err := c.terminalService.CreateCheckout(checkoutDTO.SaleID, checkoutDTO.DeviceID)
	if err != nil {
		log.Errorf("error creating checkout: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Get device name for response
	device, _ := c.terminalService.GetDevice(checkout.DeviceID)
	deviceName := ""
	if device != nil {
		deviceName = device.DeviceName
	}

	response := dto.TerminalCheckoutResponse{
		CheckoutID:       checkout.CheckoutID,
		SaleID:           checkout.SaleID,
		SquareCheckoutID: checkout.SquareCheckoutID,
		Status:           checkout.Status,
		AmountMoney:      checkout.AmountMoney,
		Currency:         checkout.Currency,
		DeviceID:         checkout.DeviceID,
		DeviceName:       deviceName,
	}

	return ctx.JSON(http.StatusCreated, response)
}

// GetCheckout retrieves a checkout by ID
// GET /api/terminal/checkout/:id
func (c *TerminalController) GetCheckout(ctx echo.Context) error {
	checkoutID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid checkout ID",
		})
	}

	checkout, err := c.terminalService.GetCheckout(checkoutID)
	if err != nil {
		log.Errorf("error getting checkout: %v", err)
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "Checkout not found",
		})
	}

	return ctx.JSON(http.StatusOK, checkout)
}

// GetCheckoutBySaleID retrieves a checkout for a sale
// GET /api/sales/:saleId/checkout
func (c *TerminalController) GetCheckoutBySaleID(ctx echo.Context) error {
	saleID, err := strconv.ParseInt(ctx.Param("saleId"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid sale ID",
		})
	}

	checkout, err := c.terminalService.GetCheckoutBySaleID(saleID)
	if err != nil {
		log.Errorf("error getting checkout for sale: %v", err)
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "Checkout not found for this sale",
		})
	}

	return ctx.JSON(http.StatusOK, checkout)
}

// CancelCheckout cancels a Terminal checkout
// POST /api/terminal/checkout/:id/cancel
func (c *TerminalController) CancelCheckout(ctx echo.Context) error {
	checkoutID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid checkout ID",
		})
	}

	if err := c.terminalService.CancelCheckout(checkoutID); err != nil {
		log.Errorf("error canceling checkout: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Checkout canceled successfully",
	})
}

// PollCheckoutStatus polls Square for the latest checkout status
// GET /api/terminal/checkout/:id/poll
func (c *TerminalController) PollCheckoutStatus(ctx echo.Context) error {
	checkoutID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid checkout ID",
		})
	}

	checkout, err := c.terminalService.PollCheckoutStatus(checkoutID)
	if err != nil {
		log.Errorf("error polling checkout status: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, checkout)
}

// ===== Refund Management =====

// ProcessRefund handles refund requests
// POST /api/terminal/refund
func (c *TerminalController) ProcessRefund(ctx echo.Context) error {
	var refundDTO dto.RefundCreate
	if err := ctx.Bind(&refundDTO); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if refundDTO.SaleID == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Sale ID is required",
		})
	}

	log.Infof("Processing refund for sale %d", refundDTO.SaleID)

	refund, err := c.terminalService.RefundPayment(
		refundDTO.SaleID,
		refundDTO.AmountCents,
		refundDTO.Reason,
		refundDTO.RefundMessage,
	)
	if err != nil {
		log.Errorf("Failed to process refund: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, refund)
}

// GetRefundForSale retrieves refund for a sale
// GET /api/terminal/refund/sale/:saleId
func (c *TerminalController) GetRefundForSale(ctx echo.Context) error {
	saleIDStr := ctx.Param("saleId")
	saleID, err := strconv.ParseInt(saleIDStr, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid sale ID",
		})
	}

	refund, err := c.terminalService.GetRefundBySaleID(saleID)
	if err != nil {
		log.Errorf("Failed to get refund: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if refund == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "No refund found for this sale",
		})
	}

	return ctx.JSON(http.StatusOK, refund)
}
