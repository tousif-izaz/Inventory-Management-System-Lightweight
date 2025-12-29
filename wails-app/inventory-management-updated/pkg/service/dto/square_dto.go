package dto

// SquareOAuthTokenRequest represents the request to exchange authorization code for tokens
type SquareOAuthTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"` // Should be "authorization_code"
	RedirectURI  string `json:"redirect_uri,omitempty"`
}

// SquareOAuthTokenResponse represents the response from Square OAuth token endpoint
type SquareOAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	MerchantID   string `json:"merchant_id"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ShortLived   bool   `json:"short_lived,omitempty"`
}

// SquareOAuthRefreshRequest represents the request to refresh an access token
type SquareOAuthRefreshRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"` // Should be "refresh_token"
	RefreshToken string `json:"refresh_token"`
}

// SquareConfigCreate represents the request to save Square configuration
type SquareConfigCreate struct {
	ApplicationID     string `json:"application_id"`
	ApplicationSecret string `json:"application_secret"`
	Environment       string `json:"environment"`   // "sandbox" or "production"
	RedirectURI       string `json:"redirect_uri,omitempty"` // OAuth callback URL
}

// SquareAuthorizationURLResponse represents the response containing the OAuth authorization URL
type SquareAuthorizationURLResponse struct {
	AuthorizationURL string   `json:"authorization_url"`
	State            string   `json:"state,omitempty"`
	Scopes           []string `json:"scopes"`
}

// SquareOAuthCallbackRequest represents the callback parameters from Square
type SquareOAuthCallbackRequest struct {
	Code  string `query:"code"`
	State string `query:"state"`
	Error string `query:"error"`
}

// SquareLocation represents a Square business location
type SquareLocation struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
	Status      string `json:"status,omitempty"`
}
