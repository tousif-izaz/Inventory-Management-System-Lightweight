package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service/dto"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

const (
	SquareProductionOAuthURL = "https://connect.squareup.com/oauth2"
	SquareSandboxOAuthURL    = "https://connect.squareupsandbox.com/oauth2"
	SquareProductionAPIURL   = "https://connect.squareup.com/v2"
	SquareSandboxAPIURL      = "https://connect.squareupsandbox.com/v2"
)

type ISquareService interface {
	// Config management
	SaveConfig(config dto.SquareConfigCreate) error
	GetConfig() (domain.SquareConfig, error)

	// OAuth flow
	GetAuthorizationURL(scopes []string) (string, error)
	ExchangeAuthorizationCode(code string) error
	RefreshAccessToken(merchantID string) error

	// Token management
	GetMerchantToken(merchantID string) (domain.SquareOAuthToken, error)
	DeleteMerchantToken(merchantID string) error
	GetAllMerchantIDs() ([]string, error)

	// Square API operations
	GetMerchantLocations(merchantID string) ([]dto.SquareLocation, error)
}

type SquareService struct {
	repository repository.ISquareRepository
}

func NewSquareService(repository repository.ISquareRepository) ISquareService {
	return &SquareService{repository}
}

// SaveConfig saves Square application configuration
func (s *SquareService) SaveConfig(config dto.SquareConfigCreate) error {
	// Validate environment
	if config.Environment != "sandbox" && config.Environment != "production" {
		return errors.New("environment must be 'sandbox' or 'production'")
	}

	// Use provided redirect URI or set default
	redirectURI := config.RedirectURI
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/square/oauth/callback"
	}

	domainConfig := domain.SquareConfig{
		ApplicationID:     config.ApplicationID,
		ApplicationSecret: config.ApplicationSecret,
		Environment:       config.Environment,
		RedirectURI:       redirectURI,
	}

	return s.repository.SaveConfig(domainConfig)
}

// GetConfig retrieves Square application configuration
func (s *SquareService) GetConfig() (domain.SquareConfig, error) {
	return s.repository.GetConfig()
}

// GetAuthorizationURL generates the Square OAuth authorization URL
func (s *SquareService) GetAuthorizationURL(scopes []string) (string, error) {
	config, err := s.repository.GetConfig()
	if err != nil {
		log.Errorf("error getting Square config: %v", err)
		return "", err
	}

	// Default scopes for payment processing with Terminal API
	if len(scopes) == 0 {
		scopes = []string{
			"PAYMENTS_READ",
			"PAYMENTS_WRITE",
			"DEVICE_CREDENTIAL_MANAGEMENT",
		}
	}

	baseURL := SquareProductionOAuthURL
	if config.Environment == "sandbox" {
		baseURL = SquareSandboxOAuthURL
	}

	// Build authorization URL
	params := url.Values{}
	params.Add("client_id", config.ApplicationID)
	params.Add("scope", strings.Join(scopes, " "))
	params.Add("session", "false") // Set to false for long-lived tokens
	params.Add("redirect_uri", config.RedirectURI) // Required by Square OAuth

	authURL := fmt.Sprintf("%s/authorize?%s", baseURL, params.Encode())

	log.Info(fmt.Sprintf("Generated Square authorization URL: %s", authURL))
	return authURL, nil
}

// ExchangeAuthorizationCode exchanges the authorization code for access tokens
func (s *SquareService) ExchangeAuthorizationCode(code string) error {
	config, err := s.repository.GetConfig()
	if err != nil {
		log.Errorf("error getting Square config: %v", err)
		return err
	}

	baseURL := SquareProductionOAuthURL
	if config.Environment == "sandbox" {
		baseURL = SquareSandboxOAuthURL
	}

	// Prepare token request
	tokenRequest := dto.SquareOAuthTokenRequest{
		ClientID:     config.ApplicationID,
		ClientSecret: config.ApplicationSecret,
		Code:         code,
		GrantType:    "authorization_code",
		RedirectURI:  config.RedirectURI, // Must match the redirect_uri used in authorization request
	}

	jsonData, err := json.Marshal(tokenRequest)
	if err != nil {
		log.Errorf("error marshaling token request: %v", err)
		return err
	}

	// Make request to Square
	tokenURL := fmt.Sprintf("%s/token", baseURL)
	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Errorf("error making token request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		log.Errorf("Square token request failed (status %d): %v", resp.StatusCode, errorResponse)
		log.Errorf("Request details - URL: %s, Environment: %s", tokenURL, config.Environment)
		return fmt.Errorf("Square token request failed with status %d: %v", resp.StatusCode, errorResponse)
	}

	// Parse response
	var tokenResponse dto.SquareOAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		log.Errorf("error decoding token response: %v", err)
		return err
	}

	// Parse expiration time if present
	var expiresAt *time.Time
	if tokenResponse.ExpiresAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, tokenResponse.ExpiresAt)
		if err == nil {
			expiresAt = &parsedTime
		}
	}

	// Save token to database
	token := domain.SquareOAuthToken{
		MerchantID:   tokenResponse.MerchantID,
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    tokenResponse.TokenType,
		Scopes:       "", // We'll store the scopes from the request
	}

	// Check if token already exists for this merchant
	existingToken, err := s.repository.GetOAuthTokenByMerchant(tokenResponse.MerchantID)
	if err == nil && existingToken.MerchantID != "" {
		// Update existing token
		return s.repository.UpdateOAuthToken(tokenResponse.MerchantID, token)
	}

	// Save new token
	if err := s.repository.SaveOAuthToken(token); err != nil {
		return err
	}

	// Also save merchant ID to settings for easy retrieval
	if err := s.repository.SaveMerchantID(tokenResponse.MerchantID); err != nil {
		log.Errorf("failed to save merchant ID to settings: %v", err)
		// Don't fail the whole operation if this fails
	}

	return nil
}

