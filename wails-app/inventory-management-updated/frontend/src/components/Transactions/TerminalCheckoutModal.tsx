import { useState, useEffect } from 'react';
import { Button } from '../UI/Button';
import { useToast } from '../UI/Toast';
import {
  getTerminalCheckout,
  cancelTerminalCheckout,
  waitForCheckoutCompletion,
} from '../../services/terminalService';
import type { TerminalCheckout } from '../../types/terminal';
import { CHECKOUT_STATUS_INFO } from '../../types/terminal';
import {
  XMarkIcon,
  CheckCircleIcon,
  XCircleIcon,
  ClockIcon,
  ArrowPathIcon,
} from '@heroicons/react/24/outline';

interface TerminalCheckoutModalProps {
  checkoutId: number;
  onClose: () => void;
  onComplete?: (checkout: TerminalCheckout) => void;
}

export const TerminalCheckoutModal = ({
  checkoutId,
  onClose,
  onComplete,
}: TerminalCheckoutModalProps) => {
  const { showToast } = useToast();
  const [checkout, setCheckout] = useState<TerminalCheckout | null>(null);
  const [loading, setLoading] = useState(true);
  const [canceling, setCanceling] = useState(false);

  useEffect(() => {
    loadCheckout();
  }, [checkoutId]);

  const loadCheckout = async () => {
    try {
      setLoading(true);
      const data = await getTerminalCheckout(checkoutId);
      setCheckout(data);

      // Start polling if checkout is pending or in progress
      if (data.status === 'PENDING' || data.status === 'IN_PROGRESS') {
        startPolling();
      }
    } catch (error) {
      console.error('Failed to load checkout:', error);
      showToast('error', 'Failed to load checkout status');
    } finally {
      setLoading(false);
    }
  };

  const startPolling = async () => {
    try {
      const completedCheckout = await waitForCheckoutCompletion(
        checkoutId,
        (updatedCheckout) => {
          setCheckout(updatedCheckout);
        }
      );

      if (completedCheckout.status === 'COMPLETED') {
        showToast('success', 'Payment completed successfully!');
        if (onComplete) {
          onComplete(completedCheckout);
        }
      } else if (completedCheckout.status === 'FAILED') {
        showToast('error', 'Payment failed: ' + (completedCheckout.error_message || 'Unknown error'));
      }
    } catch (error) {
      console.error('Polling error:', error);
      showToast('error', 'Failed to monitor payment status');
    }
  };

  const handleCancel = async () => {
    if (!confirm('Are you sure you want to cancel this payment?')) {
      return;
    }

    try {
      setCanceling(true);
      await cancelTerminalCheckout(checkoutId);
      showToast('info', 'Payment canceled');
      await loadCheckout();
    } catch (error: any) {
      console.error('Failed to cancel checkout:', error);
      showToast('error', error.response?.data?.error || 'Failed to cancel payment');
    } finally {
      setCanceling(false);
    }
  };

  if (loading || !checkout) {
    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div className="bg-white rounded-lg p-6 max-w-md w-full">
          <div className="text-center">
            <ArrowPathIcon className="w-12 h-12 animate-spin mx-auto text-blue-600" />
            <p className="mt-4 text-gray-600">Loading checkout status...</p>
          </div>
        </div>
      </div>
    );
  }

  const statusInfo = CHECKOUT_STATUS_INFO[checkout.status];
  const isTerminalState =
    checkout.status === 'COMPLETED' ||
    checkout.status === 'CANCELED' ||
    checkout.status === 'FAILED';

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 className="text-lg font-semibold">Terminal Payment</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
            disabled={!isTerminalState}
          >
            <XMarkIcon className="w-6 h-6" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6">
          {/* Status Icon */}
          <div className="flex justify-center mb-4">
            {checkout.status === 'COMPLETED' && (
              <CheckCircleIcon className="w-20 h-20 text-green-600" />
            )}
            {(checkout.status === 'FAILED' || checkout.status === 'CANCELED') && (
              <XCircleIcon className="w-20 h-20 text-red-600" />
            )}
            {(checkout.status === 'PENDING' || checkout.status === 'IN_PROGRESS') && (
              <ArrowPathIcon className="w-20 h-20 text-blue-600 animate-spin" />
            )}
            {checkout.status === 'CANCEL_REQUESTED' && (
              <ClockIcon className="w-20 h-20 text-yellow-600" />
            )}
          </div>

          {/* Status Label */}
          <div className="text-center mb-6">
            <span
              className={`inline-flex px-4 py-2 text-lg font-semibold rounded-full ${
                statusInfo.color === 'green'
                  ? 'bg-green-100 text-green-800'
                  : statusInfo.color === 'red'
                  ? 'bg-red-100 text-red-800'
                  : statusInfo.color === 'blue'
                  ? 'bg-blue-100 text-blue-800'
                  : statusInfo.color === 'yellow'
                  ? 'bg-yellow-100 text-yellow-800'
                  : 'bg-gray-100 text-gray-800'
              }`}
            >
              {statusInfo.label}
            </span>
            <p className="mt-3 text-gray-600">{statusInfo.description}</p>
          </div>

          {/* Amount */}
          <div className="bg-gray-50 rounded-lg p-4 mb-4">
            <div className="text-sm text-gray-600">Amount</div>
            <div className="text-2xl font-bold text-gray-900">
              ${(checkout.amount_money / 100).toFixed(2)}
            </div>
          </div>

          {/* Checkout ID */}
          <div className="text-sm text-gray-600 mb-4">
            <span className="font-medium">Checkout ID:</span>
            <span className="ml-2 font-mono text-xs">{checkout.square_checkout_id}</span>
          </div>

          {/* Error Message */}
          {checkout.error_message && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-3 mb-4">
              <p className="text-sm text-red-800">{checkout.error_message}</p>
            </div>
          )}

          {/* Instructions for pending/in-progress */}
          {!isTerminalState && (
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
              <h4 className="font-medium text-blue-900 mb-2">Instructions:</h4>
              <ul className="list-disc list-inside text-sm text-blue-800 space-y-1">
                <li>Customer should complete payment on the Terminal device</li>
                <li>This dialog will update automatically when payment is complete</li>
                <li>Do not close this dialog until payment is processed</li>
              </ul>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex gap-2 p-6 border-t border-gray-200">
          {!isTerminalState && (
            <Button
              variant="secondary"
              onClick={handleCancel}
              disabled={canceling || checkout.status === 'CANCEL_REQUESTED'}
              className="flex-1"
            >
              {canceling ? 'Canceling...' : 'Cancel Payment'}
            </Button>
          )}
          {isTerminalState && (
            <Button onClick={onClose} className="flex-1">
              Close
            </Button>
          )}
        </div>
      </div>
    </div>
  );
};
