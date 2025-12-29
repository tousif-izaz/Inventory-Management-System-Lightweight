import { exportToCSV, formatDateForCSV, formatCurrencyForCSV } from '../../utils/csvExport';

export interface RecentSale {
    sale_id: number;
    receipt_no: string;
    sale_date: string;
    customer_id?: number;
    customer_name?: string;
    payment_method: string;
    payment_status: string;
    total_amount: number;
    net_amount: number;
    item_count: number;
}

export interface RecentSalesTableProps {
    sales: RecentSale[];
    onSaleClick?: (sale: RecentSale) => void;
    isLoading?: boolean;
}

export const RecentSalesTable = ({ sales, onSaleClick, isLoading = false }: RecentSalesTableProps) => {
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

    const handleExportCSV = () => {
        exportToCSV({
            filename: `recent-sales-${new Date().toISOString().split('T')[0]}.csv`,
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

    if (isLoading) {
        return (
            <div className="bg-white shadow rounded-lg p-6">
                <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Sales</h3>
                <div className="h-64 flex items-center justify-center">
                    <div className="text-center">
                        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
                        <p className="mt-4 text-gray-500">Loading sales data...</p>
                    </div>
                </div>
            </div>
        );
    }

    if (!sales || sales.length === 0) {
        return (
            <div className="bg-white shadow rounded-lg p-6">
                <div className="flex items-center justify-between mb-4">
                    <h3 className="text-lg font-medium text-gray-900">Recent Sales</h3>
                </div>
                <div className="h-64 flex items-center justify-center">
                    <p className="text-gray-500">No recent sales found</p>
                </div>
            </div>
        );
    }

    return (
        <div className="bg-white shadow rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium text-gray-900">Recent Sales</h3>
                <button
                    onClick={handleExportCSV}
                    className="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 transition-colors flex items-center gap-2"
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
            </div>

            <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Receipt No
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Date & Time
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Customer
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Payment
                            </th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Net Amount
                            </th>
                            <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Items
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        {sales.map((sale) => (
                            <tr
                                key={sale.sale_id}
                                className={`transition-colors ${
                                    onSaleClick ? 'hover:bg-gray-50 cursor-pointer' : ''
                                }`}
                                onClick={() => onSaleClick && onSaleClick(sale)}
                            >
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <span className="text-sm font-medium text-blue-600">
                                        {sale.receipt_no}
                                    </span>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <span className="text-sm text-gray-900">
                                        {formatDate(sale.sale_date)}
                                    </span>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <span className="text-sm text-gray-900">
                                        {sale.customer_name || 'Walk-in'}
                                    </span>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
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
                                <td className="px-6 py-4 whitespace-nowrap text-right">
                                    <span className="text-sm font-medium text-gray-900">
                                        {formatCurrency(sale.net_amount)}
                                    </span>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-center">
                                    <span className="text-sm text-gray-900">{sale.item_count}</span>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            {onSaleClick && (
                <div className="mt-4 text-xs text-gray-500 text-center">
                    Click on any row to view full sale details
                </div>
            )}
        </div>
    );
};
