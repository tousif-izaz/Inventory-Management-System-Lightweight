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
}

export const saleService = new SaleService();
