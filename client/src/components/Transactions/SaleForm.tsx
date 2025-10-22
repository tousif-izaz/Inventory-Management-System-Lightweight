import { useState, useEffect } from 'react';
import { Button } from '../UI/Button';
import { Input } from '../UI/Input';
import { Select } from '../UI/Select';
import { useToast } from '../UI/Toast';
import { saleService } from '../../services/saleService';
import { settingService } from '../../services/settingService';
import { api } from '../../services/api';
import { Product } from '../../types/product';
import { SaleFormData, SaleItemFormData } from '../../types/sale';
import { PlusIcon, TrashIcon } from '@heroicons/react/24/outline';

interface SaleFormProps {
    onSuccess: () => void;
    onCancel: () => void;
}

export const SaleForm = ({ onSuccess, onCancel }: SaleFormProps) => {
    const { showToast } = useToast();
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(false);
    const [defaultTaxRate, setDefaultTaxRate] = useState<string>('0');
    const [skuInput, setSkuInput] = useState<string>('');
    const [skuLoading, setSkuLoading] = useState(false);
    const [formData, setFormData] = useState<SaleFormData>({
        receipt_no: `RCP-${Date.now()}`,
        payment_status: 'paid',
        payment_method: 'cash',
        items: [
            {
                product_id: '',
                quantity: '1',
                unit_price: '',
                tax_rate: '0',
                discount_percent: '0',
            },
        ],
    });

    useEffect(() => {
        const fetchInitialData = async () => {
            try {
                // Fetch products
                const productsData = await api.get<Product[]>('/products?is_active=true');
                setProducts(productsData.filter((p) => p.is_active));

                // Fetch default tax rate from settings
                try {
                    const taxRate = await settingService.getSettingValue('sales_tax_rate');
                    setDefaultTaxRate(taxRate);

                    // Update existing items with the default tax rate
                    setFormData(prev => ({
                        ...prev,
                        items: prev.items.map(item => ({
                            ...item,
                            tax_rate: taxRate
                        }))
                    }));
                } catch (settingError) {
                    console.warn('Failed to load default tax rate, using 0:', settingError);
                }
            } catch (error) {
                console.error('Failed to load initial data:', error);
                const errorMessage = error instanceof Error ? error.message : 'Failed to load data';
                showToast('error', `Failed to load data: ${errorMessage}`);
            }
        };
        fetchInitialData();
    }, [showToast]);

    const handleAddItem = () => {
        setFormData({
            ...formData,
            items: [
                ...formData.items,
                {
                    product_id: '',
                    quantity: '1',
                    unit_price: '',
                    tax_rate: defaultTaxRate,
                    discount_percent: '0',
                },
            ],
        });
    };

    const handleRemoveItem = (index: number) => {
        const newItems = formData.items.filter((_, i) => i !== index);
        setFormData({ ...formData, items: newItems });
    };

    const handleSkuSearch = async (e?: React.FormEvent) => {
        if (e) {
            e.preventDefault();
        }

        if (!skuInput.trim()) {
            showToast('error', 'Please enter a SKU to search');
            return;
        }

        try {
            setSkuLoading(true);
            const product = await api.get<Product>(`/products/sku/${skuInput.trim()}`);

            // Check if product is already in the items list
            const existingItemIndex = formData.items.findIndex(
                (item) => item.product_id === product.product_id.toString()
            );

            if (existingItemIndex >= 0) {
                // Increment quantity of existing item
                const newItems = [...formData.items];
                const currentQty = Number(newItems[existingItemIndex].quantity) || 0;
                newItems[existingItemIndex].quantity = (currentQty + 1).toString();
                setFormData({ ...formData, items: newItems });
                showToast('success', `Increased quantity for ${product.name}`);
            } else {
                // Add as new item
                const newItem: SaleItemFormData = {
                    product_id: product.product_id.toString(),
                    quantity: '1',
                    unit_price: product.selling_price.toString(),
                    tax_rate: defaultTaxRate,
                    discount_percent: '0',
                };
                setFormData({
                    ...formData,
                    items: [...formData.items, newItem],
                });
                showToast('success', `Added ${product.name} to sale`);
            }

            // Clear the SKU input
            setSkuInput('');
        } catch (error) {
            console.error('Failed to find product by SKU:', error);
            const errorMessage = error instanceof Error ? error.message : 'Product not found';
            showToast('error', `SKU search failed: ${errorMessage}`);
        } finally {
            setSkuLoading(false);
        }
    };

    const handleItemChange = (index: number, field: keyof SaleItemFormData, value: string | number) => {
        const newItems = [...formData.items];
        newItems[index] = { ...newItems[index], [field]: value };

        // Auto-fill unit price when product is selected
        if (field === 'product_id' && value) {
            const product = products.find((p) => p.product_id === Number(value));
            if (product) {
                newItems[index].unit_price = product.selling_price.toString();
            }
        }

        setFormData({ ...formData, items: newItems });
    };

    const calculateLineTotal = (item: SaleItemFormData): number => {
        const quantity = Number(item.quantity) || 0;
        const unitPrice = Number(item.unit_price) || 0;
        const taxRate = Number(item.tax_rate) || 0;
        const discountPercent = Number(item.discount_percent) || 0;

        const subtotal = quantity * unitPrice;
        const discountAmount = subtotal * (discountPercent / 100);
        const afterDiscount = subtotal - discountAmount;
        const taxAmount = afterDiscount * (taxRate / 100);
        const lineTotal = afterDiscount + taxAmount;

        return lineTotal;
    };

    const calculateTotals = () => {
        let totalAmount = 0;
        let totalTax = 0;
        let totalDiscount = 0;

        formData.items.forEach((item) => {
            const quantity = Number(item.quantity) || 0;
            const unitPrice = Number(item.unit_price) || 0;
            const taxRate = Number(item.tax_rate) || 0;
            const discountPercent = Number(item.discount_percent) || 0;

            const subtotal = quantity * unitPrice;
            const discountAmount = subtotal * (discountPercent / 100);
            const afterDiscount = subtotal - discountAmount;
            const taxAmount = afterDiscount * (taxRate / 100);

            totalAmount += afterDiscount + taxAmount;
            totalTax += taxAmount;
            totalDiscount += discountAmount;
        });

        return {
            totalAmount: totalAmount.toFixed(2),
            totalTax: totalTax.toFixed(2),
            totalDiscount: totalDiscount.toFixed(2),
            netAmount: totalAmount.toFixed(2),
        };
    };

    const validateForm = (): boolean => {
        if (!formData.receipt_no.trim()) {
            showToast('error', 'Receipt number is required');
            return false;
        }

        if (formData.items.length === 0) {
            showToast('error', 'At least one item is required');
            return false;
        }

        for (let i = 0; i < formData.items.length; i++) {
            const item = formData.items[i];
            if (!item.product_id) {
                showToast('error', `Product is required for item ${i + 1}`);
                return false;
            }
            if (!item.quantity || Number(item.quantity) <= 0) {
                showToast('error', `Valid quantity is required for item ${i + 1}`);
                return false;
            }
            if (!item.unit_price || Number(item.unit_price) < 0) {
                showToast('error', `Valid unit price is required for item ${i + 1}`);
                return false;
            }

            // Check stock availability
            const product = products.find((p) => p.product_id === Number(item.product_id));
            if (product && Number(item.quantity) > product.current_quantity) {
                showToast(
                    'error',
                    `Insufficient stock for ${product.name}. Available: ${product.current_quantity}`
                );
                return false;
            }
        }

        return true;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!validateForm()) return;

        try {
            setLoading(true);

            // Convert datetime-local format to ISO 8601
            let saleDateISO: string | undefined = undefined;
            if (formData.sale_date) {
                // datetime-local gives us format like "2025-10-20T15:30"
                // We need to convert it to ISO 8601 with timezone
                const date = new Date(formData.sale_date);
                saleDateISO = date.toISOString();
            }

            const payload: any = {
                receipt_no: formData.receipt_no,
                payment_status: formData.payment_status,
                payment_method: formData.payment_method || undefined,
                sale_date: saleDateISO,
                notes: formData.notes || undefined,
                customer_id: formData.customer_id ? Number(formData.customer_id) : undefined,
                sold_by: formData.sold_by ? Number(formData.sold_by) : undefined,
                items: formData.items.map((item) => ({
                    product_id: Number(item.product_id),
                    quantity: Number(item.quantity),
                    unit_price: Number(item.unit_price),
                    tax_rate: Number(item.tax_rate),
                    discount_percent: Number(item.discount_percent),
                })),
            };

            // Remove undefined fields
            Object.keys(payload).forEach(key => {
                if (payload[key] === undefined) {
                    delete payload[key];
                }
            });

            await saleService.createSale(payload);
            showToast('success', 'Sale recorded successfully');
            onSuccess();
        } catch (error) {
            console.error('Failed to create sale:', error);
            const errorMessage = error instanceof Error ? error.message : 'Failed to create sale';
            showToast('error', `Sale creation failed: ${errorMessage}`);
        } finally {
            setLoading(false);
        }
    };

    const totals = calculateTotals();

    return (
        <form onSubmit={handleSubmit} className="space-y-6">
            {/* SKU Scanner Section */}
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                    Quick Add by SKU/Barcode
                </label>
                <div className="flex gap-2">
                    <Input
                        placeholder="Scan or enter SKU..."
                        value={skuInput}
                        onChange={(e) => setSkuInput(e.target.value)}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                                e.preventDefault();
                                handleSkuSearch();
                            }
                        }}
                        disabled={skuLoading}
                    />
                    <Button
                        type="button"
                        onClick={() => handleSkuSearch()}
                        disabled={skuLoading || !skuInput.trim()}
                    >
                        {skuLoading ? 'Searching...' : 'Add'}
                    </Button>
                </div>
                <p className="text-xs text-gray-500 mt-1">
                    Scan barcode or type SKU and press Enter to quickly add items
                </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Input
                    label="Receipt Number"
                    value={formData.receipt_no}
                    onChange={(e) => setFormData({ ...formData, receipt_no: e.target.value })}
                    required
                />

                <Select
                    label="Payment Status"
                    value={formData.payment_status}
                    onChange={(e) =>
                        setFormData({
                            ...formData,
                            payment_status: e.target.value as any,
                        })
                    }
                    options={[
                        { value: 'pending', label: 'Pending' },
                        { value: 'partial', label: 'Partial' },
                        { value: 'paid', label: 'Paid' },
                    ]}
                    required
                />

                <Select
                    label="Payment Method"
                    value={formData.payment_method || ''}
                    onChange={(e) =>
                        setFormData({
                            ...formData,
                            payment_method: e.target.value as any,
                        })
                    }
                    options={[
                        { value: '', label: 'Select method' },
                        { value: 'cash', label: 'Cash' },
                        { value: 'card', label: 'Card' },
                        { value: 'mobile', label: 'Mobile Payment' },
                        { value: 'bank_transfer', label: 'Bank Transfer' },
                        { value: 'credit', label: 'Credit' },
                    ]}
                />

                <Input
                    label="Sale Date"
                    type="datetime-local"
                    value={formData.sale_date || ''}
                    onChange={(e) => setFormData({ ...formData, sale_date: e.target.value })}
                />
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Notes</label>
                <textarea
                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                    rows={2}
                    value={formData.notes || ''}
                    onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                    placeholder="Optional notes..."
                />
            </div>

            <div className="border-t pt-4">
                <div className="flex items-center justify-between mb-4">
                    <h3 className="text-lg font-medium text-gray-900">Sale Items</h3>
                    <Button type="button" variant="secondary" size="sm" onClick={handleAddItem}>
                        <PlusIcon className="h-4 w-4 mr-1" />
                        Add Item
                    </Button>
                </div>

                <div className="space-y-4">
                    {formData.items.map((item, index) => {
                        const product = products.find((p) => p.product_id === Number(item.product_id));
                        const lineTotal = calculateLineTotal(item);

                        return (
                            <div key={index} className="bg-gray-50 p-4 rounded-lg">
                                <div className="grid grid-cols-1 md:grid-cols-6 gap-3">
                                    <div className="md:col-span-2">
                                        <Select
                                            label="Product"
                                            value={item.product_id}
                                            onChange={(e) =>
                                                handleItemChange(index, 'product_id', e.target.value)
                                            }
                                            options={[
                                                { value: '', label: 'Select product' },
                                                ...products.map((p) => ({
                                                    value: p.product_id.toString(),
                                                    label: `${p.name} (${p.sku}) - Stock: ${p.current_quantity}`,
                                                })),
                                            ]}
                                            required
                                        />
                                        {product && (
                                            <p className="text-xs text-gray-500 mt-1">
                                                Available: {product.current_quantity} {product.unit}
                                            </p>
                                        )}
                                    </div>

                                    <Input
                                        label="Quantity"
                                        type="number"
                                        min="1"
                                        value={item.quantity}
                                        onChange={(e) =>
                                            handleItemChange(index, 'quantity', e.target.value)
                                        }
                                        required
                                    />

                                    <Input
                                        label="Unit Price"
                                        type="number"
                                        step="0.01"
                                        min="0"
                                        value={item.unit_price}
                                        onChange={(e) =>
                                            handleItemChange(index, 'unit_price', e.target.value)
                                        }
                                        required
                                    />

                                    <Input
                                        label="Tax %"
                                        type="number"
                                        step="0.01"
                                        min="0"
                                        max="100"
                                        value={item.tax_rate}
                                        onChange={(e) =>
                                            handleItemChange(index, 'tax_rate', e.target.value)
                                        }
                                    />

                                    <Input
                                        label="Discount %"
                                        type="number"
                                        step="0.01"
                                        min="0"
                                        max="100"
                                        value={item.discount_percent}
                                        onChange={(e) =>
                                            handleItemChange(index, 'discount_percent', e.target.value)
                                        }
                                    />
                                </div>

                                <div className="flex items-center justify-between mt-3 pt-3 border-t border-gray-200">
                                    <div className="text-sm text-gray-700">
                                        Line Total: <span className="font-semibold">${lineTotal.toFixed(2)}</span>
                                    </div>
                                    {formData.items.length > 1 && (
                                        <Button
                                            type="button"
                                            variant="danger"
                                            size="sm"
                                            onClick={() => handleRemoveItem(index)}
                                        >
                                            <TrashIcon className="h-4 w-4 mr-1" />
                                            Remove
                                        </Button>
                                    )}
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>

            <div className="bg-gray-50 p-4 rounded-lg space-y-2">
                <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Subtotal:</span>
                    <span className="font-medium">
                        ${(Number(totals.totalAmount) - Number(totals.totalTax)).toFixed(2)}
                    </span>
                </div>
                <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Discount:</span>
                    <span className="font-medium text-red-600">-${totals.totalDiscount}</span>
                </div>
                <div className="flex justify-between text-sm">
                    <span className="text-gray-600">Tax:</span>
                    <span className="font-medium">${totals.totalTax}</span>
                </div>
                <div className="flex justify-between text-lg font-semibold pt-2 border-t border-gray-300">
                    <span>Total:</span>
                    <span className="text-black">${totals.netAmount}</span>
                </div>
            </div>

            <div className="flex gap-3 justify-end pt-4 border-t">
                <Button type="button" variant="secondary" onClick={onCancel} disabled={loading}>
                    Cancel
                </Button>
                <Button type="submit" disabled={loading}>
                    {loading ? 'Recording Sale...' : 'Record Sale'}
                </Button>
            </div>
        </form>
    );
};
