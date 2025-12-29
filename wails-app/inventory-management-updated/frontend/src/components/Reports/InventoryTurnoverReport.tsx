import { useState, useEffect } from 'react';
import { reportService } from '../../services/reportService';
import { InventoryTurnoverDTO, ProductMovementDTO, StockHealthDTO } from '../../types/report';

type ViewMode = 'metrics' | 'movement' | 'health';

export const InventoryTurnoverReport = () => {
    const [viewMode, setViewMode] = useState<ViewMode>('metrics');
    const [turnoverData, setTurnoverData] = useState<InventoryTurnoverDTO | null>(null);
    const [movements, setMovements] = useState<ProductMovementDTO[]>([]);
    const [stockHealth, setStockHealth] = useState<StockHealthDTO | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Filters for movement view
    const [movementFilter, setMovementFilter] = useState<'all' | 'fast' | 'medium' | 'slow' | 'dead'>('all');

    useEffect(() => {
        fetchData();
    }, [viewMode]);

    const fetchData = async () => {
        setLoading(true);
        setError(null);
        try {
            if (viewMode === 'metrics') {
                const data = await reportService.getInventoryTurnover();
                setTurnoverData(data);
            } else if (viewMode === 'movement') {
                const data = await reportService.getProductMovement();
                setMovements(data);
            } else if (viewMode === 'health') {
                const data = await reportService.getStockHealth();
                setStockHealth(data);
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

    const getMovementColor = (movementClass: string) => {
        switch (movementClass) {
            case 'fast':
                return 'bg-green-100 text-green-800';
            case 'medium':
                return 'bg-yellow-100 text-yellow-800';
            case 'slow':
                return 'bg-orange-100 text-orange-800';
            case 'dead':
                return 'bg-red-100 text-red-800';
            default:
                return 'bg-gray-100 text-gray-800';
        }
    };

    const filteredMovements = movements.filter(
        (m) => movementFilter === 'all' || m.movement_class === movementFilter
    );

    return (
        <div className="space-y-6">
            <div className="mb-6">
                <h1 className="text-2xl font-semibold text-gray-900">Inventory Turnover & Performance</h1>
                <p className="mt-1 text-sm text-gray-500">
                    Monitor inventory efficiency and identify optimization opportunities
                </p>
            </div>

            {/* View Mode Selector */}
            <div className="bg-white shadow rounded-lg p-6">
                <h2 className="text-lg font-medium mb-4">Analysis View</h2>
                <div className="flex flex-wrap gap-2">
                    <button
                        onClick={() => setViewMode('metrics')}
                        className={`px-4 py-2 rounded-md font-medium ${
                            viewMode === 'metrics'
                                ? 'bg-blue-600 text-white'
                                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                    >
                        Turnover Metrics
                    </button>
                    <button
                        onClick={() => setViewMode('movement')}
                        className={`px-4 py-2 rounded-md font-medium ${
                            viewMode === 'movement'
                                ? 'bg-blue-600 text-white'
                                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                    >
                        Product Movement
                    </button>
                    <button
                        onClick={() => setViewMode('health')}
                        className={`px-4 py-2 rounded-md font-medium ${
                            viewMode === 'health'
                                ? 'bg-blue-600 text-white'
                                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                    >
                        Stock Health
                    </button>
                </div>
            </div>

            {/* Loading and Error States */}
            {loading && (
                <div className="bg-white shadow rounded-lg p-6 text-center">
                    <p className="text-gray-600">Loading data...</p>
                </div>
            )}

            {error && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-6">
                    <p className="text-red-600">{error}</p>
                </div>
            )}

            {/* Turnover Metrics View */}
            {!loading && !error && viewMode === 'metrics' && turnoverData && (
                <>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-sm font-medium text-gray-500">Turnover Rate</h3>
                            <p className="mt-2 text-3xl font-bold text-gray-900">
                                {turnoverData.turnover_rate.toFixed(2)}x
                            </p>
                            <p className="mt-1 text-sm text-gray-500">times per year</p>
                        </div>

                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-sm font-medium text-gray-500">Days Sales Inventory</h3>
                            <p className="mt-2 text-3xl font-bold text-gray-900">
                                {turnoverData.days_sales_inventory.toFixed(0)}
                            </p>
                            <p className="mt-1 text-sm text-gray-500">days on average</p>
                        </div>

                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-sm font-medium text-gray-500">Stock-to-Sales Ratio</h3>
                            <p className="mt-2 text-3xl font-bold text-gray-900">
                                {turnoverData.stock_to_sales_ratio.toFixed(2)}
                            </p>
                            <p className="mt-1 text-sm text-gray-500">months of inventory</p>
                        </div>

                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-sm font-medium text-gray-500">Avg Inventory Value</h3>
                            <p className="mt-2 text-3xl font-bold text-gray-900">
                                {formatCurrency(turnoverData.avg_inventory_value)}
                            </p>
                            <p className="mt-1 text-sm text-gray-500">COGS: {formatCurrency(turnoverData.cogs)}</p>
                        </div>
                    </div>

                    <div className="bg-white shadow rounded-lg p-6">
                        <h2 className="text-lg font-medium mb-4">Movement Distribution</h2>
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            <div className="border-l-4 border-green-500 pl-4">
                                <p className="text-sm text-gray-500">Fast Moving</p>
                                <p className="text-2xl font-bold text-gray-900">{turnoverData.fast_moving_count}</p>
                                <p className="text-xs text-gray-500">{'<'} 30 days avg</p>
                            </div>
                            <div className="border-l-4 border-yellow-500 pl-4">
                                <p className="text-sm text-gray-500">Medium Moving</p>
                                <p className="text-2xl font-bold text-gray-900">{turnoverData.medium_moving_count}</p>
                                <p className="text-xs text-gray-500">30-90 days avg</p>
                            </div>
                            <div className="border-l-4 border-orange-500 pl-4">
                                <p className="text-sm text-gray-500">Slow Moving</p>
                                <p className="text-2xl font-bold text-gray-900">{turnoverData.slow_moving_count}</p>
                                <p className="text-xs text-gray-500">90-180 days</p>
                            </div>
                            <div className="border-l-4 border-red-500 pl-4">
                                <p className="text-sm text-gray-500">Dead Stock</p>
                                <p className="text-2xl font-bold text-gray-900">{turnoverData.dead_stock_count}</p>
                                <p className="text-xs text-gray-500">{'>'} 180 days</p>
                            </div>
                        </div>
                    </div>
                </>
            )}

            {/* Product Movement View */}
            {!loading && !error && viewMode === 'movement' && movements.length > 0 && (
                <div className="bg-white shadow rounded-lg p-6">
                    <div className="flex items-center justify-between mb-4">
                        <h2 className="text-lg font-medium">Product Movement Analysis</h2>
                        <div className="flex gap-2">
                            {(['all', 'fast', 'medium', 'slow', 'dead'] as const).map((filter) => (
                                <button
                                    key={filter}
                                    onClick={() => setMovementFilter(filter)}
                                    className={`px-3 py-1 rounded text-sm font-medium ${
                                        movementFilter === filter
                                            ? 'bg-blue-600 text-white'
                                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                    }`}
                                >
                                    {filter.charAt(0).toUpperCase() + filter.slice(1)}
                                </button>
                            ))}
                        </div>
                    </div>

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
                                    <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">
                                        Movement
                                    </th>
                                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                        Stock
                                    </th>
                                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                        Value
                                    </th>
                                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                        Sold (30d)
                                    </th>
                                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                        Days Left
                                    </th>
                                    <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">
                                        Reorder
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                                {filteredMovements.map((product) => (
                                    <tr key={product.product_id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                            {product.product_name}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {product.sku}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-center">
                                            <span className={`px-2 py-1 rounded text-xs font-medium ${getMovementColor(product.movement_class)}`}>
                                                {product.movement_class.toUpperCase()}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                                            {product.current_stock}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                                            {formatCurrency(product.stock_value)}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                                            {product.total_sold_30_days}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-500">
                                            {product.days_of_stock_left
                                                ? product.days_of_stock_left.toFixed(0)
                                                : 'N/A'}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-center">
                                            {product.reorder_recommended && (
                                                <span className="text-red-600 font-medium">⚠️ Yes</span>
                                            )}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            {/* Stock Health View */}
            {!loading && !error && viewMode === 'health' && stockHealth && (
                <>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-sm font-medium text-gray-500">Total Products</h3>
                            <p className="mt-2 text-3xl font-bold text-gray-900">{stockHealth.total_products}</p>
                        </div>
                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-sm font-medium text-gray-500">Total Stock Value</h3>
                            <p className="mt-2 text-3xl font-bold text-gray-900">
                                {formatCurrency(stockHealth.total_stock_value)}
                            </p>
                        </div>
                    </div>

                    {/* Low Stock Alert */}
                    {stockHealth.below_reorder_point.length > 0 && (
                        <div className="bg-white shadow rounded-lg p-6">
                            <h2 className="text-lg font-medium mb-4 text-red-600">
                                ⚠️ Products Below Reorder Point ({stockHealth.below_reorder_point.length})
                            </h2>
                            <div className="overflow-x-auto">
                                <table className="min-w-full divide-y divide-gray-200">
                                    <thead className="bg-gray-50">
                                        <tr>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                                Product
                                            </th>
                                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                                Current Stock
                                            </th>
                                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                                Reorder Point
                                            </th>
                                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                                Recommended Qty
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-white divide-y divide-gray-200">
                                        {stockHealth.below_reorder_point.map((product) => (
                                            <tr key={product.product_id}>
                                                <td className="px-6 py-4 text-sm font-medium text-gray-900">
                                                    {product.product_name}
                                                </td>
                                                <td className="px-6 py-4 text-sm text-right text-red-600 font-medium">
                                                    {product.current_stock}
                                                </td>
                                                <td className="px-6 py-4 text-sm text-right text-gray-900">
                                                    {product.reorder_point}
                                                </td>
                                                <td className="px-6 py-4 text-sm text-right text-green-600 font-medium">
                                                    {product.recommended_qty}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    )}

                    {/* Out of Stock */}
                    {stockHealth.out_of_stock.length > 0 && (
                        <div className="bg-white shadow rounded-lg p-6">
                            <h2 className="text-lg font-medium mb-4 text-red-600">
                                Out of Stock ({stockHealth.out_of_stock.length})
                            </h2>
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
                                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                                Days Out of Stock
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-white divide-y divide-gray-200">
                                        {stockHealth.out_of_stock.map((product) => (
                                            <tr key={product.product_id}>
                                                <td className="px-6 py-4 text-sm font-medium text-gray-900">
                                                    {product.product_name}
                                                </td>
                                                <td className="px-6 py-4 text-sm text-gray-500">
                                                    {product.sku}
                                                </td>
                                                <td className="px-6 py-4 text-sm text-right text-red-600 font-medium">
                                                    {product.days_out_of_stock}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    )}

                    {/* Overstock */}
                    {stockHealth.overstock_items.length > 0 && (
                        <div className="bg-white shadow rounded-lg p-6">
                            <h2 className="text-lg font-medium mb-4 text-orange-600">
                                Overstock Items ({stockHealth.overstock_items.length})
                            </h2>
                            <div className="overflow-x-auto">
                                <table className="min-w-full divide-y divide-gray-200">
                                    <thead className="bg-gray-50">
                                        <tr>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                                Product
                                            </th>
                                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                                Current Stock
                                            </th>
                                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                                Max Level
                                            </th>
                                            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                                Stock Value
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-white divide-y divide-gray-200">
                                        {stockHealth.overstock_items.map((product) => (
                                            <tr key={product.product_id}>
                                                <td className="px-6 py-4 text-sm font-medium text-gray-900">
                                                    {product.product_name}
                                                </td>
                                                <td className="px-6 py-4 text-sm text-right text-orange-600 font-medium">
                                                    {product.current_stock}
                                                </td>
                                                <td className="px-6 py-4 text-sm text-right text-gray-900">
                                                    {product.max_stock_level}
                                                </td>
                                                <td className="px-6 py-4 text-sm text-right text-gray-900">
                                                    {formatCurrency(product.stock_value)}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    )}
                </>
            )}
        </div>
    );
};
