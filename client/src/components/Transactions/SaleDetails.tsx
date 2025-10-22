import { SaleWithItems } from '../../types/sale';
import { Button } from '../UI/Button';
import { ArrowPathIcon, DocumentTextIcon } from '@heroicons/react/24/outline';

interface SaleDetailsProps {
    sale: SaleWithItems;
    onUpdatePaymentStatus: (saleId: number, status: 'pending' | 'partial' | 'paid' | 'refunded') => void;
    onRefund: (sale: SaleWithItems) => void;
}

export const SaleDetails = ({ sale, onUpdatePaymentStatus, onRefund }: SaleDetailsProps) => {
    const handleViewReceipt = () => {
        // Open PDF in new window
        const receiptUrl = `/api/sales/${sale.sale_id}/receipt`;
        window.open(receiptUrl, '_blank');
    };
    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    const formatCurrency = (amount: number) => {
        return `$${amount.toFixed(2)}`;
    };

    const getPaymentStatusBadge = (status: string) => {
        const styles = {
            paid: 'bg-green-100 text-green-800',
            pending: 'bg-yellow-100 text-yellow-800',
            partial: 'bg-blue-100 text-blue-800',
            refunded: 'bg-red-100 text-red-800',
        };
        return (
            <span
                className={`px-3 py-1 text-sm font-medium rounded-full ${
                    styles[status as keyof typeof styles] || 'bg-gray-100 text-gray-800'
                }`}
            >
                {status.charAt(0).toUpperCase() + status.slice(1)}
            </span>
        );
    };

    return (
        <div className="space-y-6">
            {/* Sale Header Info */}
            <div className="grid grid-cols-2 gap-4 pb-4 border-b">
                <div>
                    <label className="text-sm font-medium text-gray-500">Receipt Number</label>
                    <p className="mt-1 text-lg font-semibold text-gray-900">{sale.receipt_no}</p>
                </div>
                <div>
                    <label className="text-sm font-medium text-gray-500">Sale Date</label>
                    <p className="mt-1 text-gray-900">{formatDate(sale.sale_date)}</p>
                </div>
                <div>
                    <label className="text-sm font-medium text-gray-500">Payment Status</label>
                    <p className="mt-1">{getPaymentStatusBadge(sale.payment_status)}</p>
                </div>
                <div>
                    <label className="text-sm font-medium text-gray-500">Payment Method</label>
                    <p className="mt-1 text-gray-900 capitalize">
                        {sale.payment_method?.replace('_', ' ') || '-'}
                    </p>
                </div>
            </div>

            {/* Notes */}
            {sale.notes && (
                <div>
                    <label className="text-sm font-medium text-gray-500">Notes</label>
                    <p className="mt-1 text-gray-900">{sale.notes}</p>
                </div>
            )}

            {/* Sale Items */}
            <div>
                <h3 className="text-lg font-medium text-gray-900 mb-4">Items</h3>
                <div className="space-y-3">
                    {sale.items.map((item) => (
                        <div
                            key={item.sale_item_id}
                            className="bg-gray-50 p-4 rounded-lg flex items-center justify-between"
                        >
                            <div className="flex-1">
                                <div className="font-medium text-gray-900">
                                    {item.product_name || `Product #${item.product_id}`}
                                </div>
                                {item.product_sku && (
                                    <div className="text-sm text-gray-500">SKU: {item.product_sku}</div>
                                )}
                                <div className="text-sm text-gray-600 mt-1">
                                    {item.quantity} × {formatCurrency(item.unit_price)}
                                    {item.discount_percent > 0 && (
                                        <span className="text-red-600 ml-2">
                                            ({item.discount_percent}% off)
                                        </span>
                                    )}
                                    {item.tax_rate > 0 && (
                                        <span className="text-gray-500 ml-2">+ {item.tax_rate}% tax</span>
                                    )}
                                </div>
                            </div>
                            <div className="text-right">
                                <div className="text-lg font-semibold text-gray-900">
                                    {formatCurrency(item.line_total)}
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Totals */}
            <div className="bg-gray-50 p-4 rounded-lg space-y-2">
                <div className="flex justify-between">
                    <span className="text-gray-600">Subtotal:</span>
                    <span className="font-medium">
                        {formatCurrency(sale.total_amount - sale.tax_amount)}
                    </span>
                </div>
                {sale.discount_amount > 0 && (
                    <div className="flex justify-between">
                        <span className="text-gray-600">Discount:</span>
                        <span className="font-medium text-red-600">
                            -{formatCurrency(sale.discount_amount)}
                        </span>
                    </div>
                )}
                {sale.tax_amount > 0 && (
                    <div className="flex justify-between">
                        <span className="text-gray-600">Tax:</span>
                        <span className="font-medium">{formatCurrency(sale.tax_amount)}</span>
                    </div>
                )}
                <div className="flex justify-between text-lg font-semibold pt-2 border-t border-gray-300">
                    <span>Total Amount:</span>
                    <span className="text-black">{formatCurrency(sale.net_amount)}</span>
                </div>
            </div>

            {/* Actions */}
            <div className="flex gap-3 pt-4 border-t">
                <Button
                    variant="secondary"
                    onClick={handleViewReceipt}
                    className="flex-1"
                >
                    <DocumentTextIcon className="h-5 w-5 mr-2" />
                    View Receipt
                </Button>
                {sale.payment_status !== 'refunded' && sale.payment_status !== 'paid' && (
                    <Button
                        variant="success"
                        onClick={() => onUpdatePaymentStatus(sale.sale_id, 'paid')}
                        className="flex-1"
                    >
                        Mark as Paid
                    </Button>
                )}
                {sale.payment_status === 'paid' && (
                    <Button
                        variant="danger"
                        onClick={() => onRefund(sale)}
                        className="flex-1"
                    >
                        <ArrowPathIcon className="h-5 w-5 mr-2" />
                        Process Refund
                    </Button>
                )}
            </div>
        </div>
    );
};
