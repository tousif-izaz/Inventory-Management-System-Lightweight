package controller

import (
	"ims-intro/pkg/service"
	"ims-intro/pkg/service/dto"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type SquareController struct {
	squareService service.ISquareService
}

func NewSquareController(squareService service.ISquareService) *SquareController {
	return &SquareController{squareService}
}

// RegisterSquareRoutes registers all Square-related routes
func (c *SquareController) RegisterSquareRoutes(e *echo.Echo) {
	squareGroup := e.Group("/square")

	// Config routes
	squareGroup.POST("/config", c.SaveConfig)
	squareGroup.GET("/config", c.GetConfig)

	// OAuth routes
	squareGroup.GET("/oauth/authorize", c.GetAuthorizationURL)
	squareGroup.GET("/oauth/callback", c.OAuthCallback)
	squareGroup.POST("/oauth/refresh/:merchantId", c.RefreshToken)

	// Token management routes
	squareGroup.GET("/tokens", c.GetAllTokens)
	squareGroup.GET("/token/status", c.GetTokenStatus)
	squareGroup.GET("/token/:merchantId", c.GetMerchantToken)
	squareGroup.DELETE("/token/:merchantId", c.DeleteMerchantToken)

	// Square API routes
	squareGroup.GET("/locations/:merchantId", c.GetMerchantLocations)
}

// SaveConfig handles saving Square application configuration
// POST /api/square/config
func (c *SquareController) SaveConfig(ctx echo.Context) error {
	var config dto.SquareConfigCreate
	if err := ctx.Bind(&config); err != nil {
		log.Errorf("error binding config: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error_message": "Invalid request body",
		})
	}

	if err := c.squareService.SaveConfig(config); err != nil {
		log.Errorf("error saving config: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Square configuration saved successfully",
	})
}

// GetConfig retrieves Square application configuration
// GET /api/square/config
func (c *SquareController) GetConfig(ctx echo.Context) error {
	config, err := c.squareService.GetConfig()
	if err != nil {
		log.Errorf("error getting config: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to retrieve configuration",
		})
	}

	// Don't expose the secret in the response
	safeConfig := map[string]interface{}{
		"application_id": config.ApplicationID,
		"environment":    config.Environment,
		"redirect_uri":   config.RedirectURI,
	}

	// Also try to get merchant ID if available (indicates OAuth completed)
	// We need to add this to the service interface
	// For now, leave it out

	return ctx.JSON(http.StatusOK, safeConfig)
}

// GetAuthorizationURL generates the Square OAuth authorization URL
// GET /api/square/oauth/authorize?scopes=PAYMENTS_READ,PAYMENTS_WRITE
func (c *SquareController) GetAuthorizationURL(ctx echo.Context) error {
	scopesParam := ctx.QueryParam("scopes")
	var scopes []string

	if scopesParam != "" {
		// Parse comma-separated scopes
		for _, scope := range ctx.QueryParams()["scopes"] {
			scopes = append(scopes, scope)
		}
	}

	authURL, err := c.squareService.GetAuthorizationURL(scopes)
	if err != nil {
		log.Errorf("error generating authorization URL: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to generate authorization URL",
		})
	}

	return ctx.JSON(http.StatusOK, dto.SquareAuthorizationURLResponse{
		AuthorizationURL: authURL,
		Scopes:           scopes,
	})
}

