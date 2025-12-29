import { useState } from 'react';
import { Modal } from '../UI/Modal';
import { ConfirmDialog } from '../UI/ConfirmDialog';
import { FunnelIcon, PlusIcon } from '@heroicons/react/24/outline';
import { Button } from '../UI/Button';
import { ProductTable } from './ProductTable';
import { ProductForm } from './ProductForm';
import { ProductDetails } from './ProductDetails';
import { ProductSearchBar } from './ProductSearchBar';
import { ProductFilters } from './ProductFilters';
import { ProductActions } from './ProductActions';
import { useProductLogic } from './useProductLogic';
import { useAuth } from '../../hooks/useAuth';

export const Products = () => {
    const { canAccessInventoryManagement } = useAuth();
    const {
        products,
        allProducts,
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
        currentProduct,
        formData,
        formErrors,
        setCurrentProduct,
        setFormData,
        setCurrentPage,
        handleSort,
        handleAddProduct,
        handleEditProduct,
        handleDeleteProduct,
        handleExportCSV,
        resetForm,
        openEditModal,
        handleSearchChange,
        handleCategoryFilterChange,
        handleStatusFilterChange,
    } = useProductLogic();

    const [showAddModal, setShowAddModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [showDetailsModal, setShowDetailsModal] = useState(false);
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);

    const handleOpenAddModal = () => setShowAddModal(true);

    const handleCloseAddModal = () => {
        setShowAddModal(false);
        resetForm();
    };

    const handleOpenEditModal = (product: typeof currentProduct) => {
        if (product) {
            openEditModal(product);
            setShowEditModal(true);
        }
    };

    const handleCloseEditModal = () => {
        setShowEditModal(false);
        resetForm();
    };

    const handleOpenDetailsModal = (product: typeof currentProduct) => {
        if (product) {
            setCurrentProduct(product);
            setShowDetailsModal(true);
        }
    };

    const handleCloseDetailsModal = () => {
        setShowDetailsModal(false);
        setCurrentProduct(null);
    };

    const handleOpenDeleteDialog = (product: typeof currentProduct) => {
        if (product) {
            setCurrentProduct(product);
            setShowDeleteDialog(true);
        }
    };

    const handleCloseDeleteDialog = () => {
        setShowDeleteDialog(false);
        setCurrentProduct(null);
    };

    const handleSubmitAdd = async () => {
        const success = await handleAddProduct();
        if (success) {
            handleCloseAddModal();
        }
    };

    const handleSubmitEdit = async () => {
        const success = await handleEditProduct();
        if (success) {
            handleCloseEditModal();
        }
    };

    const handleConfirmDelete = async () => {
        const success = await handleDeleteProduct();
        if (success) {
            handleCloseDeleteDialog();
        }
    };

    return (
        <div>
            <div className="mb-6">
                <h1 className="text-2xl font-semibold text-gray-900">All Products</h1>
                <p className="mt-1 text-sm text-gray-500">
                    Manage your product inventory ({allProducts.length} products)
                </p>
            </div>

            <div className="bg-white shadow rounded-lg p-4 mb-6 space-y-4">
                <div className="flex flex-col sm:flex-row gap-4">
                    <ProductSearchBar
                        searchTerm={searchTerm}
                        onSearchChange={handleSearchChange}
                    />
                    {canAccessInventoryManagement() && (
                        <ProductActions
                            onAddProduct={handleOpenAddModal}
                            onExport={handleExportCSV}
                        />
                    )}
                </div>

                <ProductFilters
                    categoryFilter={categoryFilter}
                    statusFilter={statusFilter}
                    categories={categories}
                    onCategoryChange={handleCategoryFilterChange}
                    onStatusChange={handleStatusFilterChange}
                />
            </div>

            {loading ? (
                <div className="bg-white shadow rounded-lg p-8 text-center">
                    <div className="inline-block animate-spin rounded-full h-8 w-8 border-4 border-gray-300 border-t-black" />
                    <p className="mt-2 text-gray-600">Loading products...</p>
                </div>
            ) : products.length === 0 ? (
                <div className="bg-white shadow rounded-lg p-12 text-center">
                    <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gray-100 mb-4">
                        <FunnelIcon className="h-8 w-8 text-gray-400" />
                    </div>
                    <h3 className="text-lg font-medium text-gray-900 mb-2">No products found</h3>
                    <p className="text-gray-500 mb-6">
                        {searchTerm || categoryFilter || statusFilter
                            ? 'Try adjusting your filters'
                            : 'Get started by adding your first product'}
                    </p>
                    {!searchTerm && !categoryFilter && !statusFilter && (
                        <Button onClick={handleOpenAddModal}>
                            <PlusIcon className="h-5 w-5 mr-2" />
                            Add Your First Product
                        </Button>
                    )}
                </div>
            ) : (
                <ProductTable
                    products={products}
                    sortField={sortField}
                    sortDirection={sortDirection}
                    onSort={handleSort}
                    onView={handleOpenDetailsModal}
                    onEdit={handleOpenEditModal}
                    onDelete={handleOpenDeleteDialog}
                    currentPage={currentPage}
                    totalPages={totalPages}
                    itemsPerPage={itemsPerPage}
                    totalItems={allProducts.length}
                    onPageChange={setCurrentPage}
                    canEdit={canAccessInventoryManagement()}
                />
            )}

            <Modal isOpen={showAddModal} onClose={handleCloseAddModal} title="Add New Product" size="xl">
                <ProductForm
                    formData={formData}
                    setFormData={setFormData}
                    formErrors={formErrors}
                    categories={categories}
                    onSubmit={handleSubmitAdd}
                    onCancel={handleCloseAddModal}
                    isLoading={actionLoading}
                    submitText="Add Product"
                />
            </Modal>

            <Modal isOpen={showEditModal} onClose={handleCloseEditModal} title="Edit Product" size="xl">
                <ProductForm
                    formData={formData}
                    setFormData={setFormData}
                    formErrors={formErrors}
                    categories={categories}
                    onSubmit={handleSubmitEdit}
                    onCancel={handleCloseEditModal}
                    isLoading={actionLoading}
                    submitText="Update Product"
                />
            </Modal>

            {currentProduct && (
                <Modal isOpen={showDetailsModal} onClose={handleCloseDetailsModal} title="Product Details" size="lg">
                    <ProductDetails
                        product={currentProduct}
                        onEdit={() => {
                            handleCloseDetailsModal();
                            handleOpenEditModal(currentProduct);
                        }}
                    />
                </Modal>
            )}

            {currentProduct && (
                <ConfirmDialog
                    isOpen={showDeleteDialog}
                    onClose={handleCloseDeleteDialog}
                    onConfirm={handleConfirmDelete}
                    title="Delete Product"
                    message={`Are you sure you want to permanently delete "${currentProduct.name}" - ${currentProduct.category_name || 'No Category'}? This action cannot be undone.`}
                    confirmText="Delete"
                    isLoading={actionLoading}
                />
            )}
        </div>
    );
};
