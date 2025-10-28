import { api } from './api';
import { Product } from '../types/product';
import { Sale } from '../types/sale';

export interface DashboardStats {
    totalProducts: number;
    lowStockCount: number;
    totalInventoryValue: number;
    todaysSales: number;
}

class DashboardService {
    /**
     * Get dashboard statistics
     */
    async getDashboardStats(): Promise<DashboardStats> {
        try {
            // Fetch all products to calculate stats
            const products = await api.get<Product[]>('/products');

            // Calculate total products (active only)
            const totalProducts = products.filter(p => p.is_active).length;

            // Calculate low stock count (current_quantity <= reorder_point)
            const lowStockCount = products.filter(
                p => p.is_active && p.current_quantity <= p.reorder_point
            ).length;

            // Calculate total inventory value (cost_price * current_quantity)
            const totalInventoryValue = products
                .filter(p => p.is_active)
                .reduce((sum, p) => sum + (p.cost_price * p.current_quantity), 0);

            // Get today's sales
            const today = new Date();
            today.setHours(0, 0, 0, 0);
            const tomorrow = new Date(today);
            tomorrow.setDate(tomorrow.getDate() + 1);

            const queryParams = new URLSearchParams({
                from_date: today.toISOString(),
                to_date: tomorrow.toISOString(),
            });
            const todaySales = await api.get<Sale[]>(`/sales?${queryParams.toString()}`);

            const todaysSales = todaySales.reduce((sum, sale) => sum + sale.net_amount, 0);

            return {
                totalProducts,
                lowStockCount,
                totalInventoryValue,
                todaysSales,
            };
        } catch (error) {
            console.error('Failed to fetch dashboard stats:', error);
            throw error;
        }
    }

    /**
     * Get low stock products
     */
    async getLowStockProducts(): Promise<Product[]> {
        try {
            const products = await api.get<Product[]>('/products');
            return products
                .filter(p => p.is_active && p.current_quantity <= p.reorder_point)
                .sort((a, b) => a.current_quantity - b.current_quantity)
                .slice(0, 5); // Return top 5 most critical
        } catch (error) {
            console.error('Failed to fetch low stock products:', error);
            throw error;
        }
    }

    /**
     * Get products expiring soon (within next 30 days)
     */
    async getExpiringProducts(days: number = 30): Promise<Product[]> {
        try {
            const products = await api.get<Product[]>('/products');
            const today = new Date();
            const futureDate = new Date();
            futureDate.setDate(futureDate.getDate() + days);

            return products
                .filter(p => {
                    if (!p.is_active || !p.expiry_date) return false;
                    const expiryDate = new Date(p.expiry_date);
                    return expiryDate >= today && expiryDate <= futureDate;
                })
                .sort((a, b) => {
                    const dateA = a.expiry_date ? new Date(a.expiry_date).getTime() : Infinity;
                    const dateB = b.expiry_date ? new Date(b.expiry_date).getTime() : Infinity;
                    return dateA - dateB;
                })
                .slice(0, 5); // Return top 5 soonest to expire
        } catch (error) {
            console.error('Failed to fetch expiring products:', error);
            throw error;
        }
    }

    /**
     * Get recent sales (last 5)
     */
    async getRecentSales(): Promise<Sale[]> {
        try {
            const sales = await api.get<Sale[]>('/sales');
            return sales
                .sort((a, b) => new Date(b.sale_date).getTime() - new Date(a.sale_date).getTime())
                .slice(0, 5);
        } catch (error) {
            console.error('Failed to fetch recent sales:', error);
            throw error;
        }
    }
}

export const dashboardService = new DashboardService();
