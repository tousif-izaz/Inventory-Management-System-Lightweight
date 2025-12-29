import { api } from './api';
import { Sale, SaleWithItems, SaleFormData, SaleFilters } from '../types/sale';

class SaleService {
    /**
     * Create a new sale
     */
    async createSale(saleData: SaleFormData): Promise<SaleWithItems> {
        return api.post<SaleWithItems>('/sales', saleData);
    }

    /**
     * Get all sales with optional filters
     */
    async getAllSales(filters?: SaleFilters): Promise<Sale[]> {
        const params = new URLSearchParams();

        if (filters?.customer_id) {
            params.append('customer_id', filters.customer_id.toString());
        }
        if (filters?.payment_status) {
            params.append('payment_status', filters.payment_status);
        }
        if (filters?.from_date) {
            params.append('from_date', filters.from_date);
        }
        if (filters?.to_date) {
            params.append('to_date', filters.to_date);
        }

        const queryString = params.toString();
        const endpoint = queryString ? `/sales?${queryString}` : '/sales';

        return api.get<Sale[]>(endpoint);
    }

    /**
     * Get a single sale by ID with all items
     */
    async getSaleById(id: number): Promise<SaleWithItems> {
        return api.get<SaleWithItems>(`/sales/${id}`);
    }

    /**
     * Update payment status of a sale
     */
    async updatePaymentStatus(
        id: number,
        paymentStatus: 'pending' | 'partial' | 'paid' | 'refunded'
    ): Promise<void> {
        return api.put<void>(`/sales/${id}/payment-status`, {
            payment_status: paymentStatus,
        });
    }

    /**
     * Process a refund for a sale
     */
    async processRefund(id: number): Promise<void> {
        return api.post<void>(`/sales/${id}/refund`, {});
    }

    /**
     * Backup sales to CSV (Admin only)
     * Downloads a CSV file with all sales data
     */
    async backupSales(fromDate?: string, toDate?: string): Promise<void> {
        const params = new URLSearchParams();

        if (fromDate) {
            params.append('from_date', fromDate);
        }
        if (toDate) {
            params.append('to_date', toDate);
        }

        const queryString = params.toString();
        const endpoint = queryString ? `/sales/backup?${queryString}` : '/sales/backup';

        // Get auth token
        const token = localStorage.getItem('auth_token');
        const headers: HeadersInit = {
            'Content-Type': 'application/json',
        };
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const response = await fetch(`http://localhost:37285/api${endpoint}`, {
            headers,
            credentials: 'include',
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error_message || 'Failed to backup sales');
        }

        // Download the file
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;

        // Get filename from Content-Disposition header or use default
        const contentDisposition = response.headers.get('Content-Disposition');
        const filenameMatch = contentDisposition?.match(/filename="(.+)"/);
        const filename = filenameMatch ? filenameMatch[1] : `sales_backup_${new Date().toISOString().split('T')[0]}.csv`;

        a.download = filename;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
    }

    /**
     * Delete sales in a date range and download backup (Admin only)
     * Automatically downloads a CSV backup of deleted sales
     */
    async deleteSalesInRange(fromDate?: string, toDate?: string): Promise<number> {
        const params = new URLSearchParams();

        if (fromDate) {
            params.append('from_date', fromDate);
        }
        if (toDate) {
            params.append('to_date', toDate);
        }

        const queryString = params.toString();
        const endpoint = queryString ? `/sales/range?${queryString}` : '/sales/range';

        // Get auth token
        const token = localStorage.getItem('auth_token');
        const headers: HeadersInit = {
            'Content-Type': 'application/json',
        };
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const response = await fetch(`http://localhost:37285/api${endpoint}`, {
            method: 'DELETE',
            headers,
            credentials: 'include',
        });

        if (!response.ok) {
            // Try to parse error message from JSON response
            try {
                const error = await response.json();
                throw new Error(error.error_message || 'Failed to delete sales');
            } catch (e) {
                // If JSON parsing fails, throw generic error
                throw new Error(`Failed to delete sales: ${response.status} ${response.statusText}`);
            }
        }

        // Get deleted count from header
        const deletedCount = parseInt(response.headers.get('X-Deleted-Count') || '0', 10);

        // Download the backup file
        const blob = await response.blob();
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;

        // Get filename from Content-Disposition header or use default
        const contentDisposition = response.headers.get('Content-Disposition');
        const filenameMatch = contentDisposition?.match(/filename="(.+)"/);
        const filename = filenameMatch ? filenameMatch[1] : `sales_deleted_backup_${new Date().toISOString().split('T')[0]}.csv`;

        a.download = filename;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);

        return deletedCount;
    }
}

export const saleService = new SaleService();
