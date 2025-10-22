import { useState, useEffect } from 'react';
import { ExclamationTriangleIcon, ArrowPathIcon, ArrowDownTrayIcon } from '@heroicons/react/24/outline';
import { Button } from '../UI/Button';
import { api } from '../../services/api';
import { Product } from '../../types/product';

export const LowStock = () => {
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchLowStockProducts = async () => {
        try {
            setLoading(true);
            const allProducts = await api.get<Product[]>('/products');

            // Filter for low stock (current_quantity <= reorder_point) and active products
            const lowStock = allProducts
                .filter(p => p.is_active && p.current_quantity <= p.reorder_point)
                .sort((a, b) => {
                    // Sort by urgency: current_quantity / reorder_point ratio
                    const ratioA = a.reorder_point > 0 ? a.current_quantity / a.reorder_point : 0;
                    const ratioB = b.reorder_point > 0 ? b.current_quantity / b.reorder_point : 0;
                    return ratioA - ratioB;
                });

            setProducts(lowStock);
        } catch (error) {
            console.error('Failed to fetch low stock products:', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchLowStockProducts();
    }, []);

    const getStockStatus = (product: Product) => {
        if (product.current_quantity === 0) {
            return { label: 'Out of Stock', color: 'bg-red-100 text-red-800 border-red-300' };
        } else if (product.current_quantity < product.reorder_point * 0.5) {
            return { label: 'Critical', color: 'bg-red-100 text-red-800 border-red-300' };
        } else {
            return { label: 'Low Stock', color: 'bg-yellow-100 text-yellow-800 border-yellow-300' };
        }
    };

    const handleExportCSV = () => {
        if (products.length === 0) {
            return;
        }

        // CSV headers
        const headers = ['Product Name', 'SKU', 'Category', 'Current Stock', 'Unit', 'Reorder Point', 'Status'];

        // Convert products to CSV rows
        const rows = products.map(product => {
            const status = getStockStatus(product);
            return [
                product.name,
                product.sku,
                product.category_name || 'N/A',
                product.current_quantity.toString(),
                product.unit,
                product.reorder_point.toString(),
                status.label
            ];
        });

        // Combine headers and rows
        const csvContent = [
            headers.join(','),
            ...rows.map(row => row.map(field => `"${field}"`).join(','))
        ].join('\n');

        // Create blob and download
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        const url = URL.createObjectURL(blob);

        link.setAttribute('href', url);
        link.setAttribute('download', `low-stock-report-${new Date().toISOString().split('T')[0]}.csv`);
        link.style.visibility = 'hidden';

        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
            </div>
        );
    }

    return (
        <div>
            <div className="mb-6 flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-semibold text-gray-900">Low Stock Alert</h1>
                    <p className="mt-1 text-sm text-gray-500">
                        {products.length} product{products.length !== 1 ? 's' : ''} need{products.length === 1 ? 's' : ''} to be reordered
                    </p>
                </div>
                <div className="flex gap-2">
                    {products.length > 0 && (
                        <Button onClick={handleExportCSV} variant="secondary">
                            <ArrowDownTrayIcon className="h-5 w-5 mr-2" />
                            Export CSV
                        </Button>
                    )}
                    <button
                        onClick={fetchLowStockProducts}
                        className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                        disabled={loading}
                    >
                        <ArrowPathIcon className={`h-5 w-5 ${loading ? 'animate-spin' : ''}`} />
                    </button>
                </div>
            </div>

            {products.length === 0 ? (
                <div className="bg-white shadow rounded-lg p-12 text-center">
                    <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-100 mb-4">
                        <ExclamationTriangleIcon className="h-8 w-8 text-green-600" />
                    </div>
                    <h3 className="text-lg font-medium text-gray-900 mb-2">All Stock Levels Good!</h3>
                    <p className="text-gray-500">No products are currently below reorder point</p>
                </div>
            ) : (
                <div className="bg-white shadow rounded-lg overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Product
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Category
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Current Stock
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Reorder Point
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Status
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                                {products.map((product) => {
                                    const status = getStockStatus(product);
                                    return (
                                        <tr key={product.product_id} className="hover:bg-gray-50">
                                            <td className="px-6 py-4 whitespace-nowrap">
                                                <div>
                                                    <div className="text-sm font-medium text-gray-900">
                                                        {product.name}
                                                    </div>
                                                    <div className="text-sm text-gray-500">
                                                        SKU: {product.sku}
                                                    </div>
                                                </div>
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap">
                                                <div className="text-sm text-gray-900">
                                                    {product.category_name || 'N/A'}
                                                </div>
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap">
                                                <div className="text-sm font-semibold text-gray-900">
                                                    {product.current_quantity} {product.unit}
                                                </div>
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap">
                                                <div className="text-sm text-gray-500">
                                                    {product.reorder_point} {product.unit}
                                                </div>
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap">
                                                <span className={`px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full border ${status.color}`}>
                                                    {status.label}
                                                </span>
                                            </td>
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}
        </div>
    );
};
