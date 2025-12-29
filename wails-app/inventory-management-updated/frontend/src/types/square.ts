// Square OAuth and configuration types

export interface SquareConfig {
  application_id: string;
  environment: 'production' | 'sandbox';
  redirect_uri: string;
}

export interface SquareConfigCreate {
  application_id: string;
  application_secret: string;
  environment: 'production' | 'sandbox';
  redirect_uri?: string;
}

export interface SquareOAuthToken {
  merchant_id: string;
  token_type: string;
  scopes: string;
  created_at: string;
  updated_at: string;
  has_token: boolean;
  expires_at?: string;
}

export interface SquareAuthorizationURLResponse {
  authorization_url: string;
  scopes: string[];
  state?: string;
}

export interface SquareOAuthStatus {
  isConfigured: boolean;
  isAuthorized: boolean;
  merchantId?: string;
  config?: SquareConfig;
  token?: SquareOAuthToken;
}
