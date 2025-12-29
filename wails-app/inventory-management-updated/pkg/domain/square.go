package domain

import "time"

// SquareOAuthToken represents the OAuth token for a Square merchant
type SquareOAuthToken struct {
	TokenID      int64     `json:"token_id"`
	MerchantID   string    `json:"merchant_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	TokenType    string    `json:"token_type"`
	Scopes       string    `json:"scopes"` // Comma-separated list of scopes
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SquareConfig represents Square application configuration
// These will be stored in the Settings table with keys prefixed with "square_"
type SquareConfig struct {
	ApplicationID     string `json:"application_id"`
	ApplicationSecret string `json:"application_secret"`
	Environment       string `json:"environment"` // "sandbox" or "production"
	RedirectURI       string `json:"redirect_uri"`
}
