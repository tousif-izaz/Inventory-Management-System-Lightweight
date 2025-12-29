import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts';

export interface PaymentMethodData {
    method: string;
    transactionCount: number;
    totalAmount: number;
    percentage: number;
}

export interface PaymentMethodChartProps {
    data: PaymentMethodData[];
    onMethodClick?: (method: PaymentMethodData) => void;
    isLoading?: boolean;
}

const COLORS = ['#2563eb', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#06b6d4'];

const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
        const data = payload[0].payload as PaymentMethodData;

        return (
            <div className="bg-white border border-gray-200 shadow-lg rounded-lg p-3">
                <p className="font-medium text-gray-900 capitalize mb-2">{data.method}</p>
                <p className="text-sm text-gray-600">
                    Amount: ${data.totalAmount.toLocaleString('en-US', {
                        minimumFractionDigits: 2,
                        maximumFractionDigits: 2,
                    })}
                </p>
                <p className="text-sm text-gray-600">
                    Transactions: {data.transactionCount}
                </p>
                <p className="text-sm font-medium text-blue-600">
                    {data.percentage.toFixed(1)}% of total
                </p>
            </div>
        );
    }

    return null;
};

const renderCustomLabel = (entry: any) => {
    return `${entry.percentage.toFixed(0)}%`;
};

export const PaymentMethodChart = ({
    data,
    onMethodClick,
    isLoading = false,
}: PaymentMethodChartProps) => {
    if (isLoading) {
        return (
            <div className="bg-white shadow rounded-lg p-6">
                <h3 className="text-lg font-medium text-gray-900 mb-4">
                    Payment Method Distribution
                </h3>
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
                <h3 className="text-lg font-medium text-gray-900 mb-4">
                    Payment Method Distribution
                </h3>
                <div className="h-80 flex items-center justify-center">
                    <p className="text-gray-500">No payment data available</p>
                </div>
            </div>
        );
    }

    return (
        <div className="bg-white shadow rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium text-gray-900">Payment Method Distribution</h3>
                {onMethodClick && (
                    <p className="text-xs text-gray-500">Click on any section for details</p>
                )}
            </div>

            <ResponsiveContainer width="100%" height={320}>
                <PieChart>
                    <Pie
                        data={data as any}
                        cx="50%"
                        cy="50%"
                        labelLine={false}
                        label={renderCustomLabel}
                        outerRadius={100}
                        fill="#8884d8"
                        dataKey="totalAmount"
                        onClick={(entry: any) => {
                            if (onMethodClick) {
                                onMethodClick(entry);
                            }
                        }}
                        cursor={onMethodClick ? 'pointer' : 'default'}
                    >
                        {data.map((_entry, index) => (
                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                        ))}
                    </Pie>
                    <Tooltip content={<CustomTooltip />} />
                    <Legend
                        verticalAlign="bottom"
                        height={36}
                        formatter={(_value: any, entry: any) => (
                            <span className="capitalize">{entry.payload.method}</span>
                        )}
                    />
                </PieChart>
            </ResponsiveContainer>

            <div className="mt-4 border-t border-gray-200 pt-4">
                <div className="grid grid-cols-1 gap-2">
                    {data.map((method, index) => (
                        <div
                            key={method.method}
                            className={`flex items-center justify-between p-2 rounded ${
                                onMethodClick ? 'hover:bg-gray-50 cursor-pointer' : ''
                            }`}
                            onClick={() => onMethodClick && onMethodClick(method)}
                        >
                            <div className="flex items-center gap-2">
                                <div
                                    className="w-3 h-3 rounded-full"
                                    style={{ backgroundColor: COLORS[index % COLORS.length] }}
                                ></div>
                                <span className="text-sm font-medium text-gray-700 capitalize">
                                    {method.method}
                                </span>
                            </div>
                            <div className="text-right">
                                <p className="text-sm font-medium text-gray-900">
                                    ${method.totalAmount.toLocaleString('en-US', {
                                        minimumFractionDigits: 2,
                                        maximumFractionDigits: 2,
                                    })}
                                </p>
                                <p className="text-xs text-gray-500">
                                    {method.transactionCount} transactions
                                </p>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};
