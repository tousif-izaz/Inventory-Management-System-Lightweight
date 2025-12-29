import {
    LineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
} from 'recharts';

export interface TrendDataPoint {
    date: string;
    revenue: number;
    transactionCount?: number;
}

export interface RevenueTrendChartProps {
    data: TrendDataPoint[];
    onDataPointClick?: (dataPoint: TrendDataPoint) => void;
    isLoading?: boolean;
}

const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
        const data = payload[0].payload as TrendDataPoint;

        return (
            <div className="bg-white border border-gray-200 shadow-lg rounded-lg p-3">
                <p className="font-medium text-gray-900 mb-1">{label}</p>
                <p className="text-sm text-blue-600">
                    Revenue: ${data.revenue.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                </p>
                {data.transactionCount !== undefined && (
                    <p className="text-sm text-gray-600">
                        Transactions: {data.transactionCount}
                    </p>
                )}
            </div>
        );
    }

    return null;
};

export const RevenueTrendChart = ({
    data,
    onDataPointClick,
    isLoading = false,
}: RevenueTrendChartProps) => {
    const formatXAxis = (value: string) => {
        const date = new Date(value);
        return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    };

    const formatYAxis = (value: number) => {
        if (value >= 1000) {
            return `$${(value / 1000).toFixed(1)}k`;
        }
        return `$${value}`;
    };

    if (isLoading) {
        return (
            <div className="bg-white shadow rounded-lg p-6">
                <h3 className="text-lg font-medium text-gray-900 mb-4">Revenue Trend</h3>
                <div className="h-80 flex items-center justify-center">
                    <div className="text-center">
                        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
                        <p className="mt-4 text-gray-500">Loading chart data...</p>
                    </div>
                </div>
            </div>
        );
    }

    if (!data || data.length === 0) {
        return (
            <div className="bg-white shadow rounded-lg p-6">
                <h3 className="text-lg font-medium text-gray-900 mb-4">Revenue Trend</h3>
                <div className="h-80 flex items-center justify-center">
                    <p className="text-gray-500">No data available for the selected period</p>
                </div>
            </div>
        );
    }

    return (
        <div className="bg-white shadow rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium text-gray-900">Revenue Trend</h3>
                {onDataPointClick && (
                    <p className="text-xs text-gray-500">Click on any point for details</p>
                )}
            </div>

            <ResponsiveContainer width="100%" height={320}>
                <LineChart
                    data={data}
                    margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                    onClick={(e: any) => {
                        if (e && e.activePayload && e.activePayload.length && onDataPointClick) {
                            const dataPoint = e.activePayload[0].payload as TrendDataPoint;
                            onDataPointClick(dataPoint);
                        }
                    }}
                >
                    <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                    <XAxis
                        dataKey="date"
                        tickFormatter={formatXAxis}
                        stroke="#6b7280"
                        style={{ fontSize: '12px' }}
                    />
                    <YAxis
                        tickFormatter={formatYAxis}
                        stroke="#6b7280"
                        style={{ fontSize: '12px' }}
                    />
                    <Tooltip content={<CustomTooltip />} />
                    <Line
                        type="monotone"
                        dataKey="revenue"
                        stroke="#2563eb"
                        strokeWidth={2}
                        dot={{ fill: '#2563eb', r: 4 }}
                        activeDot={{ r: 6, cursor: onDataPointClick ? 'pointer' : 'default' }}
                    />
                </LineChart>
            </ResponsiveContainer>

            <div className="mt-4 flex items-center gap-6 text-sm text-gray-600">
                <div className="flex items-center gap-2">
                    <div className="w-4 h-0.5 bg-blue-600"></div>
                    <span>Net Revenue</span>
                </div>
                <div className="flex items-center gap-2">
                    <span className="font-medium">
                        Total: ${data.reduce((sum, d) => sum + d.revenue, 0).toLocaleString('en-US', {
                            minimumFractionDigits: 2,
                            maximumFractionDigits: 2,
                        })}
                    </span>
                </div>
            </div>
        </div>
    );
};
