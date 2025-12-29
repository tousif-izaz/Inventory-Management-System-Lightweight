import { useState, useEffect } from 'react';
import { reportService } from '../../services/reportService';
import { ABCAnalysisDTO } from '../../types/report';

export const ABCAnalysisReport = () => {
    const [analysis, setAnalysis] = useState<ABCAnalysisDTO | null>(null);
    const [selectedClass, setSelectedClass] = useState<'A' | 'B' | 'C' | 'all'>('all');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Filters
    const [startDate, setStartDate] = useState(() => {
        const date = new Date();
        date.setDate(date.getDate() - 90);
        return date.toISOString().split('T')[0];
    });
    const [endDate, setEndDate] = useState(new Date().toISOString().split('T')[0]);

    useEffect(() => {
        fetchAnalysis();
    }, [startDate, endDate]);

    const fetchAnalysis = async () => {
        setLoading(true);
        setError(null);
        try {
            const data = await reportService.getABCAnalysis({
                start_date: startDate,
                end_date: endDate,
            });
            setAnalysis(data);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to fetch ABC analysis');
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

    const getClassColor = (classification: string) => {
        switch (classification) {
            case 'A':
                return 'bg-green-100 text-green-800 border-green-200';
            case 'B':
                return 'bg-yellow-100 text-yellow-800 border-yellow-200';
            case 'C':
                return 'bg-red-100 text-red-800 border-red-200';
            default:
                return 'bg-gray-100 text-gray-800 border-gray-200';
        }
    };

    const getClassBadgeColor = (classification: string) => {
        switch (classification) {
            case 'A':
                return 'bg-green-600 text-white';
            case 'B':
                return 'bg-yellow-600 text-white';
            case 'C':
                return 'bg-red-600 text-white';
            default:
                return 'bg-gray-600 text-white';
        }
    };

    const filteredProducts = analysis?.products.filter(
        (p) => selectedClass === 'all' || p.classification === selectedClass
    );

    return (
        <div className="space-y-6">
            <div className="mb-6">
                <h1 className="text-2xl font-semibold text-gray-900">ABC Analysis</h1>
                <p className="mt-1 text-sm text-gray-500">
                    Classify inventory by revenue impact using the Pareto principle (80/20 rule)
                </p>
            </div>

            {/* Date Range Filter */}
            <div className="bg-white shadow rounded-lg p-6">
                <h2 className="text-lg font-medium mb-4">Analysis Period</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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
                </div>
            </div>

            {/* Loading and Error States */}
            {loading && (
                <div className="bg-white shadow rounded-lg p-6 text-center">
                    <p className="text-gray-600">Performing ABC analysis...</p>
                </div>
            )}

            {error && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-6">
                    <p className="text-red-600">{error}</p>
                </div>
            )}

            {/* Analysis Results */}
            {!loading && !error && analysis && (
                <>
                    {/* Summary Cards */}
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        {/* Class A */}
                        <div className={`rounded-lg p-6 border-2 ${getClassColor('A')}`}>
                            <div className="flex items-center justify-between mb-2">
                                <h3 className="text-lg font-bold">Class A</h3>
                                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getClassBadgeColor('A')}`}>
                                    High Priority
                                </span>
                            </div>
                            <div className="mt-4 space-y-2">
                                <div className="flex justify-between">
                                    <span className="text-sm">Products:</span>
                                    <span className="font-medium">{analysis.class_a.product_count}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-sm">Revenue:</span>
                                    <span className="font-medium">
                                        {formatCurrency(analysis.class_a.total_revenue)}
                                    </span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-sm">Share:</span>
                                    <span className="font-bold text-lg">
                                        {analysis.class_a.revenue_percent.toFixed(1)}%
                                    </span>
                                </div>
                            </div>
                            <p className="mt-4 text-xs italic">{analysis.class_a.recommendation}</p>
                        </div>

                        {/* Class B */}
                        <div className={`rounded-lg p-6 border-2 ${getClassColor('B')}`}>
                            <div className="flex items-center justify-between mb-2">
                                <h3 className="text-lg font-bold">Class B</h3>
                                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getClassBadgeColor('B')}`}>
                                    Medium Priority
                                </span>
                            </div>
                            <div className="mt-4 space-y-2">
                                <div className="flex justify-between">
                                    <span className="text-sm">Products:</span>
                                    <span className="font-medium">{analysis.class_b.product_count}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-sm">Revenue:</span>
                                    <span className="font-medium">
                                        {formatCurrency(analysis.class_b.total_revenue)}
                                    </span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-sm">Share:</span>
                                    <span className="font-bold text-lg">
                                        {analysis.class_b.revenue_percent.toFixed(1)}%
                                    </span>
                                </div>
                            </div>
                            <p className="mt-4 text-xs italic">{analysis.class_b.recommendation}</p>
                        </div>

                        {/* Class C */}
                        <div className={`rounded-lg p-6 border-2 ${getClassColor('C')}`}>
                            <div className="flex items-center justify-between mb-2">
                                <h3 className="text-lg font-bold">Class C</h3>
                                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getClassBadgeColor('C')}`}>
                                    Low Priority
                                </span>
                            </div>
                            <div className="mt-4 space-y-2">
                                <div className="flex justify-between">
                                    <span className="text-sm">Products:</span>
                                    <span className="font-medium">{analysis.class_c.product_count}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-sm">Revenue:</span>
                                    <span className="font-medium">
                                        {formatCurrency(analysis.class_c.total_revenue)}
                                    </span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-sm">Share:</span>
                                    <span className="font-bold text-lg">
                                        {analysis.class_c.revenue_percent.toFixed(1)}%
                                    </span>
                                </div>
                            </div>
                            <p className="mt-4 text-xs italic">{analysis.class_c.recommendation}</p>
                        </div>
                    </div>

                    {/* Product List */}
                    <div className="bg-white shadow rounded-lg p-6">
                        <div className="flex items-center justify-between mb-4">
                            <h2 className="text-lg font-medium">Product Classification</h2>
                            <div className="flex gap-2">
                                <button
                                    onClick={() => setSelectedClass('all')}
                                    className={`px-4 py-2 rounded-md text-sm font-medium ${
                                        selectedClass === 'all'
                                            ? 'bg-gray-900 text-white'
                                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                    }`}
                                >
                                    All
                                </button>
                                <button
                                    onClick={() => setSelectedClass('A')}
                                    className={`px-4 py-2 rounded-md text-sm font-medium ${
                                        selectedClass === 'A'
                                            ? 'bg-green-600 text-white'
                                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                    }`}
                                >
                                    Class A
                                </button>
                                <button
                                    onClick={() => setSelectedClass('B')}
                                    className={`px-4 py-2 rounded-md text-sm font-medium ${
                                        selectedClass === 'B'
                                            ? 'bg-yellow-600 text-white'
                                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                    }`}
                                >
                                    Class B
                                </button>
                                <button
                                    onClick={() => setSelectedClass('C')}
                                    className={`px-4 py-2 rounded-md text-sm font-medium ${
                                        selectedClass === 'C'
                                            ? 'bg-red-600 text-white'
                                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                    }`}
                                >
                                    Class C
                                </button>
                            </div>
                        </div>

                        <div className="overflow-x-auto">
                            <table className="min-w-full divide-y divide-gray-200">
                                <thead className="bg-gray-50">
                                    <tr>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                                            Class
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
                                            Revenue
                                        </th>
                                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                            Revenue %
                                        </th>
                                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                            Cumulative %
                                        </th>
                                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                                            Stock
                                        </th>
                                    </tr>
                                </thead>
                                <tbody className="bg-white divide-y divide-gray-200">
                                    {filteredProducts?.map((product) => (
                                        <tr key={product.product_id} className="hover:bg-gray-50">
                                            <td className="px-6 py-4 whitespace-nowrap">
                                                <span className={`px-2 py-1 rounded text-xs font-bold ${getClassBadgeColor(product.classification)}`}>
                                                    {product.classification}
                                                </span>
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
                                                {formatCurrency(product.total_revenue)}
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                                                {product.revenue_percent.toFixed(2)}%
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-right font-medium text-gray-900">
                                                {product.cumulative_percent.toFixed(1)}%
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
                </>
            )}
        </div>
    );
};
