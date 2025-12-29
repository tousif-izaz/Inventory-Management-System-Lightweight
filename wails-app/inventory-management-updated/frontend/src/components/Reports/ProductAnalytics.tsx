import { useState, useEffect } from 'react';
import { reportService } from '../../services/reportService';
import { ProductPerformanceDTO, CategoryPerformanceDTO } from '../../types/report';

type ViewMode = 'best-sellers' | 'worst-performers' | 'categories';

export const ProductAnalytics = () => {
    const [viewMode, setViewMode] = useState<ViewMode>('best-sellers');
    const [products, setProducts] = useState<ProductPerformanceDTO[]>([]);
    const [categories, setCategories] = useState<CategoryPerformanceDTO[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Filters
    const [startDate, setStartDate] = useState(() => {
        const date = new Date();
        date.setDate(date.getDate() - 30);
        return date.toISOString().split('T')[0];
    });
    const [endDate, setEndDate] = useState(new Date().toISOString().split('T')[0]);
    const [limit, setLimit] = useState(10);

    useEffect(() => {
        fetchData();
    }, [viewMode, startDate, endDate, limit]);

    const fetchData = async () => {
        setLoading(true);
        setError(null);
        try {
            const params = {
                start_date: startDate,
                end_date: endDate,
                limit: limit,
            };

            if (viewMode === 'best-sellers') {
                const data = await reportService.getBestSellingProducts(params);
                setProducts(data);
            } else if (viewMode === 'worst-performers') {
                const data = await reportService.getWorstPerformingProducts(params);
                setProducts(data);
            } else if (viewMode === 'categories') {
                const data = await reportService.getCategoryPerformance({
                    start_date: startDate,
                    end_date: endDate,
                });
                setCategories(data);
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to fetch data');
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

    const getPerformanceColor = (margin: number) => {
        if (margin >= 30) return 'text-green-600';
        if (margin >= 15) return 'text-yellow-600';
        return 'text-red-600';
    };

    return (
        <div className="space-y-6">
            <div className="mb-6">
                <h1 className="text-2xl font-semibold text-gray-900">Product Performance Analytics</h1>
                <p className="mt-1 text-sm text-gray-500">
                    Identify best and worst performing products to optimize inventory
                </p>
            </div>

            {/* View Mode Selector */}
            <div className="bg-white shadow rounded-lg p-6">
                <h2 className="text-lg font-medium mb-4">Analysis Type</h2>
                <div className="flex flex-wrap gap-2 mb-4">
                    <button
                        onClick={() => setViewMode('best-sellers')}
                        className={`px-4 py-2 rounded-md font-medium ${
                            viewMode === 'best-sellers'
                                ? 'bg-green-600 text-white'
                                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                    >
                        Best Sellers
                    </button>
                    <button
                        onClick={() => setViewMode('worst-performers')}
                        className={`px-4 py-2 rounded-md font-medium ${
                            viewMode === 'worst-performers'
                                ? 'bg-red-600 text-white'
                                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                    >
                        Worst Performers
                    </button>
                    <button
                        onClick={() => setViewMode('categories')}
                        className={`px-4 py-2 rounded-md font-medium ${
                            viewMode === 'categories'
                                ? 'bg-blue-600 text-white'
                                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                    >
                        Category Performance
                    </button>
                </div>

                {/* Filters */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Start Date</label>
                        <input
                            type="date"
                            value={startDate}
                            onChange={(e) => setStartDate(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">End Date</label>
                        <input
                            type="date"
                            value={endDate}
                            onChange={(e) => setEndDate(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md"
                        />
                    </div>
                    {viewMode !== 'categories' && (
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                Number of Products
                            </label>
                            <select
                                value={limit}
                                onChange={(e) => setLimit(parseInt(e.target.value))}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md"
                            >
                                <option value={10}>Top 10</option>
                                <option value={20}>Top 20</option>
                                <option value={50}>Top 50</option>
                                <option value={100}>Top 100</option>
                            </select>
                        </div>
                    )}
                </div>
            </div>

            {/* Loading and Error States */}
            {loading && (
                <div className="bg-white shadow rounded-lg p-6 text-center">
                    <p className="text-gray-600">Loading analytics...</p>
                </div>
            )}

            {error && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-6">
                    <p className="text-red-600">{error}</p>
                </div>
            )}

            {/* Product Performance Table */}
            {!loading && !error && viewMode !== 'categories' && products.length > 0 && (
                <div className="bg-white shadow rounded-lg p-6">
                    <h2 className="text-lg font-medium mb-4">
                        {viewMode === 'best-sellers' ? 'Top Performing Products' : 'Underperforming Products'}
                    </h2>
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                        Rank
                                    </th>
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
                                        Qty Sold
                                    </th>
                                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                        Revenue
                                    </th>
                                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                        Profit
                                    </th>
                                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                        Margin %
                                    </th>
                                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                        Stock
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                                {products.map((product) => (
                                    <tr key={product.product_id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                            #{product.rank}
                                        </td>
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
                                            {formatCurrency(product.total_revenue)}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                                            {formatCurrency(product.gross_profit)}
                                        </td>
                                        <td className={`px-6 py-4 whitespace-nowrap text-sm text-right font-medium ${getPerformanceColor(product.profit_margin)}`}>
                                            {product.profit_margin.toFixed(1)}%
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-500">
                                            {product.current_stock}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            {/* Category Performance */}
            {!loading && !error && viewMode === 'categories' && categories.length > 0 && (
                <div className="bg-white shadow rounded-lg p-6">
                    <h2 className="text-lg font-medium mb-4">Category Performance</h2>
                    <div className="space-y-4">
                        {categories.map((category) => (
                            <div key={category.category_id} className="border rounded-lg p-4">
                                <div className="flex items-center justify-between mb-2">
                                    <h3 className="text-lg font-medium text-gray-900">
                                        {category.category_name}
                                    </h3>
                                    <span className="text-2xl font-bold text-blue-600">
                                        {formatCurrency(category.total_revenue)}
                                    </span>
                                </div>
                                <div className="w-full bg-gray-200 rounded-full h-2 mb-3">
                                    <div
                                        className="bg-blue-600 h-2 rounded-full"
                                        style={{ width: `${category.percentage}%` }}
                                    />
                                </div>
                                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                                    <div>
                                        <p className="text-gray-500">Products</p>
                                        <p className="font-medium text-gray-900">{category.product_count}</p>
                                    </div>
                                    <div>
                                        <p className="text-gray-500">Quantity Sold</p>
                                        <p className="font-medium text-gray-900">{category.quantity_sold}</p>
                                    </div>
                                    <div>
                                        <p className="text-gray-500">Avg Price</p>
                                        <p className="font-medium text-gray-900">
                                            {formatCurrency(category.avg_price)}
                                        </p>
                                    </div>
                                    <div>
                                        <p className="text-gray-500">Revenue Share</p>
                                        <p className="font-medium text-gray-900">
                                            {category.percentage.toFixed(1)}%
                                        </p>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {/* No Data Message */}
            {!loading && !error &&
                ((viewMode !== 'categories' && products.length === 0) ||
                 (viewMode === 'categories' && categories.length === 0)) && (
                <div className="bg-white shadow rounded-lg p-6 text-center">
                    <p className="text-gray-600">No data available for the selected period.</p>
                </div>
            )}
        </div>
    );
};
