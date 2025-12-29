import { useState, useEffect } from 'react';
import { Button } from '../UI/Button';
import { Input } from '../UI/Input';
import { useToast } from '../UI/Toast';
import { squareService } from '../../services/squareService';
import { SquareOAuthStatus } from '../../types/square';
import {
  CheckCircleIcon,
  XCircleIcon,
  ArrowPathIcon,
  TrashIcon,
} from '@heroicons/react/24/outline';

export const SquareSetup = () => {
  const { showToast } = useToast();
  const [status, setStatus] = useState<SquareOAuthStatus>({
    isConfigured: false,
    isAuthorized: false,
  });
  const [loading, setLoading] = useState(true);
  const [configMode, setConfigMode] = useState(false);

  // Configuration form state
  const [applicationId, setApplicationId] = useState('');
  const [applicationSecret, setApplicationSecret] = useState('');
  const [environment, setEnvironment] = useState<'production' | 'sandbox'>('production');
  const [redirectUri, setRedirectUri] = useState('http://localhost:8080/square/oauth/callback');
  const [merchantId, setMerchantId] = useState('');
  const [locations, setLocations] = useState<any[]>([]);

  useEffect(() => {
    loadStatus();
  }, []);

  const loadStatus = async () => {
    try {
      setLoading(true);
      const oauthStatus = await squareService.getOAuthStatus(merchantId || undefined);
      setStatus(oauthStatus);

      if (oauthStatus.config) {
        setApplicationId(oauthStatus.config.application_id);
        setEnvironment(oauthStatus.config.environment);
        setRedirectUri(oauthStatus.config.redirect_uri || 'http://localhost:8080/square/oauth/callback');
      }

      if (oauthStatus.merchantId) {
        setMerchantId(oauthStatus.merchantId);
        // Load locations if authorized
        if (oauthStatus.isAuthorized) {
          await loadLocations(oauthStatus.merchantId);
        }
      }
    } catch (error) {
      console.error('Failed to load Square status:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSaveConfig = async () => {
    if (!applicationId.trim()) {
      showToast('error', 'Application ID is required');
      return;
    }

    // If updating existing config and secret is empty, user wants to keep existing secret
    if (status.isConfigured && !applicationSecret.trim()) {
      showToast('error', 'Please enter your Application Secret to update the configuration');
      return;
    }

    // If creating new config, secret is required
    if (!status.isConfigured && !applicationSecret.trim()) {
      showToast('error', 'Application Secret is required');
      return;
    }

    try {
      await squareService.saveConfig({
        application_id: applicationId,
        application_secret: applicationSecret,
        environment,
        redirect_uri: redirectUri.trim() || undefined,
      });

      showToast('success', 'Square configuration saved successfully');
      setConfigMode(false);
      setApplicationSecret(''); // Clear the secret from UI
      await loadStatus();
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to save configuration';
      showToast('error', `Failed to save configuration: ${errorMessage}`);
    }
  };

  const loadLocations = async (merchantId: string) => {
    try {
      const response = await fetch(`http://localhost:37285/api/square/locations/${merchantId}`);
      if (response.ok) {
        const data = await response.json();
        setLocations(data.locations || []);
      }
    } catch (error) {
      console.error('Failed to load locations:', error);
    }
  };

  const handleInitiateOAuth = async () => {
    try {
      await squareService.initiateOAuth([
        'PAYMENTS_READ',
        'PAYMENTS_WRITE',
        'DEVICE_CREDENTIAL_MANAGEMENT',
      ]);

      showToast('info', 'Authorization window opened. Please approve in your browser.');

      // Poll for authorization status
      const pollInterval = setInterval(async () => {
        const newStatus = await squareService.getOAuthStatus(merchantId || undefined);
        if (newStatus.isAuthorized) {
          clearInterval(pollInterval);
          setStatus(newStatus);
          if (newStatus.merchantId) {
            setMerchantId(newStatus.merchantId);
            // Load locations after successful authorization
            await loadLocations(newStatus.merchantId);
          }
          showToast('success', 'Square account connected successfully!');
        }
      }, 3000); // Poll every 3 seconds

      // Stop polling after 2 minutes
      setTimeout(() => clearInterval(pollInterval), 120000);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to initiate OAuth';
      showToast('error', `Failed to initiate OAuth: ${errorMessage}`);
    }
  };

  const handleRefreshToken = async () => {
    if (!merchantId) {
      showToast('error', 'Merchant ID is required to refresh token');
      return;
    }

    try {
      await squareService.refreshToken(merchantId);
      showToast('success', 'Access token refreshed successfully');
      await loadStatus();
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to refresh token';
      showToast('error', `Failed to refresh token: ${errorMessage}`);
    }
  };

  const handleDisconnect = async () => {
    if (!merchantId) {
      showToast('error', 'Merchant ID is required to disconnect');
      return;
    }

    if (!confirm('Are you sure you want to disconnect your Square account?')) {
      return;
    }

    try {
      await squareService.deleteMerchantToken(merchantId);
      showToast('success', 'Square account disconnected successfully');
      await loadStatus();
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to disconnect';
      showToast('error', `Failed to disconnect: ${errorMessage}`);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="bg-white shadow rounded-lg p-6">
        <h2 className="text-2xl font-bold mb-4">Square Payment Integration</h2>
        <p className="text-gray-600 mb-6">
          Connect your Square account to process payments and manage Terminal devices.
        </p>

        {/* Status Overview */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <div className="border rounded-lg p-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-gray-700">Configuration</span>
              {status.isConfigured ? (
                <CheckCircleIcon className="h-6 w-6 text-green-500" />
              ) : (
                <XCircleIcon className="h-6 w-6 text-red-500" />
              )}
            </div>
            <p className="text-xs text-gray-500 mt-1">
              {status.isConfigured ? 'Configured' : 'Not configured'}
            </p>
          </div>

          <div className="border rounded-lg p-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-gray-700">Authorization</span>
              {status.isAuthorized ? (
                <CheckCircleIcon className="h-6 w-6 text-green-500" />
              ) : (
                <XCircleIcon className="h-6 w-6 text-red-500" />
              )}
            </div>
            <p className="text-xs text-gray-500 mt-1">
              {status.isAuthorized ? 'Connected' : 'Not connected'}
            </p>
          </div>
        </div>

        {/* Configuration Section */}
        {!status.isConfigured || configMode ? (
          <div className="border-t pt-6">
            <h3 className="text-lg font-semibold mb-4">Square Configuration</h3>
            <div className="space-y-4">
              <Input
                label="Application ID"
                value={applicationId}
                onChange={(e) => setApplicationId(e.target.value)}
                placeholder="sq0idp-..."
                required
              />

              <Input
                label="Application Secret"
                type="password"
                value={applicationSecret}
                onChange={(e) => setApplicationSecret(e.target.value)}
                placeholder="sq0csp-..."
                required
                helperText={status.isConfigured ? "Re-enter your secret to update configuration" : undefined}
              />

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Environment
                </label>
                <select
                  value={environment}
                  onChange={(e) => setEnvironment(e.target.value as 'production' | 'sandbox')}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="production">Production</option>
                  <option value="sandbox">Sandbox</option>
                </select>
              </div>

              <Input
                label="OAuth Redirect URI"
                value={redirectUri}
                onChange={(e) => setRedirectUri(e.target.value)}
                placeholder="http://localhost:8080/square/oauth/callback"
                helperText="The callback URL configured in your Square Developer Dashboard"
              />

              <div className="flex gap-2">
                <Button onClick={handleSaveConfig}>
                  Save Configuration
                </Button>
                {status.isConfigured && (
                  <Button variant="secondary" onClick={() => setConfigMode(false)}>
                    Cancel
                  </Button>
                )}
              </div>
            </div>
          </div>
        ) : (
          <div className="border-t pt-6">
            <div className="flex justify-between items-center mb-4">
              <div>
                <h3 className="text-lg font-semibold">Current Configuration</h3>
                <p className="text-sm text-gray-600">
                  Application ID: {status.config?.application_id}
                </p>
                <p className="text-sm text-gray-600">
                  Environment: {status.config?.environment}
                </p>
                <p className="text-sm text-gray-600">
                  Redirect URI: {status.config?.redirect_uri}
                </p>
              </div>
              <Button variant="secondary" onClick={() => setConfigMode(true)}>
                Update Configuration
              </Button>
            </div>
          </div>
        )}

        {/* Authorization Section */}
        {status.isConfigured && (
          <div className="border-t pt-6 mt-6">
            <h3 className="text-lg font-semibold mb-4">Square Account Connection</h3>

            {!status.isAuthorized ? (
              <div>
                <p className="text-sm text-gray-600 mb-4">
                  Connect your Square account to enable payment processing and Terminal management.
                </p>
                <Button onClick={handleInitiateOAuth}>
                  Connect Square Account
                </Button>
              </div>
            ) : (
              <div className="space-y-4">
                <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                  <div className="flex items-start">
                    <CheckCircleIcon className="h-5 w-5 text-green-500 mt-0.5 mr-3" />
                    <div className="flex-1">
                      <h4 className="text-sm font-medium text-green-800">Account Connected</h4>
                      <p className="text-sm text-green-700 mt-1">
                        Merchant ID: {status.merchantId}
                      </p>
                      {status.token?.expires_at && (
                        <p className="text-xs text-green-600 mt-1">
                          Token expires: {new Date(status.token.expires_at).toLocaleString()}
                        </p>
                      )}
                    </div>
                  </div>
                </div>

                <div className="flex gap-2">
                  <Button
                    variant="secondary"
                    onClick={handleRefreshToken}
                  >
                    <ArrowPathIcon className="h-4 w-4 mr-2 inline-block" />
                    Refresh Token
                  </Button>
                  <Button
                    variant="danger"
                    onClick={handleDisconnect}
                  >
                    <TrashIcon className="h-4 w-4 mr-2 inline-block" />
                    Disconnect
                  </Button>
                </div>

                {/* Display Locations */}
                {locations.length > 0 && (
                  <div className="border-t pt-4 mt-4">
                    <h4 className="text-sm font-medium text-gray-900 mb-3">Business Locations</h4>
                    <div className="space-y-2">
                      {locations.map((location) => (
                        <div key={location.id} className="border rounded-lg p-3 bg-gray-50">
                          <div className="flex items-start justify-between">
                            <div>
                              <p className="text-sm font-medium text-gray-900">{location.name}</p>
                              {location.address && (
                                <p className="text-xs text-gray-600 mt-1">{location.address}</p>
                              )}
                              {location.phone_number && (
                                <p className="text-xs text-gray-600 mt-1">{location.phone_number}</p>
                              )}
                            </div>
                            <span className={`text-xs px-2 py-1 rounded ${
                              location.status === 'ACTIVE' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                            }`}>
                              {location.status}
                            </span>
                          </div>
                          <p className="text-xs text-gray-500 mt-2">Location ID: {location.id}</p>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        )}

        {/* Merchant ID Input for checking status */}
        {status.isConfigured && !status.isAuthorized && (
          <div className="border-t pt-6 mt-6">
            <h3 className="text-lg font-semibold mb-4">Check Authorization Status</h3>
            <div className="flex gap-2">
              <Input
                label="Merchant ID (optional)"
                value={merchantId}
                onChange={(e) => setMerchantId(e.target.value)}
                placeholder="Enter merchant ID to check status"
              />
              <Button variant="secondary" onClick={loadStatus} className="mt-6">
                Check Status
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Information Panel */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <h4 className="text-sm font-medium text-blue-800 mb-2">Setup Instructions</h4>
        <ol className="list-decimal list-inside space-y-1 text-sm text-blue-700">
          <li>Enter your Square Application ID and Secret from the Square Developer Dashboard</li>
          <li>Select the environment (Production or Sandbox)</li>
          <li>Save the configuration</li>
          <li>Click "Connect Square Account" to authorize the application</li>
          <li>Approve the permissions in the browser window that opens</li>
          <li>The window will close automatically when authorization is complete</li>
        </ol>
      </div>
    </div>
  );
};
