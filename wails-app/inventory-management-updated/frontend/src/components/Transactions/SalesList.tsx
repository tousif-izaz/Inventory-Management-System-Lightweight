import { Sale } from '../../types/sale';
import { Button } from '../UI/Button';
import { EyeIcon } from '@heroicons/react/24/outline';

interface SalesListProps {
    sales: Sale[];
    onView: (sale: Sale) => void;
    onUpdatePaymentStatus: (saleId: number, status: 'pending' | 'partial' | 'paid' | 'refunded') => void;
    onRefund: (sale: Sale) => void;
    loading: boolean;
}

export const SalesList = ({ sales, onView, loading }: SalesListProps) => {
    const getPaymentStatusBadge = (status: string) => {
        const styles = {
            paid: 'bg-green-100 text-green-800',
            pending: 'bg-yellow-100 text-yellow-800',
            partial: 'bg-blue-100 text-blue-800',
            refunded: 'bg-red-100 text-red-800',
        };
        return (
            <span
                className={`px-2 py-1 text-xs font-medium rounded-full ${
                    styles[status as keyof typeof styles] || 'bg-gray-100 text-gray-800'
                }`}
            >
                {status.charAt(0).toUpperCase() + status.slice(1)}
            </span>
        );
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    const formatCurrency = (amount: number) => {
        return `$${amount.toFixed(2)}`;
    };

    return (
        <div className="bg-white shadow rounded-lg overflow-hidden">
            <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Receipt No
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Date
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Total Amount
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Payment Status
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Payment Method
                            </th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Actions
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        {sales.map((sale) => (
                            <tr key={sale.sale_id} className="hover:bg-gray-50">
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="text-sm font-medium text-gray-900">{sale.receipt_no}</div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="text-sm text-gray-900">{formatDate(sale.sale_date)}</div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="text-sm font-medium text-gray-900">
                                        {formatCurrency(sale.total_amount)}
                                    </div>
                                    {sale.discount_amount > 0 && (
                                        <div className="text-xs text-gray-500">
                                            Discount: {formatCurrency(sale.discount_amount)}
                                        </div>
                                    )}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    {getPaymentStatusBadge(sale.payment_status)}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="text-sm text-gray-900 capitalize">
                                        {sale.payment_method?.replace('_', ' ') || '-'}
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                    <div className="flex items-center justify-end gap-2">
                                        <Button
                                            variant="secondary"
                                            size="sm"
                                            onClick={() => onView(sale)}
                                            disabled={loading}
                                        >
                                            <EyeIcon className="h-4 w-4 mr-1" />
                                            View
                                        </Button>
                                    </div>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            {sales.length === 0 && (
                <div className="text-center py-12">
                    <p className="text-gray-500">No sales found</p>
                </div>
            )}
        </div>
    );
};
