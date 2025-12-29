import { ProductFormData, Category } from '../../types/product';
import { Input } from '../UI/Input';
import { Select } from '../UI/Select';
import { Button } from '../UI/Button';

interface ProductFormProps {
    formData: ProductFormData;
    setFormData: React.Dispatch<React.SetStateAction<ProductFormData>>;
    formErrors: Partial<Record<keyof ProductFormData, string>>;
    categories: Category[];
    onSubmit: () => void;
    onCancel: () => void;
    isLoading: boolean;
    submitText: string;
}

export const ProductForm = ({
    formData,
    setFormData,
    formErrors,
    categories,
    onSubmit,
    onCancel,
    isLoading,
    submitText,
}: ProductFormProps) => {
    const handleChange = (
        e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>
    ) => {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    return (
        <div className="p-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Input
                    label="Product Name"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    error={formErrors.name}
                    required
                />

                <Input
                    label="SKU"
                    name="sku"
                    value={formData.sku}
                    onChange={handleChange}
                    error={formErrors.sku}
                    required
                />

                <Select
                    label="Category"
                    name="category_id"
                    value={formData.category_id.toString()}
                    onChange={handleChange}
                    error={formErrors.category_id}
                    options={[
                        { value: '', label: 'Select Category' },
                        ...categories.map((c) => ({
                            value: c.category_id.toString(),
                            label: c.name,
                        })),
                    ]}
                    required
                />

                <Input
                    label="Batch Number"
                    name="batch_no"
                    value={formData.batch_no}
                    onChange={handleChange}
                />

                <Input
                    label="Expiry Date"
                    name="expiry_date"
                    type="date"
                    value={formData.expiry_date}
                    onChange={handleChange}
                />

                <Input
                    label="Cost Price"
                    name="cost_price"
                    type="number"
                    step="0.01"
                    value={formData.cost_price.toString()}
                    onChange={handleChange}
                    error={formErrors.cost_price}
                    required
                />

                <Input
                    label="Selling Price"
                    name="selling_price"
                    type="number"
                    step="0.01"
                    value={formData.selling_price.toString()}
                    onChange={handleChange}
                    error={formErrors.selling_price}
                    required
                />

                <Input
                    label="Current Quantity"
                    name="current_quantity"
                    type="number"
                    value={formData.current_quantity.toString()}
                    onChange={handleChange}
                    error={formErrors.current_quantity}
                    required
                />

                <Input
                    label="Minimum Stock Level"
                    name="min_stock_level"
                    type="number"
                    value={formData.min_stock_level.toString()}
                    onChange={handleChange}
                    error={formErrors.min_stock_level}
                    required
                />

                <Input
                    label="Maximum Stock Level"
                    name="max_stock_level"
                    type="number"
                    value={formData.max_stock_level ? formData.max_stock_level.toString() : ''}
                    onChange={handleChange}
                />

                <Input
                    label="Reorder Point"
                    name="reorder_point"
                    type="number"
                    value={formData.reorder_point.toString()}
                    onChange={handleChange}
                    error={formErrors.reorder_point}
                    required
                />

                <Select
                    label="Unit"
                    name="unit"
                    value={formData.unit}
                    onChange={handleChange}
                    error={formErrors.unit}
                    options={[
                        { value: 'pcs', label: 'Pieces (pcs)' },
                        { value: 'kg', label: 'Kilogram (kg)' },
                        { value: 'liter', label: 'Liter' },
                        { value: 'box', label: 'Box' },
                        { value: 'carton', label: 'Carton' },
                    ]}
                    required
                />

                <Input
                    label="Shelf Location"
                    name="shelf_location"
                    value={formData.shelf_location || ''}
                    onChange={handleChange}
                    helperText="Optional: Physical location in warehouse"
                />

                <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-900 mb-2">
                        Description
                    </label>
                    <textarea
                        name="description"
                        rows={3}
                        value={formData.description}
                        onChange={handleChange}
                        className="block w-full rounded-md border-0 px-3 py-2 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm"
                    />
                </div>
            </div>

            <div className="mt-6 flex justify-end gap-3">
                <Button variant="secondary" onClick={onCancel} disabled={isLoading}>
                    Cancel
                </Button>
                <Button onClick={onSubmit} isLoading={isLoading}>
                    {submitText}
                </Button>
            </div>
        </div>
    );
};
