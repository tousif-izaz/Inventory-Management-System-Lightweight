import { useState, useEffect } from 'react';
import { MetricWidget } from './MetricWidget';
import { TimeRangePicker, TimeRange } from './TimeRangePicker';
import { RevenueTrendChart, TrendDataPoint } from './RevenueTrendChart';
import { PaymentMethodChart, PaymentMethodData } from './PaymentMethodChart';
import { RecentSalesTable, RecentSale } from './RecentSalesTable';
import { DailySummaryTable, DailySummary } from './DailySummaryTable';
import { SalesDetailModal } from './SalesDetailModal';
import { reportService } from '../../services/reportService';

export const SalesAnalytics = () => {
    // State for time range
    const [timeRange, setTimeRange] = useState<TimeRange>(() => {
        const now = new Date();
        const dayOfWeek = now.getDay();
        const daysToMonday = dayOfWeek === 0 ? 6 : dayOfWeek - 1;
        const monday = new Date(now);
        monday.setDate(now.getDate() - daysToMonday);
        monday.setHours(0, 0, 0, 0);

        const sunday = new Date(monday);
        sunday.setDate(monday.getDate() + 6);
        sunday.setHours(23, 59, 59, 999);

        return {
            startDate: monday,
            endDate: sunday,
            type: 'week',
        };
    });

    // State for data
    const [salesReport, setSalesReport] = useState<any>(null);
    const [recentSales, setRecentSales] = useState<RecentSale[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // State for modal
    const [modalState, setModalState] = useState<{
        isOpen: boolean;
        title: string;
        description?: string;
        sales: RecentSale[];
        isLoading: boolean;
    }>({
        isOpen: false,
        title: '',
        description: '',
        sales: [],
        isLoading: false,
    });

    // Fetch data whenever time range changes
    useEffect(() => {
        fetchAllData();
    }, [timeRange]);

    const fetchAllData = async () => {
        setLoading(true);
        setError(null);

        try {
            // Use ISO format for dates (includes time)
            const startDateISO = timeRange.startDate.toISOString();
            const endDateISO = timeRange.endDate.toISOString();

            console.log('Fetching sales data for range:', { startDateISO, endDateISO });

            // Fetch sales report
            const report = await reportService.getCustomSalesReport({
                start_date: startDateISO,
                end_date: endDateISO,
            });

            console.log('Sales report received:', report);
            setSalesReport(report);

            // Fetch recent sales
            const sales = await reportService.getRecentSales({
                limit: 10,
                start_date: startDateISO,
                end_date: endDateISO,
            });

            console.log('Recent sales received:', sales);
            setRecentSales(sales || []);
        } catch (err) {
            console.error('Error fetching sales data:', err);
            setError(err instanceof Error ? err.message : 'Failed to fetch sales data');
        } finally {
            setLoading(false);
        }
    };

    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
        }).format(amount);
    };

    const handleMetricClick = async (metricType: string) => {
        setModalState({
            isOpen: true,
            title: `${metricType} - Detailed View`,
            description: `All sales for the selected period`,
            sales: [],
            isLoading: true,
        });

        try {
            const startDate = timeRange.startDate.toISOString().split('T')[0];
            const endDate = timeRange.endDate.toISOString().split('T')[0];

            const sales = await reportService.getRecentSales({
                limit: 100,
                start_date: startDate,
                end_date: endDate,
            });

            setModalState((prev) => ({
                ...prev,
                sales,
                isLoading: false,
            }));
        } catch (err) {
            console.error('Error fetching detailed sales:', err);
            setModalState((prev) => ({
                ...prev,
                sales: [],
                isLoading: false,
            }));
        }
    };

    const handlePaymentMethodClick = async (method: PaymentMethodData) => {
        setModalState({
            isOpen: true,
            title: `${method.method.toUpperCase()} Payments`,
            description: `${method.transactionCount} transactions totaling ${formatCurrency(
                method.totalAmount
            )}`,
            sales: [],
            isLoading: true,
        });

        try {
            const startDate = timeRange.startDate.toISOString().split('T')[0];
            const endDate = timeRange.endDate.toISOString().split('T')[0];

            const sales = await reportService.getRecentSales({
                limit: 100,
                start_date: startDate,
                end_date: endDate,
                payment_method: method.method,
            });

            setModalState((prev) => ({
                ...prev,
                sales,
                isLoading: false,
            }));
        } catch (err) {
            console.error('Error fetching payment method sales:', err);
            setModalState((prev) => ({
                ...prev,
                sales: [],
                isLoading: false,
            }));
        }
    };

    const handleDataPointClick = async (dataPoint: TrendDataPoint) => {
        setModalState({
            isOpen: true,
            title: `Sales for ${new Date(dataPoint.date).toLocaleDateString('en-US', {
                month: 'long',
                day: 'numeric',
                year: 'numeric',
            })}`,
            description: `${dataPoint.transactionCount || 0} transactions totaling ${formatCurrency(
                dataPoint.revenue
            )}`,
            sales: [],
            isLoading: true,
        });

        try {
            const date = new Date(dataPoint.date);
            const startDate = new Date(date.getFullYear(), date.getMonth(), date.getDate(), 0, 0, 0);
            const endDate = new Date(
                date.getFullYear(),
                date.getMonth(),
                date.getDate(),
                23,
                59,
                59
            );

            const sales = await reportService.getRecentSales({
                limit: 100,
                start_date: startDate.toISOString().split('T')[0],
                end_date: endDate.toISOString().split('T')[0],
            });

            setModalState((prev) => ({
                ...prev,
                sales,
                isLoading: false,
            }));
        } catch (err) {
            console.error('Error fetching sales for date:', err);
            setModalState((prev) => ({
                ...prev,
                sales: [],
                isLoading: false,
            }));
        }
    };

    const handleDayClick = async (summary: DailySummary) => {
        setModalState({
            isOpen: true,
            title: `Sales for ${new Date(summary.date).toLocaleDateString('en-US', {
                month: 'long',
                day: 'numeric',
                year: 'numeric',
            })}`,
            description: `${summary.transaction_count} transactions totaling ${formatCurrency(
                summary.revenue
            )}`,
            sales: [],
            isLoading: true,
        });

        try {
            const date = new Date(summary.date);
            const startDate = new Date(date.getFullYear(), date.getMonth(), date.getDate(), 0, 0, 0);
            const endDate = new Date(
                date.getFullYear(),
                date.getMonth(),
                date.getDate(),
                23,
                59,
                59
            );

            const sales = await reportService.getRecentSales({
                limit: 100,
                start_date: startDate.toISOString().split('T')[0],
                end_date: endDate.toISOString().split('T')[0],
            });

            setModalState((prev) => ({
                ...prev,
                sales,
                isLoading: false,
            }));
        } catch (err) {
            console.error('Error fetching sales for day:', err);
            setModalState((prev) => ({
                ...prev,
                sales: [],
                isLoading: false,
            }));
        }
    };

    // Prepare chart data with safe defaults
    const trendData: TrendDataPoint[] =
        salesReport?.daily_breakdown?.map((day: any) => ({
            date: day.date || new Date().toISOString(),
            revenue: day.revenue || 0,
            transactionCount: day.transaction_count || 0,
        })) || [];

    const paymentMethodData: PaymentMethodData[] =
        salesReport?.payment_methods?.map((pm: any) => ({
            method: pm.method || 'unknown',
            transactionCount: pm.transaction_count || 0,
            totalAmount: pm.total_amount || 0,
            percentage: pm.percentage || 0,
        })) || [];

    return (
        <div className="space-y-6 p-6">
            <div className="mb-6">
                <h1 className="text-2xl font-semibold text-gray-900">Sales Analytics</h1>
                <p className="mt-1 text-sm text-gray-500">
                    Comprehensive sales performance metrics and insights
                </p>
            </div>

            {error && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                    <p className="text-red-600">{error}</p>
                </div>
            )}

            {/* Time Range Picker */}
            <TimeRangePicker value={timeRange} onChange={setTimeRange} />

            {/* Top Metrics */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
                <MetricWidget
                    title="Total Sales"
                    value={formatCurrency(salesReport?.gross_sales || 0)}
                    subtitle={`Before discounts & tax`}
                    onClick={salesReport?.transaction_count > 0 ? () => handleMetricClick('Total Sales') : undefined}
                    isLoading={loading}
                />

                <MetricWidget
                    title="Net Revenue"
                    value={formatCurrency(salesReport?.total_revenue || 0)}
                    subtitle="After discounts & tax"
                    onClick={salesReport?.transaction_count > 0 ? () => handleMetricClick('Net Revenue') : undefined}
                    isLoading={loading}
                />

                <MetricWidget
                    title="Transactions"
                    value={salesReport?.transaction_count || 0}
                    subtitle={`Avg: ${formatCurrency(salesReport?.avg_transaction || 0)}`}
                    onClick={salesReport?.transaction_count > 0 ? () => handleMetricClick('Transactions') : undefined}
                    isLoading={loading}
                />

                <MetricWidget
                    title="Total Tax"
                    value={formatCurrency(salesReport?.tax_collected || 0)}
                    subtitle={`Collected`}
                    onClick={salesReport?.transaction_count > 0 ? () => handleMetricClick('Tax Collected') : undefined}
                    isLoading={loading}
                />

                <MetricWidget
                    title="Total Discounts"
                    value={formatCurrency(salesReport?.discounts_given || 0)}
                    subtitle="Given to customers"
                    onClick={salesReport?.transaction_count > 0 ? () => handleMetricClick('Discounts') : undefined}
                    isLoading={loading}
                />
            </div>

            {/* Charts Row */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <RevenueTrendChart
                    data={trendData}
                    onDataPointClick={trendData.length > 0 ? handleDataPointClick : undefined}
                    isLoading={loading}
                />

                <PaymentMethodChart
                    data={paymentMethodData}
                    onMethodClick={paymentMethodData.length > 0 ? handlePaymentMethodClick : undefined}
                    isLoading={loading}
                />
            </div>

            {/* Tables Row */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <DailySummaryTable
                    data={salesReport?.daily_breakdown || []}
                    totalRevenue={salesReport?.total_revenue}
                    totalTransactions={salesReport?.transaction_count}
                    onDayClick={salesReport?.daily_breakdown?.length > 0 ? handleDayClick : undefined}
                    isLoading={loading}
                />

                <RecentSalesTable
                    sales={recentSales}
                    onSaleClick={(sale) => {
                        setModalState({
                            isOpen: true,
                            title: `Sale Details - ${sale.receipt_no}`,
                            description: `${new Date(sale.sale_date).toLocaleString()}`,
                            sales: [sale],
                            isLoading: false,
                        });
                    }}
                    isLoading={loading}
                />
            </div>

            {/* Detail Modal */}
            <SalesDetailModal
                isOpen={modalState.isOpen}
                onClose={() => setModalState((prev) => ({ ...prev, isOpen: false }))}
                title={modalState.title}
                description={modalState.description}
                sales={modalState.sales}
                isLoading={modalState.isLoading}
            />
        </div>
    );
};
