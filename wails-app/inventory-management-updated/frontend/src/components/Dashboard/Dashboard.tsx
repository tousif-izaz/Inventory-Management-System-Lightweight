import { useState, useEffect } from 'react';
import { ChartBarIcon, CubeIcon, CurrencyDollarIcon, ExclamationTriangleIcon, ArrowPathIcon } from '@heroicons/react/24/outline';
import { dashboardService, DashboardStats } from '../../services/dashboardService';
import { Product } from '../../types/product';
import { Sale } from '../../types/sale';

interface StatCardProps {
    title: string;
    value: string | number;
    icon: any;
    color: string;
}

const StatCard = ({ title, value, icon: Icon, color }: StatCardProps) => {
    return (
        <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
                <div>
                    <p className="text-sm font-medium text-gray-600">{title}</p>
                    <p className="mt-2 text-3xl font-semibold text-gray-900">{value}</p>
                </div>
                <div className={`p-3 rounded-full ${color}`}>
                    <Icon className="h-8 w-8 text-white" />
                </div>
            </div>
        </div>
    );
};

export const Dashboard = () => {
    const [stats, setStats] = useState<DashboardStats | null>(null);
    const [lowStockProducts, setLowStockProducts] = useState<Product[]>([]);
    const [expiringProducts, setExpiringProducts] = useState<Product[]>([]);
    const [recentSales, setRecentSales] = useState<Sale[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchDashboardData = async () => {
        try {
            setLoading(true);
            const [statsData, lowStock, expiring, recent] = await Promise.all([
                dashboardService.getDashboardStats(),
                dashboardService.getLowStockProducts(),
                dashboardService.getExpiringProducts(30),
                dashboardService.getRecentSales(),
            ]);

            setStats(statsData);
            setLowStockProducts(lowStock);
            setExpiringProducts(expiring);
            setRecentSales(recent);
        } catch (error) {
            console.error('Failed to fetch dashboard data:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchDashboardData();
    }, []);

    const formatCurrency = (amount: number) => `$${amount.toFixed(2)}`;
    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
        });
    };

    const statCards = stats ? [
        {
            title: 'Total Products',
            value: stats.totalProducts,
            icon: CubeIcon,
            color: 'bg-indigo-600',
        },
        {
            title: 'Low Stock Items',
            value: stats.lowStockCount,
            icon: ExclamationTriangleIcon,
            color: 'bg-red-600',
        },
        {
            title: 'Total Inventory Value',
            value: formatCurrency(stats.totalInventoryValue),
            icon: CurrencyDollarIcon,
            color: 'bg-green-600',
        },
        {
            title: "Today's Sales",
            value: formatCurrency(stats.todaysSales),
            icon: ChartBarIcon,
            color: 'bg-blue-600',
        },
    ] : [];

    if (loading) {
        return (
            <div className="flex justify-center items-center h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
            </div>
        );
    }

    return (
        <div>
            {/* Page header */}
            <div className="mb-6 flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-semibold text-gray-900">Dashboard</h1>
                    <p className="mt-1 text-sm text-gray-500">
                        Welcome back! Here's an overview of your inventory.
                    </p>
                </div>
                <button
                    onClick={fetchDashboardData}
                    className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                >
                    <ArrowPathIcon className={`h-5 w-5 ${loading ? 'animate-spin' : ''}`} />
                </button>
            </div>

            {/* Stats grid */}
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4 mb-8">
                {statCards.map((stat, index) => (
                    <StatCard key={index} {...stat} />
                ))}
            </div>

            {/* Alerts section */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                {/* Low Stock Alerts */}
                <div className="bg-white rounded-lg shadow p-6">
                    <h2 className="text-lg font-semibold text-gray-900 mb-4">Low Stock Alerts</h2>
                    {lowStockProducts.length === 0 ? (
                        <div className="text-center py-8">
                            <ExclamationTriangleIcon className="h-12 w-12 text-gray-400 mx-auto mb-2" />
                            <p className="text-gray-500">No low stock alerts</p>
                        </div>
                    ) : (
                        <div className="space-y-3">
                            {lowStockProducts.map((product) => (
                                <div
                                    key={product.product_id}
                                    className="flex items-center justify-between p-3 bg-red-50 rounded-lg border border-red-200"
                                >
                                    <div className="flex-1">
                                        <p className="text-sm font-medium text-gray-900">{product.name}</p>
                                        <p className="text-xs text-gray-500">SKU: {product.sku}</p>
                                    </div>
                                    <div className="text-right">
                                        <p className="text-sm font-semibold text-red-600">
                                            {product.current_quantity} {product.unit}
                                        </p>
                                        <p className="text-xs text-gray-500">
                                            Reorder: {product.reorder_point}
                                        </p>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                {/* Expiring Soon */}
                <div className="bg-white rounded-lg shadow p-6">
                    <h2 className="text-lg font-semibold text-gray-900 mb-4">Expiring Soon (30 days)</h2>
                    {expiringProducts.length === 0 ? (
                        <div className="text-center py-8">
                            <CubeIcon className="h-12 w-12 text-gray-400 mx-auto mb-2" />
                            <p className="text-gray-500">No products expiring soon</p>
                        </div>
                    ) : (
                        <div className="space-y-3">
                            {expiringProducts.map((product) => (
                                <div
                                    key={product.product_id}
                                    className="flex items-center justify-between p-3 bg-yellow-50 rounded-lg border border-yellow-200"
                                >
                                    <div className="flex-1">
                                        <p className="text-sm font-medium text-gray-900">{product.name}</p>
                                        <p className="text-xs text-gray-500">Batch: {product.batch_no || 'N/A'}</p>
                                    </div>
                                    <div className="text-right">
                                        <p className="text-sm font-semibold text-yellow-700">
                                            {product.expiry_date ? formatDate(product.expiry_date) : 'N/A'}
                                        </p>
                                        <p className="text-xs text-gray-500">
                                            {product.current_quantity} {product.unit}
                                        </p>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>

            {/* Recent Sales */}
            <div className="bg-white rounded-lg shadow p-6">
                <h2 className="text-lg font-semibold text-gray-900 mb-4">Recent Sales</h2>
                {recentSales.length === 0 ? (
                    <div className="text-center py-8">
                        <ChartBarIcon className="h-12 w-12 text-gray-400 mx-auto mb-2" />
                        <p className="text-gray-500">No recent sales</p>
                    </div>
                ) : (
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Receipt #
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Date
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Amount
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Status
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                                {recentSales.map((sale) => (
                                    <tr key={sale.sale_id}>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                            {sale.receipt_no}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {formatDate(sale.sale_date)}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-900">
                                            {formatCurrency(sale.net_amount)}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span
                                                className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                                    sale.payment_status === 'paid'
                                                        ? 'bg-green-100 text-green-800'
                                                        : sale.payment_status === 'pending'
                                                        ? 'bg-yellow-100 text-yellow-800'
                                                        : sale.payment_status === 'partial'
                                                        ? 'bg-blue-100 text-blue-800'
                                                        : 'bg-red-100 text-red-800'
                                                }`}
                                            >
                                                {sale.payment_status}
                                            </span>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
};
