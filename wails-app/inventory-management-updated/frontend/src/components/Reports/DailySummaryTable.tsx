import { exportToCSV, formatDateForCSV, formatCurrencyForCSV } from '../../utils/csvExport';

export interface DailySummary {
    date: string;
    revenue: number;
    transaction_count: number;
    items_sold: number;
}

export interface DailySummaryTableProps {
    data: DailySummary[];
    totalRevenue?: number;
    totalTransactions?: number;
    onDayClick?: (summary: DailySummary) => void;
    isLoading?: boolean;
}

export const DailySummaryTable = ({
    data,
    totalRevenue,
    totalTransactions,
    onDayClick,
    isLoading = false,
}: DailySummaryTableProps) => {
    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
        }).format(amount);
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            weekday: 'short',
            month: 'short',
            day: 'numeric',
            year: 'numeric',
        });
    };

    const handleExportCSV = () => {
        exportToCSV({
            filename: `daily-summary-${new Date().toISOString().split('T')[0]}.csv`,
            data: data.map(day => ({
                'Date': formatDateForCSV(day.date),
                'Revenue': formatCurrencyForCSV(day.revenue),
                'Transactions': day.transaction_count,
                'Items Sold': day.items_sold,
                'Avg Transaction': formatCurrencyForCSV(
                    day.transaction_count > 0 ? day.revenue / day.transaction_count : 0
                ),
            })),
        });
    };

    if (isLoading) {
        return (
            <div className="bg-white shadow rounded-lg p-6">
                <h3 className="text-lg font-medium text-gray-900 mb-4">Daily Summary</h3>
                <div className="h-64 flex items-center justify-center">
                    <div className="text-center">
                        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
                        <p className="mt-4 text-gray-500">Loading daily summary...</p>
                    </div>
                </div>
            </div>
        );
    }

    if (!data || data.length === 0) {
        return (
            <div className="bg-white shadow rounded-lg p-6">
                <div className="flex items-center justify-between mb-4">
                    <h3 className="text-lg font-medium text-gray-900">Daily Summary</h3>
                </div>
                <div className="h-64 flex items-center justify-center">
                    <p className="text-gray-500">No daily summary data available</p>
                </div>
            </div>
        );
    }

    return (
        <div className="bg-white shadow rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium text-gray-900">Daily Summary</h3>
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
                                Date
                            </th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Revenue
                            </th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Transactions
                            </th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Items Sold
                            </th>
                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Avg Transaction
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        {data.map((day, index) => {
                            const avgTransaction =
                                day.transaction_count > 0 ? day.revenue / day.transaction_count : 0;

                            return (
                                <tr
                                    key={index}
                                    className={`transition-colors ${
                                        onDayClick ? 'hover:bg-gray-50 cursor-pointer' : ''
                                    }`}
                                    onClick={() => onDayClick && onDayClick(day)}
                                >
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span className="text-sm font-medium text-gray-900">
                                            {formatDate(day.date)}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-right">
                                        <span className="text-sm font-medium text-gray-900">
                                            {formatCurrency(day.revenue)}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-right">
                                        <span className="text-sm text-gray-900">
                                            {day.transaction_count}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-right">
                                        <span className="text-sm text-gray-900">{day.items_sold}</span>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-right">
                                        <span className="text-sm text-gray-600">
                                            {formatCurrency(avgTransaction)}
                                        </span>
                                    </td>
                                </tr>
                            );
                        })}
                    </tbody>
                    {(totalRevenue !== undefined || totalTransactions !== undefined) && (
                        <tfoot className="bg-gray-50 font-medium">
                            <tr>
                                <td className="px-6 py-4 text-sm text-gray-900">Total</td>
                                <td className="px-6 py-4 text-sm text-gray-900 text-right">
                                    {totalRevenue !== undefined && formatCurrency(totalRevenue)}
                                </td>
                                <td className="px-6 py-4 text-sm text-gray-900 text-right">
                                    {totalTransactions !== undefined && totalTransactions}
                                </td>
                                <td className="px-6 py-4 text-sm text-gray-900 text-right">
                                    {data.reduce((sum, d) => sum + d.items_sold, 0)}
                                </td>
                                <td className="px-6 py-4 text-sm text-gray-900 text-right">
                                    {totalRevenue !== undefined &&
                                        totalTransactions !== undefined &&
                                        totalTransactions > 0 &&
                                        formatCurrency(totalRevenue / totalTransactions)}
                                </td>
                            </tr>
                        </tfoot>
                    )}
                </table>
            </div>

            {onDayClick && (
                <div className="mt-4 text-xs text-gray-500 text-center">
                    Click on any row to view sales for that day
                </div>
            )}
        </div>
    );
};