// RefreshAccessToken refreshes an expired access token
func (s *SquareService) RefreshAccessToken(merchantID string) error {
	config, err := s.repository.GetConfig()
	if err != nil {
		log.Errorf("error getting Square config: %v", err)
		return err
	}

	token, err := s.repository.GetOAuthTokenByMerchant(merchantID)
	if err != nil {
		log.Errorf("error getting merchant token: %v", err)
		return err
	}

	if token.RefreshToken == "" {
		return errors.New("no refresh token available")
	}

	baseURL := SquareProductionOAuthURL
	if config.Environment == "sandbox" {
		baseURL = SquareSandboxOAuthURL
	}

	// Prepare refresh request
	refreshRequest := dto.SquareOAuthRefreshRequest{
		ClientID:     config.ApplicationID,
		ClientSecret: config.ApplicationSecret,
		GrantType:    "refresh_token",
		RefreshToken: token.RefreshToken,
	}

	jsonData, err := json.Marshal(refreshRequest)
	if err != nil {
		log.Errorf("error marshaling refresh request: %v", err)
		return err
	}

	// Make request to Square
	tokenURL := fmt.Sprintf("%s/token", baseURL)
	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Errorf("error making refresh request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		log.Errorf("Square refresh request failed: %v", errorResponse)
		return fmt.Errorf("Square refresh request failed with status %d", resp.StatusCode)
	}

	// Parse response
	var tokenResponse dto.SquareOAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		log.Errorf("error decoding refresh response: %v", err)
		return err
	}

	// Parse expiration time if present
	var expiresAt *time.Time
	if tokenResponse.ExpiresAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, tokenResponse.ExpiresAt)
		if err == nil {
			expiresAt = &parsedTime
		}
	}

	// Update token in database
	updatedToken := domain.SquareOAuthToken{
		MerchantID:   tokenResponse.MerchantID,
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    tokenResponse.TokenType,
		Scopes:       token.Scopes, // Keep existing scopes
	}

	return s.repository.UpdateOAuthToken(merchantID, updatedToken)
}

// GetMerchantToken retrieves the OAuth token for a merchant
func (s *SquareService) GetMerchantToken(merchantID string) (domain.SquareOAuthToken, error) {
	return s.repository.GetOAuthTokenByMerchant(merchantID)
}

// DeleteMerchantToken deletes the OAuth token for a merchant
func (s *SquareService) DeleteMerchantToken(merchantID string) error {
	return s.repository.DeleteOAuthToken(merchantID)
}

// GetAllMerchantIDs retrieves all merchant IDs that have OAuth tokens
func (s *SquareService) GetAllMerchantIDs() ([]string, error) {
	return s.repository.GetAllMerchantIDs()
}

// GetMerchantLocations retrieves all locations for a merchant
func (s *SquareService) GetMerchantLocations(merchantID string) ([]dto.SquareLocation, error) {
	// Get the access token for this merchant
	token, err := s.repository.GetOAuthTokenByMerchant(merchantID)
	if err != nil {
		log.Errorf("error getting merchant token: %v", err)
		return nil, err
	}

	// Get config to determine environment
	config, err := s.repository.GetConfig()
	if err != nil {
		log.Errorf("error getting Square config: %v", err)
		return nil, err
	}

	baseURL := SquareProductionAPIURL
	if config.Environment == "sandbox" {
		baseURL = SquareSandboxAPIURL
	}

	// Make request to Square Locations API
	locationsURL := fmt.Sprintf("%s/locations", baseURL)
	req, err := http.NewRequest("GET", locationsURL, nil)
	if err != nil {
		log.Errorf("error creating locations request: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Square-Version", "2024-01-18") // Use a recent API version

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error making locations request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		log.Errorf("Square locations request failed (status %d): %v", resp.StatusCode, errorResponse)
		return nil, fmt.Errorf("Square locations request failed with status %d: %v", resp.StatusCode, errorResponse)
	}

	// Parse response
	var locationsResponse struct {
		Locations []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Address     *struct {
				AddressLine1 string `json:"address_line_1"`
				Locality     string `json:"locality"`
				AdminDistrict string `json:"administrative_district_level_1"`
				PostalCode   string `json:"postal_code"`
			} `json:"address"`
			PhoneNumber string `json:"phone_number"`
			Timezone    string `json:"timezone"`
			Status      string `json:"status"`
		} `json:"locations"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&locationsResponse); err != nil {
		log.Errorf("error decoding locations response: %v", err)
		return nil, err
	}

	// Convert to DTO
	locations := make([]dto.SquareLocation, 0, len(locationsResponse.Locations))
	for _, loc := range locationsResponse.Locations {
		address := ""
		if loc.Address != nil {
			address = fmt.Sprintf("%s, %s, %s %s",
				loc.Address.AddressLine1,
				loc.Address.Locality,
				loc.Address.AdminDistrict,
				loc.Address.PostalCode,
			)
		}

		locations = append(locations, dto.SquareLocation{
			ID:          loc.ID,
			Name:        loc.Name,
			Address:     address,
			PhoneNumber: loc.PhoneNumber,
			Timezone:    loc.Timezone,
			Status:      loc.Status,
		})
	}

	log.Info(fmt.Sprintf("Retrieved %d locations for merchant %s", len(locations), merchantID))
	return locations, nil
}
