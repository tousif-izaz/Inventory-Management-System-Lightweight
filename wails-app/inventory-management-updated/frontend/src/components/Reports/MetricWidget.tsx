import { ReactNode } from 'react';

export interface MetricWidgetProps {
    title: string;
    value: string | number;
    subtitle?: string;
    icon?: ReactNode;
    trend?: {
        value: number;
        isPositive: boolean;
    };
    onClick?: () => void;
    isLoading?: boolean;
}

export const MetricWidget = ({
    title,
    value,
    subtitle,
    icon,
    trend,
    onClick,
    isLoading = false,
}: MetricWidgetProps) => {
    const isClickable = !!onClick;

    return (
        <div
            className={`bg-white shadow rounded-lg p-6 transition-all duration-200 ${
                isClickable ? 'cursor-pointer hover:shadow-lg hover:scale-105' : ''
            }`}
            onClick={onClick}
            role={isClickable ? 'button' : undefined}
            tabIndex={isClickable ? 0 : undefined}
            onKeyDown={(e) => {
                if (isClickable && (e.key === 'Enter' || e.key === ' ')) {
                    e.preventDefault();
                    onClick();
                }
            }}
        >
            <div className="flex items-center justify-between">
                <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                        {icon && <div className="text-gray-500">{icon}</div>}
                        <h3 className="text-sm font-medium text-gray-500">{title}</h3>
                    </div>

                    {isLoading ? (
                        <div className="animate-pulse">
                            <div className="h-8 bg-gray-200 rounded w-24 mb-2"></div>
                            {subtitle && <div className="h-4 bg-gray-200 rounded w-32"></div>}
                        </div>
                    ) : (
                        <>
                            <p className="text-3xl font-bold text-gray-900">{value}</p>
                            {subtitle && (
                                <p className="mt-1 text-sm text-gray-500">{subtitle}</p>
                            )}
                        </>
                    )}
                </div>

                {trend && !isLoading && (
                    <div
                        className={`flex items-center gap-1 px-2 py-1 rounded-full text-sm font-medium ${
                            trend.isPositive
                                ? 'bg-green-100 text-green-800'
                                : 'bg-red-100 text-red-800'
                        }`}
                    >
                        <span>{trend.isPositive ? '↑' : '↓'}</span>
                        <span>{Math.abs(trend.value)}%</span>
                    </div>
                )}
            </div>

            {isClickable && !isLoading && (
                <div className="mt-4 text-xs text-blue-600 font-medium flex items-center gap-1">
                    <span>View details</span>
                    <span>→</span>
                </div>
            )}
        </div>
    );
};