// OAuthCallback handles the OAuth callback from Square
// GET /api/square/oauth/callback?code=...&state=...
func (c *SquareController) OAuthCallback(ctx echo.Context) error {
	var callbackReq dto.SquareOAuthCallbackRequest
	if err := ctx.Bind(&callbackReq); err != nil {
		log.Errorf("error binding callback params: %v", err)
		return ctx.HTML(http.StatusBadRequest, `
			<html>
				<head><title>Authorization Failed</title></head>
				<body>
					<h1>Authorization Failed</h1>
					<p>Invalid callback parameters.</p>
					<script>setTimeout(() => window.close(), 3000);</script>
				</body>
			</html>
		`)
	}

	// Check for error from Square
	if callbackReq.Error != "" {
		log.Errorf("Square authorization error: %s", callbackReq.Error)
		return ctx.HTML(http.StatusBadRequest, `
			<html>
				<head><title>Authorization Failed</title></head>
				<body>
					<h1>Authorization Failed</h1>
					<p>Error: `+callbackReq.Error+`</p>
					<script>setTimeout(() => window.close(), 5000);</script>
				</body>
			</html>
		`)
	}

	// Exchange authorization code for tokens
	if err := c.squareService.ExchangeAuthorizationCode(callbackReq.Code); err != nil {
		log.Errorf("error exchanging authorization code: %v", err)
		return ctx.HTML(http.StatusInternalServerError, `
			<html>
				<head><title>Authorization Failed</title></head>
				<body>
					<h1>Authorization Failed</h1>
					<p>Failed to complete authorization. Please try again.</p>
					<script>setTimeout(() => window.close(), 5000);</script>
				</body>
			</html>
		`)
	}

	// Success - show success page
	return ctx.HTML(http.StatusOK, `
		<html>
			<head>
				<title>Authorization Successful</title>
				<style>
					body {
						font-family: Arial, sans-serif;
						text-align: center;
						padding: 50px;
						background-color: #f0f0f0;
					}
					.success-box {
						background-color: white;
						padding: 30px;
						border-radius: 10px;
						box-shadow: 0 2px 10px rgba(0,0,0,0.1);
						max-width: 500px;
						margin: 0 auto;
					}
					h1 { color: #4CAF50; }
					p { color: #666; }
				</style>
			</head>
			<body>
				<div class="success-box">
					<h1>✓ Authorization Successful!</h1>
					<p>Your Square account has been connected successfully.</p>
					<p>You can now close this window and return to the application.</p>
				</div>
				<script>
					// Auto-close after 3 seconds
					setTimeout(() => {
						window.close();
					}, 3000);
				</script>
			</body>
		</html>
	`)
}

// RefreshToken refreshes an expired access token
// POST /api/square/oauth/refresh/:merchantId
func (c *SquareController) RefreshToken(ctx echo.Context) error {
	merchantID := ctx.Param("merchantId")

	if err := c.squareService.RefreshAccessToken(merchantID); err != nil {
		log.Errorf("error refreshing token: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to refresh access token",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Access token refreshed successfully",
	})
}

// GetMerchantToken retrieves the OAuth token for a merchant
// GET /api/square/token/:merchantId
func (c *SquareController) GetMerchantToken(ctx echo.Context) error {
	merchantID := ctx.Param("merchantId")

	token, err := c.squareService.GetMerchantToken(merchantID)
	if err != nil {
		log.Errorf("error getting merchant token: %v", err)
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error_message": "Token not found",
		})
	}

	// Return sanitized token info (not the actual access token)
	response := map[string]interface{}{
		"merchant_id": token.MerchantID,
		"token_type":  token.TokenType,
		"scopes":      token.Scopes,
		"created_at":  token.CreatedAt,
		"updated_at":  token.UpdatedAt,
		"has_token":   token.AccessToken != "",
	}

	if token.ExpiresAt != nil {
		response["expires_at"] = token.ExpiresAt
	}

	return ctx.JSON(http.StatusOK, response)
}

// DeleteMerchantToken deletes the OAuth token for a merchant
// DELETE /api/square/token/:merchantId
func (c *SquareController) DeleteMerchantToken(ctx echo.Context) error {
	merchantID := ctx.Param("merchantId")

	if err := c.squareService.DeleteMerchantToken(merchantID); err != nil {
		log.Errorf("error deleting merchant token: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to delete token",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Token deleted successfully",
	})
}

// GetAllTokens retrieves all merchant IDs that have OAuth tokens
// GET /api/square/tokens
func (c *SquareController) GetAllTokens(ctx echo.Context) error {
	merchantIDs, err := c.squareService.GetAllMerchantIDs()
	if err != nil {
		log.Errorf("error getting merchant IDs: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to retrieve merchant IDs",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"merchants": merchantIDs,
	})
}

// GetTokenStatus checks if any OAuth token exists and returns basic status
// GET /api/square/token/status
func (c *SquareController) GetTokenStatus(ctx echo.Context) error {
	// Try to get the saved merchant ID from config
	// Note: This requires adding GetMerchantID to the service interface
	// For now, we'll return a helpful error message

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Use GET /api/square/config to check configuration status",
	})
}

// GetMerchantLocations retrieves all locations for a merchant
// GET /api/square/locations/:merchantId
func (c *SquareController) GetMerchantLocations(ctx echo.Context) error {
	merchantID := ctx.Param("merchantId")

	locations, err := c.squareService.GetMerchantLocations(merchantID)
	if err != nil {
		log.Errorf("error getting merchant locations: %v", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error_message": "Failed to retrieve locations",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"locations": locations,
	})
}
