import { RecentSale } from './RecentSalesTable';
import { exportToCSV, formatDateForCSV, formatCurrencyForCSV } from '../../utils/csvExport';

export interface SalesDetailModalProps {
    isOpen: boolean;
    onClose: () => void;
    title: string;
    description?: string;
    sales: RecentSale[];
    isLoading?: boolean;
}

export const SalesDetailModal = ({
    isOpen,
    onClose,
    title,
    description,
    sales,
    isLoading = false,
}: SalesDetailModalProps) => {
    if (!isOpen) return null;

    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
        }).format(amount);
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    const totalAmount = sales.reduce((sum, sale) => sum + sale.net_amount, 0);
    const totalItems = sales.reduce((sum, sale) => sum + sale.item_count, 0);

    const handleExportCSV = () => {
        exportToCSV({
            filename: `sales-detail-${new Date().toISOString().split('T')[0]}.csv`,
            data: sales.map(sale => ({
                'Receipt No': sale.receipt_no,
                'Sale Date': formatDateForCSV(sale.sale_date),
                'Customer': sale.customer_name || 'Walk-in',
                'Payment Method': sale.payment_method,
                'Status': sale.payment_status,
                'Total Amount': formatCurrencyForCSV(sale.total_amount),
                'Net Amount': formatCurrencyForCSV(sale.net_amount),
                'Items': sale.item_count,
            })),
        });
    };

    return (
        <div className="fixed inset-0 z-50 overflow-y-auto">
            {/* Backdrop */}
            <div
                className="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
                onClick={onClose}
            ></div>

            {/* Modal */}
            <div className="flex min-h-full items-center justify-center p-4">
                <div className="relative bg-white rounded-lg shadow-xl max-w-6xl w-full max-h-[90vh] flex flex-col">
                    {/* Header */}
                    <div className="border-b border-gray-200 px-6 py-4 flex items-center justify-between">
                        <div>
                            <h2 className="text-xl font-semibold text-gray-900">{title}</h2>
                            {description && (
                                <p className="mt-1 text-sm text-gray-500">{description}</p>
                            )}
                        </div>
                        <button
                            onClick={onClose}
                            className="text-gray-400 hover:text-gray-500 transition-colors"
                        >
                            <svg
                                className="w-6 h-6"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d="M6 18L18 6M6 6l12 12"
                                />
                            </svg>
                        </button>
                    </div>

                    {/* Summary Cards */}
                    <div className="px-6 py-4 bg-gray-50 border-b border-gray-200">
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                            <div className="bg-white rounded-lg p-4 shadow-sm">
                                <p className="text-sm text-gray-500">Total Sales</p>
                                <p className="text-2xl font-bold text-gray-900">{sales.length}</p>
                            </div>
                            <div className="bg-white rounded-lg p-4 shadow-sm">
                                <p className="text-sm text-gray-500">Total Revenue</p>
                                <p className="text-2xl font-bold text-green-600">
                                    {formatCurrency(totalAmount)}
                                </p>
                            </div>
                            <div className="bg-white rounded-lg p-4 shadow-sm">
                                <p className="text-sm text-gray-500">Total Items</p>
                                <p className="text-2xl font-bold text-blue-600">{totalItems}</p>
                            </div>
                        </div>
                    </div>

                    {/* Content */}
                    <div className="flex-1 overflow-y-auto px-6 py-4">
                        {isLoading ? (
                            <div className="flex items-center justify-center h-64">
                                <div className="text-center">
                                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
                                    <p className="mt-4 text-gray-500">Loading sales data...</p>
                                </div>
                            </div>
                        ) : sales.length === 0 ? (
                            <div className="flex items-center justify-center h-64">
                                <p className="text-gray-500">No sales found</p>
                            </div>
                        ) : (
                            <div className="overflow-x-auto">
                                <table className="min-w-full divide-y divide-gray-200">
                                    <thead className="bg-gray-50">
                                        <tr>
                                            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                                Receipt No
                                            </th>
                                            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                                Date & Time
                                            </th>
                                            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                                Customer
                                            </th>
                                            <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                                Payment
                                            </th>
                                            <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                                Amount
                                            </th>
                                            <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase">
                                                Items
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-white divide-y divide-gray-200">
                                        {sales.map((sale) => (
                                            <tr key={sale.sale_id} className="hover:bg-gray-50">
                                                <td className="px-4 py-3 whitespace-nowrap">
                                                    <span className="text-sm font-medium text-blue-600">
                                                        {sale.receipt_no}
                                                    </span>
                                                </td>
                                                <td className="px-4 py-3 whitespace-nowrap">
                                                    <span className="text-sm text-gray-900">
                                                        {formatDate(sale.sale_date)}
                                                    </span>
                                                </td>
                                                <td className="px-4 py-3 whitespace-nowrap">
                                                    <span className="text-sm text-gray-900">
                                                        {sale.customer_name || 'Walk-in'}
                                                    </span>
                                                </td>
                                                <td className="px-4 py-3 whitespace-nowrap">
                                                    <div className="flex flex-col">
                                                        <span className="text-sm text-gray-900 capitalize">
                                                            {sale.payment_method}
                                                        </span>
                                                        <span
                                                            className={`text-xs ${
                                                                sale.payment_status === 'paid'
                                                                    ? 'text-green-600'
                                                                    : sale.payment_status === 'pending'
                                                                    ? 'text-yellow-600'
                                                                    : 'text-red-600'
                                                            }`}
                                                        >
                                                            {sale.payment_status}
                                                        </span>
                                                    </div>
                                                </td>
                                                <td className="px-4 py-3 whitespace-nowrap text-right">
                                                    <span className="text-sm font-medium text-gray-900">
                                                        {formatCurrency(sale.net_amount)}
                                                    </span>
                                                </td>
                                                <td className="px-4 py-3 whitespace-nowrap text-center">
                                                    <span className="text-sm text-gray-900">
                                                        {sale.item_count}
                                                    </span>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>

                    {/* Footer */}
                    <div className="border-t border-gray-200 px-6 py-4 flex items-center justify-between">
                        <div className="text-sm text-gray-500">
                            Showing {sales.length} sale{sales.length !== 1 ? 's' : ''}
                        </div>
                        <div className="flex gap-3">
                            <button
                                onClick={handleExportCSV}
                                disabled={sales.length === 0}
                                className="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
                            >
                                <svg
                                    className="w-4 h-4"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                                    />
                                </svg>
                                Export CSV
                            </button>
                            <button
                                onClick={onClose}
                                className="px-4 py-2 bg-gray-200 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-300 transition-colors"
                            >
                                Close
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};
