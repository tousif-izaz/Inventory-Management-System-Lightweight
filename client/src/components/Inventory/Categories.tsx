import { useState, useEffect } from 'react';
import { PlusIcon, PencilIcon, TrashIcon, ArrowPathIcon } from '@heroicons/react/24/outline';
import { Button } from '../UI/Button';
import { Input } from '../UI/Input';
import { Modal } from '../UI/Modal';
import { ConfirmDialog } from '../UI/ConfirmDialog';
import { useToast } from '../UI/Toast';
import { api } from '../../services/api';
import { Category } from '../../types/product';

interface CategoryFormData {
    name: string;
    description: string;
}

export const Categories = () => {
    const { showToast } = useToast();
    const [categories, setCategories] = useState<Category[]>([]);
    const [loading, setLoading] = useState(true);
    const [showAddModal, setShowAddModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);
    const [currentCategory, setCurrentCategory] = useState<Category | null>(null);
    const [formData, setFormData] = useState<CategoryFormData>({
        name: '',
        description: '',
    });
    const [actionLoading, setActionLoading] = useState(false);

    const fetchCategories = async () => {
        try {
            setLoading(true);
            const data = await api.get<Category[]>('/categories');
            setCategories(data);
        } catch (error) {
            console.error('Failed to fetch categories:', error);
            showToast('error', 'Failed to load categories');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchCategories();
    }, []);

    const handleOpenAddModal = () => {
        setFormData({ name: '', description: '' });
        setShowAddModal(true);
    };

    const handleOpenEditModal = (category: Category) => {
        setCurrentCategory(category);
        setFormData({
            name: category.name,
            description: category.description || '',
        });
        setShowEditModal(true);
    };

    const handleOpenDeleteDialog = (category: Category) => {
        setCurrentCategory(category);
        setShowDeleteDialog(true);
    };

    const handleCloseModals = () => {
        setShowAddModal(false);
        setShowEditModal(false);
        setShowDeleteDialog(false);
        setCurrentCategory(null);
        setFormData({ name: '', description: '' });
    };

    const handleCreateCategory = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!formData.name.trim()) {
            showToast('error', 'Category name is required');
            return;
        }

        try {
            setActionLoading(true);
            await api.post('/categories', {
                name: formData.name.trim(),
                description: formData.description.trim() || undefined,
            });
            showToast('success', 'Category created successfully');
            handleCloseModals();
            fetchCategories();
        } catch (error) {
            console.error('Failed to create category:', error);
            showToast('error', 'Failed to create category');
        } finally {
            setActionLoading(false);
        }
    };

    const handleUpdateCategory = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!currentCategory || !formData.name.trim()) {
            showToast('error', 'Category name is required');
            return;
        }

        try {
            setActionLoading(true);
            await api.put(`/categories/${currentCategory.category_id}`, {
                name: formData.name.trim(),
                description: formData.description.trim() || undefined,
            });
            showToast('success', 'Category updated successfully');
            handleCloseModals();
            fetchCategories();
        } catch (error) {
            console.error('Failed to update category:', error);
            showToast('error', 'Failed to update category');
        } finally {
            setActionLoading(false);
        }
    };

    const handleDeleteCategory = async () => {
        if (!currentCategory) return;

        try {
            setActionLoading(true);
            await api.delete(`/categories/${currentCategory.category_id}`);
            showToast('success', 'Category deleted successfully');
            handleCloseModals();
            fetchCategories();
        } catch (error) {
            console.error('Failed to delete category:', error);
            showToast('error', 'Failed to delete category. It may have associated products.');
        } finally {
            setActionLoading(false);
        }
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
        });
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
                    <h1 className="text-2xl font-semibold text-gray-900">Categories</h1>
                    <p className="mt-1 text-sm text-gray-500">
                        Manage product categories ({categories.length} total)
                    </p>
                </div>
                <div className="flex gap-2">
                    <button
                        onClick={fetchCategories}
                        className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                        disabled={loading}
                    >
                        <ArrowPathIcon className={`h-5 w-5 ${loading ? 'animate-spin' : ''}`} />
                    </button>
                    <Button onClick={handleOpenAddModal}>
                        <PlusIcon className="h-5 w-5 mr-2" />
                        Add Category
                    </Button>
                </div>
            </div>

            {categories.length === 0 ? (
                <div className="bg-white shadow rounded-lg p-12 text-center">
                    <h3 className="text-lg font-medium text-gray-900 mb-2">No Categories</h3>
                    <p className="text-gray-500 mb-6">Get started by creating your first category</p>
                    <Button onClick={handleOpenAddModal}>
                        <PlusIcon className="h-5 w-5 mr-2" />
                        Add Category
                    </Button>
                </div>
            ) : (
                <div className="bg-white shadow rounded-lg overflow-hidden">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Name
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Description
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Created
                                </th>
                                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Actions
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {categories.map((category) => (
                                <tr key={category.category_id} className="hover:bg-gray-50">
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <div className="text-sm font-medium text-gray-900">
                                            {category.name}
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="text-sm text-gray-500">
                                            {category.description || '-'}
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <div className="text-sm text-gray-500">
                                            {formatDate(category.created_at)}
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                        <button
                                            onClick={() => handleOpenEditModal(category)}
                                            className="text-indigo-600 hover:text-indigo-900 mr-4"
                                        >
                                            <PencilIcon className="h-5 w-5 inline" />
                                        </button>
                                        <button
                                            onClick={() => handleOpenDeleteDialog(category)}
                                            className="text-red-600 hover:text-red-900"
                                        >
                                            <TrashIcon className="h-5 w-5 inline" />
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}

            {/* Add Category Modal */}
            <Modal
                isOpen={showAddModal}
                onClose={handleCloseModals}
                title="Add New Category"
                size="md"
            >
                <form onSubmit={handleCreateCategory} className="p-6 space-y-4">
                    <Input
                        label="Category Name"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        required
                        placeholder="e.g., Fragrances"
                    />
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Description (Optional)
                        </label>
                        <textarea
                            value={formData.description}
                            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                            rows={3}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                            placeholder="Brief description of the category"
                        />
                    </div>
                    <div className="flex gap-3 justify-end pt-4">
                        <Button type="button" variant="secondary" onClick={handleCloseModals}>
                            Cancel
                        </Button>
                        <Button type="submit" disabled={actionLoading}>
                            {actionLoading ? 'Creating...' : 'Create Category'}
                        </Button>
                    </div>
                </form>
            </Modal>

            {/* Edit Category Modal */}
            <Modal
                isOpen={showEditModal}
                onClose={handleCloseModals}
                title="Edit Category"
                size="md"
            >
                <form onSubmit={handleUpdateCategory} className="p-6 space-y-4">
                    <Input
                        label="Category Name"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        required
                        placeholder="e.g., Fragrances"
                    />
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Description (Optional)
                        </label>
                        <textarea
                            value={formData.description}
                            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                            rows={3}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                            placeholder="Brief description of the category"
                        />
                    </div>
                    <div className="flex gap-3 justify-end pt-4">
                        <Button type="button" variant="secondary" onClick={handleCloseModals}>
                            Cancel
                        </Button>
                        <Button type="submit" disabled={actionLoading}>
                            {actionLoading ? 'Updating...' : 'Update Category'}
                        </Button>
                    </div>
                </form>
            </Modal>

            {/* Delete Confirmation Dialog */}
            {currentCategory && (
                <ConfirmDialog
                    isOpen={showDeleteDialog}
                    onClose={handleCloseModals}
                    onConfirm={handleDeleteCategory}
                    title="Delete Category"
                    message={`Are you sure you want to delete "${currentCategory.name}"? This action cannot be undone and will fail if there are products in this category.`}
                    confirmText="Delete"
                    isLoading={actionLoading}
                />
            )}
        </div>
    );
};
