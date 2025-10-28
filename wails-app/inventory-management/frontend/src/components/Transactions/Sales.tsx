import { useState, useEffect } from 'react';
import { Modal } from '../UI/Modal';
import { ConfirmDialog } from '../UI/ConfirmDialog';
import { PlusIcon, EyeIcon, ArrowPathIcon, XMarkIcon, FunnelIcon } from '@heroicons/react/24/outline';
import { Button } from '../UI/Button';
import { Input } from '../UI/Input';
import { Select } from '../UI/Select';
import { useToast } from '../UI/Toast';
import { saleService } from '../../services/saleService';
import { Sale, SaleWithItems } from '../../types/sale';
import { SaleForm } from './SaleForm';
import { SalesList } from './SalesList';
import { SaleDetails } from './SaleDetails';

export const Sales = () => {
    const { showToast } = useToast();
    const [sales, setSales] = useState<Sale[]>([]);
    const [currentSale, setCurrentSale] = useState<SaleWithItems | null>(null);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState(false);

    const [showAddModal, setShowAddModal] = useState(false);
    const [showDetailsModal, setShowDetailsModal] = useState(false);
    const [showRefundDialog, setShowRefundDialog] = useState(false);

    // Filter states
    const [fromDate, setFromDate] = useState<string>('');
    const [toDate, setToDate] = useState<string>('');
    const [paymentStatus, setPaymentStatus] = useState<string>('');

    const fetchSales = async () => {
        try {
            setLoading(true);

            // Build filters object
            const filters: any = {};
            if (fromDate) {
                filters.from_date = new Date(fromDate).toISOString();
            }
            if (toDate) {
                // Set to end of day for to_date
                const endOfDay = new Date(toDate);
                endOfDay.setHours(23, 59, 59, 999);
                filters.to_date = endOfDay.toISOString();
            }
            if (paymentStatus) {
                filters.payment_status = paymentStatus;
            }

            const data = await saleService.getAllSales(filters);
            setSales(data);
        } catch (error) {
            console.error('Failed to fetch sales:', error);
            const errorMessage = error instanceof Error ? error.message : 'Failed to fetch sales';
            showToast('error', `Failed to load sales: ${errorMessage}`);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchSales();
    }, [fromDate, toDate, paymentStatus]);

    const handleClearFilters = () => {
        setFromDate('');
        setToDate('');
        setPaymentStatus('');
    };

    const handleOpenAddModal = () => setShowAddModal(true);

    const handleCloseAddModal = () => {
        setShowAddModal(false);
    };

    const handleOpenDetailsModal = async (sale: Sale) => {
        try {
            setActionLoading(true);
            const saleWithItems = await saleService.getSaleById(sale.sale_id);
            setCurrentSale(saleWithItems);
            setShowDetailsModal(true);
        } catch (error) {
            console.error('Failed to fetch sale details:', error);
            const errorMessage = error instanceof Error ? error.message : 'Failed to fetch sale details';
            showToast('error', `Failed to load sale details: ${errorMessage}`);
        } finally {
            setActionLoading(false);
        }
    };

    const handleCloseDetailsModal = () => {
        setShowDetailsModal(false);
        setCurrentSale(null);
    };

    const handleOpenRefundDialog = (sale: Sale | SaleWithItems) => {
        if ('items' in sale) {
            setCurrentSale(sale);
        }
        setShowRefundDialog(true);
    };

    const handleCloseRefundDialog = () => {
        setShowRefundDialog(false);
    };

    const handleSubmitSale = async () => {
        await fetchSales();
        handleCloseAddModal();
        showToast('success', 'Sale recorded successfully');
    };

    const handleUpdatePaymentStatus = async (
        saleId: number,
        status: 'pending' | 'partial' | 'paid' | 'refunded'
    ) => {
        try {
            setActionLoading(true);
            await saleService.updatePaymentStatus(saleId, status);
            showToast('success', 'Payment status updated successfully');
            await fetchSales();
        } catch (error) {
            console.error('Failed to update payment status:', error);
            const errorMessage = error instanceof Error ? error.message : 'Failed to update payment status';
            showToast('error', `Payment status update failed: ${errorMessage}`);
        } finally {
            setActionLoading(false);
        }
    };

    const handleProcessRefund = async () => {
        if (!currentSale) return;

        try {
            setActionLoading(true);
            await saleService.processRefund(currentSale.sale_id);
            showToast('success', 'Refund processed successfully and inventory restored');
            handleCloseRefundDialog();
            setCurrentSale(null);
            await fetchSales();
        } catch (error) {
            console.error('Failed to process refund:', error);
            const errorMessage = error instanceof Error ? error.message : 'Failed to process refund';
            showToast('error', `Refund failed: ${errorMessage}`);
        } finally {
            setActionLoading(false);
        }
    };

    return (
        <div>
            <div className="mb-6 flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-semibold text-gray-900">Sales</h1>
                    <p className="mt-1 text-sm text-gray-500">
                        Manage sales transactions and receipts ({sales.length} sales)
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="secondary" onClick={fetchSales} disabled={loading}>
                        <ArrowPathIcon className={`h-5 w-5 ${loading ? 'animate-spin' : ''}`} />
                        Refresh
                    </Button>
                    <Button onClick={handleOpenAddModal}>
                        <PlusIcon className="h-5 w-5 mr-2" />
                        New Sale
                    </Button>
                </div>
            </div>

            {/* Filters Section */}
            <div className="bg-white shadow rounded-lg p-4 mb-6">
                <div className="flex items-center gap-2 mb-4">
                    <FunnelIcon className="h-5 w-5 text-gray-600" />
                    <h2 className="text-lg font-medium text-gray-900">Filters</h2>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <Input
                        label="From Date"
                        type="date"
                        value={fromDate}
                        onChange={(e) => setFromDate(e.target.value)}
                    />
                    <Input
                        label="To Date"
                        type="date"
                        value={toDate}
                        onChange={(e) => setToDate(e.target.value)}
                    />
                    <Select
                        label="Payment Status"
                        value={paymentStatus}
                        onChange={(e) => setPaymentStatus(e.target.value)}
                        options={[
                            { value: '', label: 'All Statuses' },
                            { value: 'pending', label: 'Pending' },
                            { value: 'partial', label: 'Partial' },
                            { value: 'paid', label: 'Paid' },
                            { value: 'refunded', label: 'Refunded' },
                        ]}
                    />
                    <div className="flex items-end">
                        <Button
                            variant="secondary"
                            onClick={handleClearFilters}
                            disabled={!fromDate && !toDate && !paymentStatus}
                            className="w-full"
                        >
                            <XMarkIcon className="h-5 w-5 mr-2" />
                            Clear Filters
                        </Button>
                    </div>
                </div>
            </div>

            {loading ? (
                <div className="bg-white shadow rounded-lg p-8 text-center">
                    <div className="inline-block animate-spin rounded-full h-8 w-8 border-4 border-gray-300 border-t-black" />
                    <p className="mt-2 text-gray-600">Loading sales...</p>
                </div>
            ) : sales.length === 0 ? (
                <div className="bg-white shadow rounded-lg p-12 text-center">
                    <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gray-100 mb-4">
                        <EyeIcon className="h-8 w-8 text-gray-400" />
                    </div>
                    <h3 className="text-lg font-medium text-gray-900 mb-2">No sales recorded</h3>
                    <p className="text-gray-500 mb-6">Get started by recording your first sale</p>
                    <Button onClick={handleOpenAddModal}>
                        <PlusIcon className="h-5 w-5 mr-2" />
                        Record Your First Sale
                    </Button>
                </div>
            ) : (
                <SalesList
                    sales={sales}
                    onView={handleOpenDetailsModal}
                    onUpdatePaymentStatus={handleUpdatePaymentStatus}
                    onRefund={handleOpenRefundDialog}
                    loading={actionLoading}
                />
            )}

            <Modal
                isOpen={showAddModal}
                onClose={handleCloseAddModal}
                title="Record New Sale"
                size="xl"
            >
                <div className="p-6">
                    <SaleForm onSuccess={handleSubmitSale} onCancel={handleCloseAddModal} />
                </div>
            </Modal>

            {currentSale && (
                <Modal
                    isOpen={showDetailsModal}
                    onClose={handleCloseDetailsModal}
                    title="Sale Details"
                    size="lg"
                >
                    <div className="p-6">
                        <SaleDetails
                            sale={currentSale}
                            onUpdatePaymentStatus={handleUpdatePaymentStatus}
                            onRefund={handleOpenRefundDialog}
                        />
                    </div>
                </Modal>
            )}

            {currentSale && (
                <ConfirmDialog
                    isOpen={showRefundDialog}
                    onClose={handleCloseRefundDialog}
                    onConfirm={handleProcessRefund}
                    title="Process Refund"
                    message={`Are you sure you want to refund sale ${currentSale.receipt_no}? This will restore the inventory for all items in this sale.`}
                    confirmText="Process Refund"
                    isLoading={actionLoading}
                />
            )}
        </div>
    );
};
