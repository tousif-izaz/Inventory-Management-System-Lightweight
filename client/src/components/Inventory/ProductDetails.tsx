import { Product } from '../../types/product';
import { Button } from '../UI/Button';
import { PencilIcon } from '@heroicons/react/24/outline';

interface ProductDetailsProps {
    product: Product;
    onEdit?: () => void;
}

export const ProductDetails = ({ product, onEdit }: ProductDetailsProps) => (
    <div className="p-6">
        {onEdit && (
            <div className="mb-6 flex justify-end">
                <Button onClick={onEdit}>
                    <PencilIcon className="h-4 w-4 mr-2" />
                    Edit Product
                </Button>
            </div>
        )}
        <div className="space-y-6">
            {/* Basic Information */}
            <div>
                <h3 className="text-lg font-medium text-gray-900 mb-3 pb-2 border-b">Basic Information</h3>
                <div className="grid grid-cols-2 gap-4">
                    <DetailRow label="Name" value={product.name} />
                    <DetailRow label="SKU" value={product.sku} />
                    <DetailRow label="Category" value={product.category_name || '-'} />
                    <DetailRow label="Unit" value={product.unit} />
                    <DetailRow
                        label="Status"
                        value={
                            <span
                                className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                    product.is_active
                                        ? 'bg-green-100 text-green-800'
                                        : 'bg-red-100 text-red-800'
                                }`}
                            >
                                {product.is_active ? 'Active' : 'Inactive'}
                            </span>
                        }
                    />
                    <DetailRow label="Shelf Location" value={product.shelf_location || '-'} />
                </div>
                {product.description && (
                    <div className="mt-4">
                        <DetailRow label="Description" value={product.description} />
                    </div>
                )}
            </div>

            {/* Batch & Expiry */}
            <div>
                <h3 className="text-lg font-medium text-gray-900 mb-3 pb-2 border-b">Batch Information</h3>
                <div className="grid grid-cols-2 gap-4">
                    <DetailRow label="Batch Number" value={product.batch_no || '-'} />
                    <DetailRow
                        label="Expiry Date"
                        value={
                            product.expiry_date
                                ? new Date(product.expiry_date).toLocaleDateString()
                                : '-'
                        }
                    />
                </div>
            </div>

            {/* Pricing */}
            <div>
                <h3 className="text-lg font-medium text-gray-900 mb-3 pb-2 border-b">Pricing</h3>
                <div className="grid grid-cols-2 gap-4">
                    <DetailRow label="Cost Price" value={`$${product.cost_price.toFixed(2)}`} />
                    <DetailRow label="Selling Price" value={`$${product.selling_price.toFixed(2)}`} />
                    <DetailRow
                        label="Profit Margin"
                        value={`$${(product.selling_price - product.cost_price).toFixed(2)} (${(
                            ((product.selling_price - product.cost_price) / product.cost_price) *
                            100
                        ).toFixed(1)}%)`}
                    />
                </div>
            </div>

            {/* Inventory Levels */}
            <div>
                <h3 className="text-lg font-medium text-gray-900 mb-3 pb-2 border-b">Inventory Levels</h3>
                <div className="grid grid-cols-2 gap-4">
                    <DetailRow label="Current Stock" value={`${product.current_quantity || 0} ${product.unit}`} />
                    <DetailRow label="Min Stock Level" value={`${product.min_stock_level} ${product.unit}`} />
                    <DetailRow
                        label="Max Stock Level"
                        value={
                            product.max_stock_level
                                ? `${product.max_stock_level} ${product.unit}`
                                : '-'
                        }
                    />
                    <DetailRow label="Reorder Point" value={`${product.reorder_point} ${product.unit}`} />
                    <DetailRow
                        label="Stock Status"
                        value={
                            <span
                                className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                    product.current_quantity <= product.reorder_point
                                        ? 'bg-red-100 text-red-800'
                                        : product.current_quantity <= product.min_stock_level
                                        ? 'bg-yellow-100 text-yellow-800'
                                        : 'bg-green-100 text-green-800'
                                }`}
                            >
                                {product.current_quantity <= product.reorder_point
                                    ? 'Reorder Required'
                                    : product.current_quantity <= product.min_stock_level
                                    ? 'Low Stock'
                                    : 'In Stock'}
                            </span>
                        }
                    />
                </div>
            </div>

            {/* Timestamps */}
            <div>
                <h3 className="text-lg font-medium text-gray-900 mb-3 pb-2 border-b">Timestamps</h3>
                <div className="grid grid-cols-2 gap-4">
                    <DetailRow
                        label="Created At"
                        value={new Date(product.created_at).toLocaleString()}
                    />
                    <DetailRow
                        label="Updated At"
                        value={
                            product.updated_at ? new Date(product.updated_at).toLocaleString() : '-'
                        }
                    />
                </div>
            </div>
        </div>
    </div>
);

const DetailRow = ({ label, value }: { label: string; value: React.ReactNode }) => (
    <div>
        <dt className="text-sm font-medium text-gray-500">{label}</dt>
        <dd className="mt-1 text-sm text-gray-900">{value}</dd>
    </div>
);
