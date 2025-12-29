export interface Sale {
    sale_id: number;
    sale_date: string;
    receipt_no: string;
    customer_id?: number;
    total_amount: number;
    tax_amount: number;
    discount_amount: number;
    net_amount: number;
    payment_status: 'pending' | 'partial' | 'paid' | 'refunded';
    payment_method?: 'cash' | 'card' | 'mobile' | 'bank_transfer' | 'credit' | 'square_terminal';
    notes?: string;
    sold_by?: number;
    created_at: string;
}

export interface SaleItem {
    sale_item_id: number;
    sale_id: number;
    product_id: number;
    product_name?: string;
    product_sku?: string;
    quantity: number;
    unit_price: number;
    tax_rate: number;
    discount_percent: number;
    line_total: number;
}

export interface SaleWithItems extends Sale {
    items: SaleItem[];
}

export interface SaleItemFormData {
    product_id: number | string;
    quantity: number | string;
    unit_price: number | string;
    tax_rate: number | string;
    discount_percent: number | string;
}

export interface SaleFormData {
    sale_date?: string;
    receipt_no: string;
    customer_id?: number | string;
    payment_status: 'pending' | 'partial' | 'paid' | 'refunded';
    payment_method?: 'cash' | 'card' | 'mobile' | 'bank_transfer' | 'credit' | 'square_terminal';
    notes?: string;
    sold_by?: number | string;
    items: SaleItemFormData[];
}

export interface SaleFilters {
    customer_id?: number;
    payment_status?: string;
    from_date?: string;
    to_date?: string;
}
