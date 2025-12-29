import { useState, useEffect } from 'react';
import { reportService } from '../../services/reportService';
import { SalesReportDTO } from '../../types/report';

type ReportPeriod = 'daily' | 'weekly' | 'monthly' | 'quarterly' | 'custom';

export const SalesReports = () => {
    const [period, setPeriod] = useState<ReportPeriod>('monthly');
    const [report, setReport] = useState<SalesReportDTO | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Date filters
    const [selectedDate, setSelectedDate] = useState(new Date().toISOString().split('T')[0]);
    const [selectedYear, setSelectedYear] = useState(new Date().getFullYear());
    const [selectedMonth, setSelectedMonth] = useState(new Date().getMonth() + 1);
    const [selectedWeek, setSelectedWeek] = useState(1);
    const [selectedQuarter, setSelectedQuarter] = useState(1);
    const [customStartDate, setCustomStartDate] = useState('');
    const [customEndDate, setCustomEndDate] = useState('');

    useEffect(() => {
        fetchReport();
    }, [period, selectedDate, selectedYear, selectedMonth, selectedWeek, selectedQuarter]);

    const fetchReport = async () => {
        setLoading(true);
        setError(null);
        try {
            let reportData: SalesReportDTO;

            switch (period) {
                case 'daily':
                    reportData = await reportService.getDailySalesReport(selectedDate);
                    break;
                case 'weekly':
                    reportData = await reportService.getWeeklySalesReport(selectedYear, selectedWeek);
                    break;
                case 'monthly':
                    reportData = await reportService.getMonthlySalesReport(selectedYear, selectedMonth);
                    break;
                case 'quarterly':
                    reportData = await reportService.getQuarterlySalesReport(selectedYear, selectedQuarter);
                    break;
                case 'custom':
                    if (customStartDate && customEndDate) {
                        reportData = await reportService.getCustomSalesReport({
                            start_date: customStartDate,
                            end_date: customEndDate,
                        });
                    } else {
                        return;
                    }
                    break;
                default:
                    return;
            }

            setReport(reportData);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to fetch report');
        } finally {
            setLoading(false);
        }
    };

    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
        }).format(amount);
    };

    return (
        <div className="space-y-6">
            <div className="mb-6">
                <h1 className="text-2xl font-semibold text-gray-900">Sales Performance Reports</h1>
                <p className="mt-1 text-sm text-gray-500">Analyze sales trends and performance metrics</p>
            </div>

            {/* Period Selector */}
            <div className="bg-white shadow rounded-lg p-6">
                <h2 className="text-lg font-medium mb-4">Report Period</h2>
                <div className="flex flex-wrap gap-2 mb-4">
                    {(['daily', 'weekly', 'monthly', 'quarterly', 'custom'] as ReportPeriod[]).map((p) => (
                        <button
                            key={p}
                            onClick={() => setPeriod(p)}
                            className={`px-4 py-2 rounded-md font-medium ${
                                period === p
                                    ? 'bg-blue-600 text-white'
                                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                            }`}
                        >
                            {p.charAt(0).toUpperCase() + p.slice(1)}
                        </button>
                    ))}
                </div>

                {/* Period-specific filters */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {period === 'daily' && (
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">Date</label>
                            <input
                                type="date"
                                value={selectedDate}
                                onChange={(e) => setSelectedDate(e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md"
                            />
                        </div>
                    )}

                    {period === 'weekly' && (
                        <>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">Year</label>
                                <input
                                    type="number"
                                    value={selectedYear}
                                    onChange={(e) => setSelectedYear(parseInt(e.target.value))}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">Week</label>
                                <input
                                    type="number"
                                    min="1"
                                    max="53"
                                    value={selectedWeek}
                                    onChange={(e) => setSelectedWeek(parseInt(e.target.value))}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                                />
                            </div>
                        </>
                    )}

                    {period === 'monthly' && (
                        <>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">Year</label>
                                <input
                                    type="number"
                                    value={selectedYear}
                                    onChange={(e) => setSelectedYear(parseInt(e.target.value))}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">Month</label>
                                <select
                                    value={selectedMonth}
                                    onChange={(e) => setSelectedMonth(parseInt(e.target.value))}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                                >
                                    {Array.from({ length: 12 }, (_, i) => i + 1).map((m) => (
                                        <option key={m} value={m}>
                                            {new Date(2024, m - 1).toLocaleDateString('en-US', { month: 'long' })}
                                        </option>
                                    ))}
                                </select>
                            </div>
                        </>
                    )}

                    {period === 'quarterly' && (
                        <>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">Year</label>
                                <input
                                    type="number"
                                    value={selectedYear}
                                    onChange={(e) => setSelectedYear(parseInt(e.target.value))}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">Quarter</label>
                                <select
                                    value={selectedQuarter}
                                    onChange={(e) => setSelectedQuarter(parseInt(e.target.value))}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                                >
                                    <option value={1}>Q1 (Jan-Mar)</option>
                                    <option value={2}>Q2 (Apr-Jun)</option>
                                    <option value={3}>Q3 (Jul-Sep)</option>
                                    <option value={4}>Q4 (Oct-Dec)</option>
                                </select>
                            </div>
                        </>
                    )}

                    {period === 'custom' && (
                        <>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">Start Date</label>
                                <input
                                    type="date"
                                    value={customStartDate}
                                    onChange={(e) => setCustomStartDate(e.target.value)}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-2">End Date</label>
                                <input
                                    type="date"
                                    value={customEndDate}
                                    onChange={(e) => setCustomEndDate(e.target.value)}
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md"
                                />
                            </div>
                            <div className="flex items-end">
                                <button
                                    onClick={fetchReport}
                                    className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                                >
                                    Generate Report
                                </button>
                            </div>
                        </>
                    )}
                </div>
            </div>

            {/* Loading and Error States */}
            {loading && (
                <div className="bg-white shadow rounded-lg p-6 text-center">
                    <p className="text-gray-600">Loading report...</p>
                </div>
            )}

            {error && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-6">
                    <p className="text-red-600">{error}</p>
                </div>
            )}

            {/* Report Display */}
            {!loading && !error && report && (
                <>
                    {/* Key Metrics */}
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-sm font-medium text-gray-500">Total Revenue</h3>
                            <p className="mt-2 text-3xl font-bold text-gray-900">
                                {formatCurrency(report.total_revenue)}
                            </p>
                            <p className="mt-1 text-sm text-gray-500">
                                {report.transaction_count} transactions
                            </p>
                        </div>

                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-sm font-medium text-gray-500">Average Transaction</h3>
                            <p className="mt-2 text-3xl font-bold text-gray-900">
                                {formatCurrency(report.avg_transaction)}
                            </p>
                            <p className="mt-1 text-sm text-gray-500">
                                {report.items_sold} items sold
                            </p>
                        </div>

                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-sm font-medium text-gray-500">Tax Collected</h3>
                            <p className="mt-2 text-3xl font-bold text-gray-900">
                                {formatCurrency(report.tax_collected)}
                            </p>
                            <p className="mt-1 text-sm text-gray-500">
                                Discounts: {formatCurrency(report.discounts_given)}
                            </p>
                        </div>

                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-sm font-medium text-gray-500">Unique Customers</h3>
                            <p className="mt-2 text-3xl font-bold text-gray-900">
                                {report.unique_customers}
                            </p>
                            <p className="mt-1 text-sm text-gray-500">
                                {formatCurrency(report.revenue_per_customer)} avg per customer
                            </p>
                        </div>
                    </div>

                    {/* Top Products */}
                    {report.top_products && report.top_products.length > 0 && (
                        <div className="bg-white shadow rounded-lg p-6">
                            <h2 className="text-lg font-medium mb-4">Top Selling Products</h2>
                            <div className="overflow-x-auto">
                                <table className="min-w-full divide-y divide-gray-200">
                                    <thead className="bg-gray-50">
                                        <tr>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                                Product
                                            </th>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                                SKU
                                            </th>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                                Category
                                            </th>
                                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                                Quantity
                                            </th>
                                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                                Revenue
                                            </th>
                                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                                %
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-white divide-y divide-gray-200">
                                        {report.top_products.map((product) => (
                                            <tr key={product.product_id}>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                                    {product.product_name}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                                    {product.sku}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                                    {product.category_name}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                                                    {product.quantity_sold}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                                                    {formatCurrency(product.revenue)}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-500">
                                                    {product.percentage.toFixed(1)}%
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    )}

                    {/* Payment Methods */}
                    {report.payment_methods && report.payment_methods.length > 0 && (
                        <div className="bg-white shadow rounded-lg p-6">
                            <h2 className="text-lg font-medium mb-4">Payment Methods</h2>
                            <div className="space-y-3">
                                {report.payment_methods.map((method, index) => (
                                    <div key={index} className="flex items-center justify-between">
                                        <div className="flex-1">
                                            <div className="flex items-center justify-between mb-1">
                                                <span className="text-sm font-medium text-gray-900 capitalize">
                                                    {method.method}
                                                </span>
                                                <span className="text-sm text-gray-500">
                                                    {method.percentage.toFixed(1)}%
                                                </span>
                                            </div>
                                            <div className="w-full bg-gray-200 rounded-full h-2">
                                                <div
                                                    className="bg-blue-600 h-2 rounded-full"
                                                    style={{ width: `${method.percentage}%` }}
                                                />
                                            </div>
                                            <div className="flex items-center justify-between mt-1">
                                                <span className="text-xs text-gray-500">
                                                    {method.transaction_count} transactions
                                                </span>
                                                <span className="text-xs text-gray-900 font-medium">
                                                    {formatCurrency(method.total_amount)}
                                                </span>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </>
            )}
        </div>
    );
};
