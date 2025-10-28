import { Product } from '../../types/product';
import { Button } from '../UI/Button';
import { EyeIcon, PencilIcon, TrashIcon, ChevronUpIcon, ChevronDownIcon } from '@heroicons/react/24/outline';

type SortField = 'name' | 'sku' | 'cost_price' | 'selling_price' | 'current_quantity';
type SortDirection = 'asc' | 'desc';

interface ProductTableProps {
    products: Product[];
    selectedProducts: number[];
    sortField: SortField;
    sortDirection: SortDirection;
    onSort: (field: SortField) => void;
    onSelectAll: () => void;
    onToggleSelection: (id: number) => void;
    onView: (product: Product) => void;
    onEdit: (product: Product) => void;
    onDelete: (product: Product) => void;
    currentPage: number;
    totalPages: number;
    itemsPerPage: number;
    totalItems: number;
    onPageChange: (page: number) => void;
}

export const ProductTable = ({
    products,
    selectedProducts,
    sortField,
    sortDirection,
    onSort,
    onSelectAll,
    onToggleSelection,
    onView,
    onEdit,
    onDelete,
    currentPage,
    totalPages,
    itemsPerPage,
    totalItems,
    onPageChange,
}: ProductTableProps) => {
    return (
        <div className="bg-white shadow rounded-lg overflow-hidden">
            <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="px-6 py-3 text-left">
                                <input
                                    type="checkbox"
                                    checked={selectedProducts.length === products.length}
                                    onChange={onSelectAll}
                                    className="rounded border-gray-300 text-black focus:ring-black"
                                />
                            </th>
                            <SortableHeader
                                label="Name"
                                field="name"
                                currentField={sortField}
                                direction={sortDirection}
                                onClick={onSort}
                            />
                            <SortableHeader
                                label="SKU"
                                field="sku"
                                currentField={sortField}
                                direction={sortDirection}
                                onClick={onSort}
                            />
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Category
                            </th>
                            <SortableHeader
                                label="Cost Price"
                                field="cost_price"
                                currentField={sortField}
                                direction={sortDirection}
                                onClick={onSort}
                            />
                            <SortableHeader
                                label="Selling Price"
                                field="selling_price"
                                currentField={sortField}
                                direction={sortDirection}
                                onClick={onSort}
                            />
                            <SortableHeader
                                label="Stock"
                                field="current_quantity"
                                currentField={sortField}
                                direction={sortDirection}
                                onClick={onSort}
                            />
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Status
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Actions
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        {products.map((product) => (
                            <tr key={product.product_id} className="hover:bg-gray-50">
                                <td className="px-6 py-4">
                                    <input
                                        type="checkbox"
                                        checked={selectedProducts.includes(product.product_id)}
                                        onChange={() => onToggleSelection(product.product_id)}
                                        className="rounded border-gray-300 text-black focus:ring-black"
                                    />
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="text-sm font-medium text-gray-900">{product.name}</div>
                                    {product.description && (
                                        <div className="text-sm text-gray-500 truncate max-w-xs">
                                            {product.description}
                                        </div>
                                    )}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    {product.sku}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    {product.category_name || '-'}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    ${product.cost_price.toFixed(2)}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    ${product.selling_price.toFixed(2)}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <span
                                        className={`text-sm font-medium ${
                                            (product.current_quantity || 0) <= product.reorder_point
                                                ? 'text-red-600'
                                                : 'text-gray-900'
                                        }`}
                                    >
                                        {product.current_quantity || 0} {product.unit}
                                    </span>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <span
                                        className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                            product.is_active
                                                ? 'bg-green-100 text-green-800'
                                                : 'bg-red-100 text-red-800'
                                        }`}
                                    >
                                        {product.is_active ? 'Active' : 'Inactive'}
                                    </span>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                                    <div className="flex gap-2">
                                        <button
                                            onClick={() => onView(product)}
                                            className="text-indigo-600 hover:text-indigo-900"
                                            title="View Details"
                                        >
                                            <EyeIcon className="h-5 w-5" />
                                        </button>
                                        <button
                                            onClick={() => onEdit(product)}
                                            className="text-blue-600 hover:text-blue-900"
                                            title="Edit"
                                        >
                                            <PencilIcon className="h-5 w-5" />
                                        </button>
                                        <button
                                            onClick={() => onDelete(product)}
                                            className="text-red-600 hover:text-red-900"
                                            title="Delete"
                                        >
                                            <TrashIcon className="h-5 w-5" />
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            {totalPages > 1 && (
                <div className="bg-white px-4 py-3 border-t border-gray-200 sm:px-6">
                    <div className="flex items-center justify-between">
                        <div className="text-sm text-gray-700">
                            Showing <span className="font-medium">{(currentPage - 1) * itemsPerPage + 1}</span> to{' '}
                            <span className="font-medium">{Math.min(currentPage * itemsPerPage, totalItems)}</span> of{' '}
                            <span className="font-medium">{totalItems}</span> results
                        </div>
                        <div className="flex gap-2">
                            <Button
                                variant="secondary"
                                size="sm"
                                onClick={() => onPageChange(Math.max(1, currentPage - 1))}
                                disabled={currentPage === 1}
                            >
                                Previous
                            </Button>
                            <div className="flex gap-1">
                                {Array.from({ length: totalPages }, (_, i) => i + 1)
                                    .filter(
                                        (page) =>
                                            page === 1 ||
                                            page === totalPages ||
                                            Math.abs(page - currentPage) <= 1
                                    )
                                    .map((page, index, array) => (
                                        <div key={page}>
                                            {index > 0 && array[index - 1] !== page - 1 && (
                                                <span className="px-2 py-1">...</span>
                                            )}
                                            <button
                                                onClick={() => onPageChange(page)}
                                                className={`px-3 py-1 text-sm rounded-md ${
                                                    currentPage === page
                                                        ? 'bg-black text-white'
                                                        : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-300'
                                                }`}
                                            >
                                                {page}
                                            </button>
                                        </div>
                                    ))}
                            </div>
                            <Button
                                variant="secondary"
                                size="sm"
                                onClick={() => onPageChange(Math.min(totalPages, currentPage + 1))}
                                disabled={currentPage === totalPages}
                            >
                                Next
                            </Button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

interface SortableHeaderProps {
    label: string;
    field: SortField;
    currentField: SortField;
    direction: SortDirection;
    onClick: (field: SortField) => void;
}

const SortableHeader = ({ label, field, currentField, direction, onClick }: SortableHeaderProps) => (
    <th
        className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
        onClick={() => onClick(field)}
    >
        <div className="flex items-center gap-1">
            {label}
            {currentField === field && (
                direction === 'asc' ? (
                    <ChevronUpIcon className="h-4 w-4" />
                ) : (
                    <ChevronDownIcon className="h-4 w-4" />
                )
            )}
        </div>
    </th>
);
