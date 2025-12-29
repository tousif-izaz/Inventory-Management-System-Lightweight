// Terminal Device Types
export interface TerminalDevice {
  device_id: number;
  merchant_id: string;
  square_device_id: string;
  device_name: string;
  location_id?: string;
  device_code?: string;
  status: 'active' | 'inactive' | 'paired' | 'unpaired';
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

export interface TerminalDeviceCreate {
  square_device_id: string;
  device_name: string;
  location_id?: string;
}

export interface TerminalDeviceUpdate {
  device_name?: string;
  location_id?: string;
  status?: 'active' | 'inactive';
}

// Terminal Checkout Types
export interface TerminalCheckout {
  checkout_id: number;
  sale_id: number;
  device_id: number;
  device_name?: string;
  merchant_id: string;
  square_checkout_id: string;
  amount_money: number;
  currency: string;
  status: 'PENDING' | 'IN_PROGRESS' | 'CANCEL_REQUESTED' | 'CANCELED' | 'COMPLETED' | 'FAILED';
  square_payment_id?: string;
  error_code?: string;
  error_message?: string;
  created_at: string;
  updated_at: string;
  completed_at?: string;
}

export interface TerminalCheckoutCreate {
  sale_id: number;
  device_id?: number;
}

export interface TerminalCheckoutResponse {
  checkout: TerminalCheckout;
  message?: string;
}

// Payment Method Type
export type PaymentMethod = 'cash' | 'credit_card' | 'debit_card' | 'square_terminal' | 'other';

// Checkout Status Display
export interface CheckoutStatusInfo {
  status: TerminalCheckout['status'];
  label: string;
  color: 'blue' | 'yellow' | 'red' | 'green' | 'gray';
  icon: string;
  description: string;
}

export const CHECKOUT_STATUS_INFO: Record<TerminalCheckout['status'], CheckoutStatusInfo> = {
  PENDING: {
    status: 'PENDING',
    label: 'Pending',
    color: 'blue',
    icon: 'clock',
    description: 'Waiting to be sent to terminal',
  },
  IN_PROGRESS: {
    status: 'IN_PROGRESS',
    label: 'In Progress',
    color: 'yellow',
    icon: 'loading',
    description: 'Customer is completing payment on terminal',
  },
  CANCEL_REQUESTED: {
    status: 'CANCEL_REQUESTED',
    label: 'Canceling',
    color: 'yellow',
    icon: 'loading',
    description: 'Cancellation requested',
  },
  CANCELED: {
    status: 'CANCELED',
    label: 'Canceled',
    color: 'gray',
    icon: 'close',
    description: 'Payment was canceled',
  },
  COMPLETED: {
    status: 'COMPLETED',
    label: 'Completed',
    color: 'green',
    icon: 'check',
    description: 'Payment completed successfully',
  },
  FAILED: {
    status: 'FAILED',
    label: 'Failed',
    color: 'red',
    icon: 'alert',
    description: 'Payment failed',
  },
};

// Refund Types
export interface SquareRefund {
  refund_id: number;
  sale_id: number;
  payment_id: string;
  square_refund_id: string;
  amount_money: number;
  status: 'PENDING' | 'COMPLETED' | 'FAILED' | 'REJECTED';
  reason?: string;
  refund_message?: string;
  created_at: string;
  updated_at: string;
  completed_at?: string;
}

export interface RefundCreate {
  sale_id: number;
  amount_cents?: number; // Optional for partial refunds
  reason?: string;
  refund_message?: string;
}

export interface RefundStatusInfo {
  status: SquareRefund['status'];
  label: string;
  color: 'blue' | 'yellow' | 'red' | 'green' | 'gray';
  icon: string;
  description: string;
}

export const REFUND_STATUS_INFO: Record<SquareRefund['status'], RefundStatusInfo> = {
  PENDING: {
    status: 'PENDING',
    label: 'Processing',
    color: 'blue',
    icon: 'clock',
    description: 'Refund is being processed',
  },
  COMPLETED: {
    status: 'COMPLETED',
    label: 'Completed',
    color: 'green',
    icon: 'check',
    description: 'Refund completed successfully',
  },
  FAILED: {
    status: 'FAILED',
    label: 'Failed',
    color: 'red',
    icon: 'alert',
    description: 'Refund failed',
  },
  REJECTED: {
    status: 'REJECTED',
    label: 'Rejected',
    color: 'red',
    icon: 'close',
    description: 'Refund was rejected',
  },
};
