package controller

import (
	"ims-intro/pkg/service"
	"ims-intro/pkg/service/dto"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type WebhookController struct {
	terminalService service.ITerminalService
	saleService     *service.SaleService
}

func NewWebhookController(terminalService service.ITerminalService, saleService *service.SaleService) *WebhookController {
	return &WebhookController{
		terminalService: terminalService,
		saleService:     saleService,
	}
}

// HandleSquareWebhook handles Square webhook events
func (wc *WebhookController) HandleSquareWebhook(c echo.Context) error {
	var webhookPayload dto.SquareWebhookPayload
	if err := c.Bind(&webhookPayload); err != nil {
		log.Errorf("Failed to bind webhook event: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid webhook payload"})
	}

	log.Infof("Received Square webhook event: type=%s", webhookPayload.Type)

	// Handle Terminal checkout updated event
	if webhookPayload.Type == "terminal.checkout.updated" {
		checkout := webhookPayload.Data.Object.Checkout
		if checkout == nil {
			log.Warn("Webhook event missing checkout data")
			return c.JSON(http.StatusOK, map[string]string{"status": "ignored"})
		}

		// Extract payment ID (if any)
		var paymentID *string
		if len(checkout.PaymentIDs) > 0 {
			paymentID = &checkout.PaymentIDs[0]
		}

		// Update checkout status in database
		err := wc.terminalService.UpdateCheckoutFromWebhook(
			checkout.ID,
			checkout.Status,
			paymentID,
		)
		if err != nil {
			log.Errorf("Failed to update checkout from webhook: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process webhook"})
		}

		// If checkout is completed, update sale payment status
		if checkout.Status == "COMPLETED" && checkout.ID != "" {
			log.Infof("Checkout %s completed, will update sale payment status via UpdateCheckoutFromWebhook", checkout.ID)
			// Note: UpdateCheckoutFromWebhook already handles updating the sale payment status
			// by calling the sale service internally through the terminal service
		}

		log.Infof("Successfully processed webhook for checkout %s", checkout.ID)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "processed"})
}

// RegisterWebhookRoutes registers webhook routes
func (wc *WebhookController) RegisterWebhookRoutes(e *echo.Echo) {
	webhooks := e.Group("/api/webhook")

	// Square webhook endpoint (no auth required as it's called by Square)
	webhooks.POST("/square/payment", wc.HandleSquareWebhook)

	log.Info("Webhook routes registered")
}
