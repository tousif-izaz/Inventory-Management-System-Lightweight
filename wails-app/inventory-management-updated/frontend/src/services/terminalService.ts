import { api } from './api';
import type {
  TerminalDevice,
  TerminalDeviceCreate,
  TerminalDeviceUpdate,
  TerminalCheckout,
  TerminalCheckoutCreate,
  TerminalCheckoutResponse,
  SquareRefund,
  RefundCreate,
} from '../types/terminal';

// Device Management APIs

export const registerTerminalDevice = async (
  merchantId: string,
  device: TerminalDeviceCreate
): Promise<TerminalDevice> => {
  return api.post<TerminalDevice>(`/terminal/devices?merchant_id=${merchantId}`, device);
};

export const getTerminalDevices = async (merchantId: string): Promise<TerminalDevice[]> => {
  return api.get<TerminalDevice[]>(`/terminal/devices?merchant_id=${merchantId}`);
};

export const getTerminalDevice = async (deviceId: number): Promise<TerminalDevice> => {
  return api.get<TerminalDevice>(`/terminal/devices/${deviceId}`);
};

export const updateTerminalDevice = async (
  deviceId: number,
  device: TerminalDeviceUpdate
): Promise<void> => {
  return api.put<void>(`/terminal/devices/${deviceId}`, device);
};

export const deleteTerminalDevice = async (deviceId: number): Promise<void> => {
  return api.delete<void>(`/terminal/devices/${deviceId}`);
};

export const setDefaultTerminalDevice = async (
  deviceId: number,
  merchantId: string
): Promise<void> => {
  return api.post<void>(`/terminal/devices/${deviceId}/default?merchant_id=${merchantId}`, {});
};

// Checkout Management APIs

export const createTerminalCheckout = async (
  checkout: TerminalCheckoutCreate
): Promise<TerminalCheckoutResponse> => {
  return api.post<TerminalCheckoutResponse>(`/terminal/checkout`, checkout);
};

export const getTerminalCheckout = async (checkoutId: number): Promise<TerminalCheckout> => {
  return api.get<TerminalCheckout>(`/terminal/checkout/${checkoutId}`);
};

export const getCheckoutBySaleId = async (saleId: number): Promise<TerminalCheckout> => {
  return api.get<TerminalCheckout>(`/terminal/checkout/sale/${saleId}`);
};

export const cancelTerminalCheckout = async (checkoutId: number): Promise<void> => {
  return api.post<void>(`/terminal/checkout/${checkoutId}/cancel`, {});
};

export const pollCheckoutStatus = async (checkoutId: number): Promise<TerminalCheckout> => {
  return api.get<TerminalCheckout>(`/terminal/checkout/${checkoutId}/poll`);
};

// Helper: Poll checkout status until completion
export const waitForCheckoutCompletion = async (
  checkoutId: number,
  onStatusUpdate?: (checkout: TerminalCheckout) => void,
  maxAttempts: number = 60, // 5 minutes with 5 second intervals
  intervalMs: number = 5000
): Promise<TerminalCheckout> => {
  let attempts = 0;

  return new Promise((resolve, reject) => {
    const poll = async () => {
      try {
        const checkout = await pollCheckoutStatus(checkoutId);

        if (onStatusUpdate) {
          onStatusUpdate(checkout);
        }

        // Check if checkout is in a terminal state
        if (
          checkout.status === 'COMPLETED' ||
          checkout.status === 'CANCELED' ||
          checkout.status === 'FAILED'
        ) {
          resolve(checkout);
          return;
        }

        // Continue polling if not done
        attempts++;
        if (attempts >= maxAttempts) {
          reject(new Error('Checkout polling timeout'));
          return;
        }

        setTimeout(poll, intervalMs);
      } catch (error) {
        reject(error);
      }
    };

    poll();
  });
};

// Refund Management APIs

export const processRefund = async (refund: RefundCreate): Promise<SquareRefund> => {
  return api.post<SquareRefund>(`/terminal/refund`, refund);
};

export const getRefundForSale = async (saleId: number): Promise<SquareRefund | null> => {
  try {
    return await api.get<SquareRefund>(`/terminal/refund/sale/${saleId}`);
  } catch (error) {
    return null;
  }
};
