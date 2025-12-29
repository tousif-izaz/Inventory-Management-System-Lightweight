import { api } from './api';
import {
  SquareConfig,
  SquareConfigCreate,
  SquareOAuthToken,
  SquareAuthorizationURLResponse,
  SquareOAuthStatus,
} from '../types/square';

class SquareService {
  // Configuration Management
  async saveConfig(config: SquareConfigCreate): Promise<{ message: string }> {
    return api.post<{ message: string }>('/square/config', config);
  }

  async getConfig(): Promise<SquareConfig> {
    return api.get<SquareConfig>('/square/config');
  }

  // OAuth Flow
  async getAuthorizationURL(scopes?: string[]): Promise<SquareAuthorizationURLResponse> {
    const scopeQuery = scopes && scopes.length > 0 ? `?scopes=${scopes.join(',')}` : '';
    return api.get<SquareAuthorizationURLResponse>(`/square/oauth/authorize${scopeQuery}`);
  }

  async initiateOAuth(scopes?: string[]): Promise<void> {
    try {
      const response = await this.getAuthorizationURL(scopes);
      // Open authorization URL in default browser (external to Wails app)
      window.open(response.authorization_url, '_blank');
    } catch (error) {
      console.error('Error initiating OAuth:', error);
      throw error;
    }
  }

  async refreshToken(merchantId: string): Promise<{ message: string }> {
    return api.post<{ message: string }>(`/square/oauth/refresh/${merchantId}`, {});
  }

  // Token Management
  async getMerchantToken(merchantId: string): Promise<SquareOAuthToken> {
    return api.get<SquareOAuthToken>(`/square/token/${merchantId}`);
  }

  async deleteMerchantToken(merchantId: string): Promise<{ message: string }> {
    return api.delete<{ message: string }>(`/square/token/${merchantId}`);
  }

  // Status Check
  async getOAuthStatus(merchantId?: string): Promise<SquareOAuthStatus> {
    try {
      const config = await this.getConfig();

      let token: SquareOAuthToken | undefined;
      let isAuthorized = false;
      let actualMerchantId: string | undefined = merchantId;

      // If merchantId not provided, try to get list of all merchant IDs from backend
      if (!merchantId) {
        try {
          // Call a new endpoint that lists all tokens/merchants
          const tokensResponse = await api.get<any>('/square/tokens');
          if (tokensResponse && tokensResponse.merchants && tokensResponse.merchants.length > 0) {
            // Use the first merchant ID
            actualMerchantId = tokensResponse.merchants[0];
          }
        } catch (error) {
          // No tokens found yet
          isAuthorized = false;
        }
      }

      // Fetch token if we have a merchant ID
      if (actualMerchantId) {
        try {
          const merchantToken = await this.getMerchantToken(actualMerchantId);
          token = merchantToken;
          isAuthorized = merchantToken.has_token;
        } catch (error) {
          // Token not found - not authorized yet
          isAuthorized = false;
        }
      }

      return {
        isConfigured: !!config.application_id,
        isAuthorized,
        merchantId: actualMerchantId,
        config,
        token,
      };
    } catch (error) {
      // Configuration not found
      return {
        isConfigured: false,
        isAuthorized: false,
      };
    }
  }
}

export const squareService = new SquareService();
