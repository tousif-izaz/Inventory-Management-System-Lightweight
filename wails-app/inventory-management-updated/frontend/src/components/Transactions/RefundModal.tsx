import React, { useState } from 'react';
import { XMarkIcon, ExclamationTriangleIcon, CheckCircleIcon } from '@heroicons/react/24/outline';
import { Button } from '../UI/Button';
import { useToast } from '../UI/Toast';
import { processRefund } from '../../services/terminalService';
import { RefundCreate, SquareRefund, REFUND_STATUS_INFO } from '../../types/terminal';

interface RefundModalProps {
  saleId: number;
  saleAmount: number;
  paymentMethod: string;
  onClose: () => void;
  onRefundComplete?: (refund: SquareRefund) => void;
}

export const RefundModal: React.FC<RefundModalProps> = ({
  saleId,
  saleAmount,
  paymentMethod,
  onClose,
  onRefundComplete
}) => {
  const { showToast } = useToast();
  const [loading, setLoading] = useState(false);
  const [refundStatus, setRefundStatus] = useState<SquareRefund | null>(null);
  const [formData, setFormData] = useState({
    amount: (saleAmount / 100).toFixed(2), // Convert cents to dollars
    isPartialRefund: false,
    reason: '',
    refundMessage: ''
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { checked } = e.target;
    setFormData(prev => ({
      ...prev,
      isPartialRefund: checked,
      amount: checked ? '' : (saleAmount / 100).toFixed(2)
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validation
    const amountCents = Math.round(parseFloat(formData.amount) * 100);
    if (isNaN(amountCents) || amountCents <= 0) {
      showToast('error', 'Please enter a valid refund amount');
      return;
    }

    if (amountCents > saleAmount) {
      showToast('error', 'Refund amount cannot exceed sale amount');
      return;
    }

    setLoading(true);

    try {
      const refundRequest: RefundCreate = {
        sale_id: saleId,
        amount_cents: formData.isPartialRefund ? amountCents : undefined,
        reason: formData.reason || undefined,
        refund_message: formData.refundMessage || undefined
      };

      const refund = await processRefund(refundRequest);
      setRefundStatus(refund);

      if (refund.status === 'COMPLETED') {
        showToast('success', 'Refund processed successfully');
        if (onRefundComplete) {
          onRefundComplete(refund);
        }
      } else if (refund.status === 'PENDING') {
        showToast('info', 'Refund is being processed');
      } else {
        showToast('warning', `Refund status: ${refund.status}`);
      }
    } catch (error: any) {
      console.error('Refund error:', error);
      showToast('error', error.message || 'Failed to process refund');
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (cents: number) => {
    return `$${(cents / 100).toFixed(2)}`;
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <h2 className="text-xl font-semibold text-gray-900">Process Refund</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors"
            disabled={loading}
          >
            <XMarkIcon className="w-6 h-6" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6">
          {refundStatus ? (
            // Refund Status Display
            <div className="space-y-4">
              <div className="flex items-center justify-center">
                {refundStatus.status === 'COMPLETED' && (
                  <CheckCircleIcon className="w-16 h-16 text-green-500" />
                )}
                {(refundStatus.status === 'FAILED' || refundStatus.status === 'REJECTED') && (
                  <ExclamationTriangleIcon className="w-16 h-16 text-red-500" />
                )}
                {refundStatus.status === 'PENDING' && (
                  <div className="w-16 h-16 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
                )}
              </div>

              <div className="text-center">
                <p className="text-lg font-semibold text-gray-900">
                  {REFUND_STATUS_INFO[refundStatus.status].label}
                </p>
                <p className="text-sm text-gray-600 mt-1">
                  Refund Amount: {formatCurrency(refundStatus.amount_money)}
                </p>
                {refundStatus.reason && (
                  <p className="text-sm text-gray-600 mt-1">
                    Reason: {refundStatus.reason}
                  </p>
                )}
              </div>

              <div className="bg-gray-50 rounded-lg p-4">
                <p className="text-xs text-gray-600">
                  <strong>Refund ID:</strong> {refundStatus.square_refund_id}
                </p>
                <p className="text-xs text-gray-600 mt-1">
                  <strong>Created:</strong> {new Date(refundStatus.created_at).toLocaleString()}
                </p>
              </div>

              <Button
                onClick={onClose}
                className="w-full"
              >
                Close
              </Button>
            </div>
          ) : (
            // Refund Form
            <form onSubmit={handleSubmit} className="space-y-4">
              {/* Sale Info */}
              <div className="bg-gray-50 rounded-lg p-4">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Original Sale Amount:</span>
                  <span className="font-semibold text-gray-900">{formatCurrency(saleAmount)}</span>
                </div>
                <div className="flex justify-between text-sm mt-1">
                  <span className="text-gray-600">Payment Method:</span>
                  <span className="font-semibold text-gray-900 capitalize">
                    {paymentMethod.replace('_', ' ')}
                  </span>
                </div>
              </div>

              {/* Partial Refund Checkbox */}
              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="isPartialRefund"
                  checked={formData.isPartialRefund}
                  onChange={handleCheckboxChange}
                  className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  disabled={loading}
                />
                <label htmlFor="isPartialRefund" className="ml-2 text-sm text-gray-700">
                  Process partial refund
                </label>
              </div>

              {/* Amount */}
              <div>
                <label htmlFor="amount" className="block text-sm font-medium text-gray-700 mb-1">
                  Refund Amount *
                </label>
                <div className="relative">
                  <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">
                    $
                  </span>
                  <input
                    type="number"
                    id="amount"
                    name="amount"
                    value={formData.amount}
                    onChange={handleInputChange}
                    step="0.01"
                    min="0.01"
                    max={(saleAmount / 100).toFixed(2)}
                    disabled={!formData.isPartialRefund || loading}
                    required
                    className="w-full pl-7 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100 disabled:cursor-not-allowed"
                  />
                </div>
              </div>

              {/* Reason */}
              <div>
                <label htmlFor="reason" className="block text-sm font-medium text-gray-700 mb-1">
                  Reason
                </label>
                <input
                  type="text"
                  id="reason"
                  name="reason"
                  value={formData.reason}
                  onChange={handleInputChange}
                  placeholder="e.g., Customer return, Defective product"
                  disabled={loading}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100"
                />
              </div>

              {/* Refund Message */}
              <div>
                <label htmlFor="refundMessage" className="block text-sm font-medium text-gray-700 mb-1">
                  Internal Note
                </label>
                <textarea
                  id="refundMessage"
                  name="refundMessage"
                  value={formData.refundMessage}
                  onChange={handleInputChange}
                  placeholder="Optional internal note for audit trail"
                  rows={3}
                  disabled={loading}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-100 resize-none"
                />
              </div>

              {/* Warning */}
              <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-3 flex items-start gap-2">
                <ExclamationTriangleIcon className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
                <p className="text-sm text-yellow-800">
                  This will process a refund through Square and restore the inventory. This action cannot be undone.
                </p>
              </div>

              {/* Actions */}
              <div className="flex gap-3 pt-2">
                <Button
                  type="button"
                  onClick={onClose}
                  variant="secondary"
                  className="flex-1"
                  disabled={loading}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  className="flex-1"
                  disabled={loading}
                >
                  {loading ? (
                    <div className="flex items-center justify-center gap-2">
                      <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                      Processing...
                    </div>
                  ) : (
                    'Process Refund'
                  )}
                </Button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
};
