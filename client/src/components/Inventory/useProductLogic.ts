import { useState, useEffect, useCallback, useMemo } from 'react';
import { Product, ProductFormData, Category } from '../../types/product';
import { api } from '../../services/api';
import { useToast } from '../UI/Toast';

type SortField = 'name' | 'sku' | 'cost_price' | 'selling_price' | 'current_quantity';
type SortDirection = 'asc' | 'desc';

const createEmptyFormData = (): ProductFormData => ({
    name: '',
    description: '',
    sku: '',
    category_id: '',
    batch_no: '',
    expiry_date: '',
    cost_price: '',
    selling_price: '',
    current_quantity: '0',
    min_stock_level: '',
    max_stock_level: '',
    reorder_point: '',
    unit: 'pcs',
    shelf_location: '',
});

export const useProductLogic = () => {
    const { showToast } = useToast();
    const [products, setProducts] = useState<Product[]>([]);
    const [categories, setCategories] = useState<Category[]>([]);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState(false);
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage] = useState(25);
    const [searchTerm, setSearchTerm] = useState('');
    const [categoryFilter, setCategoryFilter] = useState<string>('');
    const [statusFilter, setStatusFilter] = useState<string>('');
    const [sortField, setSortField] = useState<SortField>('name');
    const [sortDirection, setSortDirection] = useState<SortDirection>('asc');
    const [selectedProducts, setSelectedProducts] = useState<number[]>([]);
    const [currentProduct, setCurrentProduct] = useState<Product | null>(null);
    const [formData, setFormData] = useState<ProductFormData>(createEmptyFormData());
    const [formErrors, setFormErrors] = useState<Partial<Record<keyof ProductFormData, string>>>({});

    const fetchProducts = useCallback(async () => {
        try {
            setLoading(true);
            const data = await api.get<Product[]>('/products');
            setProducts(data);
        } catch (error) {
            showToast('error', (error as Error).message);
        } finally {
            setLoading(false);
        }
    }, [showToast]);

    const fetchCategories = useCallback(async () => {
        try {
            const data = await api.get<Category[]>('/categories');
            setCategories(data);
        } catch (error) {
            showToast('error', 'Failed to load categories');
        }
    }, [showToast]);

    useEffect(() => {
        fetchProducts();
        fetchCategories();
    }, [fetchProducts, fetchCategories]);

    const filteredAndSortedProducts = useMemo(() => {
        return products
            .filter((product) => {
                const matchesSearch =
                    product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                    product.sku.toLowerCase().includes(searchTerm.toLowerCase());
                const matchesCategory = categoryFilter ? product.category_id === parseInt(categoryFilter) : true;
                const matchesStatus = statusFilter ? (statusFilter === 'active' ? product.is_active : !product.is_active) : true;
                return matchesSearch && matchesCategory && matchesStatus;
            })
            .sort((a, b) => {
                const aValue = a[sortField];
                const bValue = b[sortField];
                if (aValue == null) return 1;
                if (bValue == null) return -1;
                return sortDirection === 'asc' ? (aValue > bValue ? 1 : -1) : (aValue < bValue ? 1 : -1);
            });
    }, [products, searchTerm, categoryFilter, statusFilter, sortField, sortDirection]);

    const paginatedProducts = useMemo(() => {
        const start = (currentPage - 1) * itemsPerPage;
        const end = currentPage * itemsPerPage;
        return filteredAndSortedProducts.slice(start, end);
    }, [filteredAndSortedProducts, currentPage, itemsPerPage]);

    const totalPages = Math.ceil(filteredAndSortedProducts.length / itemsPerPage);

    const handleSort = (field: SortField) => {
        if (sortField === field) {
            setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
        } else {
            setSortField(field);
            setSortDirection('asc');
        }
    };

    const validateForm = (): boolean => {
        const errors: Partial<Record<keyof ProductFormData, string>> = {};
        if (!formData.name.trim()) errors.name = 'Name is required';
        if (!formData.sku.trim()) errors.sku = 'SKU is required';
        if (!formData.category_id) errors.category_id = 'Category is required';
        if (!formData.cost_price || parseFloat(formData.cost_price.toString()) < 0)
            errors.cost_price = 'Valid cost price is required';
        if (!formData.selling_price || parseFloat(formData.selling_price.toString()) < 0)
            errors.selling_price = 'Valid selling price is required';
        if (!formData.min_stock_level || parseInt(formData.min_stock_level.toString()) < 0)
            errors.min_stock_level = 'Valid minimum stock level is required';
        if (!formData.reorder_point || parseInt(formData.reorder_point.toString()) < 0)
            errors.reorder_point = 'Valid reorder point is required';
        if (!formData.unit.trim()) errors.unit = 'Unit is required';
        setFormErrors(errors);
        return Object.keys(errors).length === 0;
    };

    const handleAddProduct = async () => {
        if (!validateForm()) return;
        try {
            setActionLoading(true);

            const payload: any = {
                name: formData.name,
                sku: formData.sku,
                category_id: parseInt(formData.category_id.toString()),
                cost_price: parseFloat(formData.cost_price.toString()),
                selling_price: parseFloat(formData.selling_price.toString()),
                current_quantity: parseInt(formData.current_quantity.toString()),
                min_stock_level: parseInt(formData.min_stock_level.toString()),
                reorder_point: parseInt(formData.reorder_point.toString()),
                unit: formData.unit,
            };

            // Add optional fields only if they have values
            if (formData.description?.trim()) payload.description = formData.description;
            if (formData.batch_no?.trim()) payload.batch_no = formData.batch_no;
            if (formData.expiry_date?.trim()) payload.expiry_date = formData.expiry_date;
            if (formData.shelf_location?.trim()) payload.shelf_location = formData.shelf_location;
            if (formData.max_stock_level) payload.max_stock_level = parseInt(formData.max_stock_level.toString());

            await api.post('/products', payload);
            showToast('success', 'Product added successfully');
            resetForm();
            fetchProducts();
            return true;
        } catch (error) {
            showToast('error', (error as Error).message);
            return false;
        } finally {
            setActionLoading(false);
        }
    };

    const handleEditProduct = async () => {
        if (!validateForm() || !currentProduct) return;
        try {
            setActionLoading(true);

            const payload: any = {
                name: formData.name,
                sku: formData.sku,
                category_id: parseInt(formData.category_id.toString()),
                cost_price: parseFloat(formData.cost_price.toString()),
                selling_price: parseFloat(formData.selling_price.toString()),
                current_quantity: parseInt(formData.current_quantity.toString()),
                min_stock_level: parseInt(formData.min_stock_level.toString()),
                reorder_point: parseInt(formData.reorder_point.toString()),
                unit: formData.unit,
            };

            // Add optional fields only if they have values
            if (formData.description?.trim()) payload.description = formData.description;
            if (formData.batch_no?.trim()) payload.batch_no = formData.batch_no;
            if (formData.expiry_date?.trim()) payload.expiry_date = formData.expiry_date;
            if (formData.shelf_location?.trim()) payload.shelf_location = formData.shelf_location;
            if (formData.max_stock_level) payload.max_stock_level = parseInt(formData.max_stock_level.toString());

            await api.put(`/products/${currentProduct.product_id}`, payload);
            showToast('success', 'Product updated successfully');
            resetForm();
            fetchProducts();
            return true;
        } catch (error) {
            showToast('error', (error as Error).message);
            return false;
        } finally {
            setActionLoading(false);
        }
    };

    const handleDeleteProduct = async () => {
        if (!currentProduct) return;
        try {
            setActionLoading(true);
            await api.delete(`/products/${currentProduct.product_id}`);
            showToast('success', 'Product deleted successfully');
            setCurrentProduct(null);
            fetchProducts();
            return true;
        } catch (error) {
            showToast('error', (error as Error).message);
            return false;
        } finally {
            setActionLoading(false);
        }
    };

    const handleBulkDelete = async () => {
        try {
            setActionLoading(true);
            await Promise.all(selectedProducts.map((id) => api.delete(`/products/${id}`)));
            showToast('success', `${selectedProducts.length} products deleted`);
            setSelectedProducts([]);
            fetchProducts();
        } catch (error) {
            showToast('error', 'Failed to delete some products');
        } finally {
            setActionLoading(false);
        }
    };

    const handleExportCSV = () => {
        const headers = ['Name', 'SKU', 'Category', 'Cost Price', 'Selling Price', 'Current Quantity', 'Unit', 'Status'];
        const rows = filteredAndSortedProducts.map((p) => [
            p.name,
            p.sku,
            p.category_name || '',
            p.cost_price,
            p.selling_price,
            p.current_quantity || 0,
            p.unit,
            p.is_active ? 'Active' : 'Inactive',
        ]);
        const csvContent = [
            headers.join(','),
            ...rows.map((row) => row.map((cell) => `"${cell}"`).join(',')),
        ].join('\n');
        const blob = new Blob([csvContent], { type: 'text/csv' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `products_${new Date().toISOString().split('T')[0]}.csv`;
        a.click();
        URL.revokeObjectURL(url);
        showToast('success', 'Products exported successfully');
    };

    const resetForm = () => {
        setFormData(createEmptyFormData());
        setFormErrors({});
        setCurrentProduct(null);
    };

    const openEditModal = (product: Product) => {
        setCurrentProduct(product);
        setFormData({
            name: product.name,
            description: product.description || '',
            sku: product.sku,
            category_id: product.category_id,
            batch_no: product.batch_no || '',
            expiry_date: product.expiry_date?.split('T')[0] || '',
            cost_price: product.cost_price,
            selling_price: product.selling_price,
            current_quantity: product.current_quantity,
            min_stock_level: product.min_stock_level,
            max_stock_level: product.max_stock_level || '',
            reorder_point: product.reorder_point,
            unit: product.unit,
            shelf_location: product.shelf_location || '',
        });
    };

    const handleSelectAll = () => {
        if (selectedProducts.length === paginatedProducts.length) {
            setSelectedProducts([]);
        } else {
            setSelectedProducts(paginatedProducts.map((p) => p.product_id));
        }
    };

    const toggleSelection = (id: number) => {
        setSelectedProducts((prev) =>
            prev.includes(id) ? prev.filter((pid) => pid !== id) : [...prev, id]
        );
    };

    const handleSearchChange = (value: string) => {
        setSearchTerm(value);
        setCurrentPage(1);
    };

    const handleCategoryFilterChange = (value: string) => {
        setCategoryFilter(value);
        setCurrentPage(1);
    };

    const handleStatusFilterChange = (value: string) => {
        setStatusFilter(value);
        setCurrentPage(1);
    };

    return {
        // State
        products: paginatedProducts,
        allProducts: filteredAndSortedProducts,
        categories,
        loading,
        actionLoading,
        currentPage,
        totalPages,
        itemsPerPage,
        searchTerm,
        categoryFilter,
        statusFilter,
        sortField,
        sortDirection,
        selectedProducts,
        currentProduct,
        formData,
        formErrors,

        // Actions
        setCurrentProduct,
        setFormData,
        setCurrentPage,
        handleSort,
        handleAddProduct,
        handleEditProduct,
        handleDeleteProduct,
        handleBulkDelete,
        handleExportCSV,
        resetForm,
        openEditModal,
        handleSelectAll,
        toggleSelection,
        handleSearchChange,
        handleCategoryFilterChange,
        handleStatusFilterChange,
    };
};
