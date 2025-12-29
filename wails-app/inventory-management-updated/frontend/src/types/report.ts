// Report Types

export interface SalesReportRequest {
    start_date: string;
    end_date: string;
    customer_id?: number;
    product_id?: number;
    category_id?: number;
    payment_status?: string;
    payment_method?: string;
    user_id?: number;
}

export interface DailySalesDTO {
    date: string;
    revenue: number;
    transaction_count: number;
    items_sold: number;
}

export interface ProductSalesDTO {
    product_id: number;
    product_name: string;
    sku: string;
    category_name: string;
    quantity_sold: number;
    revenue: number;
    percentage: number;
}

export interface PaymentMethodDTO {
    method: string;
    transaction_count: number;
    total_amount: number;
    percentage: number;
}

export interface SalesReportDTO {
    period: string;
    start_date: string;
    end_date: string;
    total_revenue: number;
    gross_sales: number;
    tax_collected: number;
    discounts_given: number;
    transaction_count: number;
    avg_transaction: number;
    items_sold: number;
    unique_customers: number;
    revenue_per_customer: number;
    daily_breakdown?: DailySalesDTO[];
    top_products?: ProductSalesDTO[];
    payment_methods?: PaymentMethodDTO[];
}

// Product Performance Types

export interface ProductPerformanceRequest {
    start_date: string;
    end_date: string;
    category_id?: number;
    supplier_id?: number;
    limit?: number;
}

export interface ProductPerformanceDTO {
    product_id: number;
    product_name: string;
    sku: string;
    category_name: string;
    quantity_sold: number;
    total_revenue: number;
    avg_unit_price: number;
    total_cost: number;
    gross_profit: number;
    profit_margin: number;
    avg_discount: number;
    growth_rate?: number;
    current_stock: number;
    rank: number;
}

export interface CategoryPerformanceDTO {
    category_id: number;
    category_name: string;
    product_count: number;
    quantity_sold: number;
    total_revenue: number;
    avg_price: number;
    percentage: number;
    growth_rate?: number;
}

// ABC Analysis Types

export interface ABCAnalysisRequest {
    start_date: string;
    end_date: string;
    category_id?: number;
    supplier_id?: number;
}

export interface ABCClassDTO {
    product_count: number;
    total_revenue: number;
    revenue_percent: number;
    recommendation: string;
}

export interface ABCProductDTO {
    product_id: number;
    product_name: string;
    sku: string;
    category_name: string;
    total_revenue: number;
    quantity_sold: number;
    revenue_percent: number;
    cumulative_percent: number;
    classification: 'A' | 'B' | 'C';
    current_stock: number;
    recommendation: string;
}

export interface ABCAnalysisDTO {
    analysis_period: string;
    start_date: string;
    end_date: string;
    total_revenue: number;
    class_a: ABCClassDTO;
    class_b: ABCClassDTO;
    class_c: ABCClassDTO;
    products: ABCProductDTO[];
}

// Inventory Turnover Types

export interface ProductMovementDTO {
    product_id: number;
    product_name: string;
    sku: string;
    category_name: string;
    current_stock: number;
    stock_value: number;
    last_sale_date?: string;
    days_since_last_sale?: number;
    avg_days_between_sales?: number;
    total_sold_30_days: number;
    total_sold_90_days: number;
    days_of_stock_left?: number;
    movement_class: 'fast' | 'medium' | 'slow' | 'dead';
    reorder_recommended: boolean;
}

export interface InventoryTurnoverDTO {
    period: string;
    turnover_rate: number;
    days_sales_inventory: number;
    stock_to_sales_ratio: number;
    avg_inventory_value: number;
    cogs: number;
    fast_moving_count: number;
    medium_moving_count: number;
    slow_moving_count: number;
    dead_stock_count: number;
    product_movement?: ProductMovementDTO[];
}

export interface LowStockProductDTO {
    product_id: number;
    product_name: string;
    sku: string;
    current_stock: number;
    reorder_point: number;
    recommended_qty: number;
}

export interface OverstockProductDTO {
    product_id: number;
    product_name: string;
    sku: string;
    current_stock: number;
    max_stock_level: number;
    days_of_supply: number;
    stock_value: number;
}

export interface OutOfStockProductDTO {
    product_id: number;
    product_name: string;
    sku: string;
    last_sale_date?: string;
    avg_daily_sales: number;
    days_out_of_stock: number;
}

export interface ExpiringProductDTO {
    product_id: number;
    product_name: string;
    sku: string;
    batch_no: string;
    expiry_date: string;
    days_to_expiry: number;
    quantity: number;
    stock_value: number;
}

export interface StockHealthDTO {
    total_products: number;
    total_stock_value: number;
    below_reorder_point: LowStockProductDTO[];
    overstock_items: OverstockProductDTO[];
    out_of_stock: OutOfStockProductDTO[];
    expiring_products?: ExpiringProductDTO[];
}

// Trend Analysis Types

export interface TrendAnalysisRequest {
    metric: 'revenue' | 'quantity' | 'transactions';
    start_date: string;
    end_date: string;
    granularity: 'daily' | 'weekly' | 'monthly';
    product_id?: number;
    category_id?: number;
}

export interface TrendPointDTO {
    date: string;
    value: number;
    label?: string;
}

export interface TrendDataDTO {
    metric: string;
    period: string;
    granularity: string;
    data_points: TrendPointDTO[];
    moving_avg_7?: number[];
    moving_avg_30?: number[];
    growth_rate: number;
    trend: 'increasing' | 'decreasing' | 'stable';
}
