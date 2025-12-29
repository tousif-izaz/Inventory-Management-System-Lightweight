import { api } from './api';
import {
    SalesReportDTO,
    ProductPerformanceDTO,
    CategoryPerformanceDTO,
    ABCAnalysisDTO,
    ABCProductDTO,
    InventoryTurnoverDTO,
    ProductMovementDTO,
    StockHealthDTO,
} from '../types/report';

class ReportService {
    // ========== Sales Reports ==========

    async getDailySalesReport(date: string): Promise<SalesReportDTO> {
        return api.get<SalesReportDTO>(`/reports/sales/daily?date=${date}`);
    }

    async getWeeklySalesReport(year: number, week: number): Promise<SalesReportDTO> {
        return api.get<SalesReportDTO>(`/reports/sales/weekly?year=${year}&week=${week}`);
    }

    async getMonthlySalesReport(year: number, month: number): Promise<SalesReportDTO> {
        return api.get<SalesReportDTO>(`/reports/sales/monthly?year=${year}&month=${month}`);
    }

    async getQuarterlySalesReport(year: number, quarter: number): Promise<SalesReportDTO> {
        return api.get<SalesReportDTO>(`/reports/sales/quarterly?year=${year}&quarter=${quarter}`);
    }

    async getCustomSalesReport(params: {
        start_date: string;
        end_date: string;
        customer_id?: number;
        payment_status?: string;
        payment_method?: string;
    }): Promise<SalesReportDTO> {
        const queryParams = new URLSearchParams();
        queryParams.append('start_date', params.start_date);
        queryParams.append('end_date', params.end_date);

        if (params.customer_id) queryParams.append('customer_id', params.customer_id.toString());
        if (params.payment_status) queryParams.append('payment_status', params.payment_status);
        if (params.payment_method) queryParams.append('payment_method', params.payment_method);

        return api.get<SalesReportDTO>(`/reports/sales/custom?${queryParams.toString()}`);
    }

    async getRecentSales(params?: {
        limit?: number;
        start_date?: string;
        end_date?: string;
        payment_method?: string;
        payment_status?: string;
    }): Promise<any[]> {
        const queryParams = new URLSearchParams();

        if (params?.limit) queryParams.append('limit', params.limit.toString());
        if (params?.start_date) queryParams.append('start_date', params.start_date);
        if (params?.end_date) queryParams.append('end_date', params.end_date);
        if (params?.payment_method) queryParams.append('payment_method', params.payment_method);
        if (params?.payment_status) queryParams.append('payment_status', params.payment_status);

        return api.get<any[]>(`/reports/sales/recent?${queryParams.toString()}`);
    }

    // ========== Product Analytics ==========

    async getBestSellingProducts(params: {
        start_date?: string;
        end_date?: string;
        category_id?: number;
        supplier_id?: number;
        limit?: number;
    }): Promise<ProductPerformanceDTO[]> {
        const queryParams = new URLSearchParams();

        if (params.start_date) queryParams.append('start_date', params.start_date);
        if (params.end_date) queryParams.append('end_date', params.end_date);
        if (params.category_id) queryParams.append('category_id', params.category_id.toString());
        if (params.supplier_id) queryParams.append('supplier_id', params.supplier_id.toString());
        if (params.limit) queryParams.append('limit', params.limit.toString());

        return api.get<ProductPerformanceDTO[]>(
            `/reports/products/best-sellers?${queryParams.toString()}`
        );
    }

    async getWorstPerformingProducts(params: {
        start_date?: string;
        end_date?: string;
        category_id?: number;
        supplier_id?: number;
        limit?: number;
    }): Promise<ProductPerformanceDTO[]> {
        const queryParams = new URLSearchParams();

        if (params.start_date) queryParams.append('start_date', params.start_date);
        if (params.end_date) queryParams.append('end_date', params.end_date);
        if (params.category_id) queryParams.append('category_id', params.category_id.toString());
        if (params.supplier_id) queryParams.append('supplier_id', params.supplier_id.toString());
        if (params.limit) queryParams.append('limit', params.limit.toString());

        return api.get<ProductPerformanceDTO[]>(
            `/reports/products/worst-performers?${queryParams.toString()}`
        );
    }

    async getCategoryPerformance(params: {
        start_date?: string;
        end_date?: string;
    }): Promise<CategoryPerformanceDTO[]> {
        const queryParams = new URLSearchParams();

        if (params.start_date) queryParams.append('start_date', params.start_date);
        if (params.end_date) queryParams.append('end_date', params.end_date);

        return api.get<CategoryPerformanceDTO[]>(
            `/reports/categories/performance?${queryParams.toString()}`
        );
    }

    // ========== ABC Analysis ==========

    async getABCAnalysis(params: {
        start_date?: string;
        end_date?: string;
        category_id?: number;
        supplier_id?: number;
    }): Promise<ABCAnalysisDTO> {
        const queryParams = new URLSearchParams();

        if (params.start_date) queryParams.append('start_date', params.start_date);
        if (params.end_date) queryParams.append('end_date', params.end_date);
        if (params.category_id) queryParams.append('category_id', params.category_id.toString());
        if (params.supplier_id) queryParams.append('supplier_id', params.supplier_id.toString());

        return api.get<ABCAnalysisDTO>(`/reports/abc-analysis?${queryParams.toString()}`);
    }

    async getABCProductsByClass(
        classType: 'A' | 'B' | 'C',
        params: {
            start_date?: string;
            end_date?: string;
            category_id?: number;
            supplier_id?: number;
        }
    ): Promise<ABCProductDTO[]> {
        const queryParams = new URLSearchParams();
        queryParams.append('class', classType);

        if (params.start_date) queryParams.append('start_date', params.start_date);
        if (params.end_date) queryParams.append('end_date', params.end_date);
        if (params.category_id) queryParams.append('category_id', params.category_id.toString());
        if (params.supplier_id) queryParams.append('supplier_id', params.supplier_id.toString());

        return api.get<ABCProductDTO[]>(`/reports/abc-analysis/products?${queryParams.toString()}`);
    }

    // ========== Inventory Turnover ==========

    async getInventoryTurnover(): Promise<InventoryTurnoverDTO> {
        return api.get<InventoryTurnoverDTO>('/reports/inventory/turnover');
    }

    async getProductMovement(): Promise<ProductMovementDTO[]> {
        return api.get<ProductMovementDTO[]>('/reports/inventory/movement');
    }

    async getStockHealth(): Promise<StockHealthDTO> {
        return api.get<StockHealthDTO>('/reports/inventory/stock-health');
    }
}

export const reportService = new ReportService();
