export interface Product {
    product_id: number;
    name: string;
    description?: string;
    sku: string;
    category_id: number;
    category_name?: string;
    batch_no?: string;
    expiry_date?: string;
    cost_price: number;
    selling_price: number;
    current_quantity: number;
    min_stock_level: number;
    max_stock_level?: number;
    reorder_point: number;
    unit: string;
    shelf_location?: string;
    is_active: boolean;
    created_at: string;
    updated_at?: string;
}

export interface ProductFormData {
    name: string;
    description?: string;
    sku: string;
    category_id: number | string;
    batch_no?: string;
    expiry_date?: string;
    cost_price: number | string;
    selling_price: number | string;
    current_quantity: number | string;
    min_stock_level: number | string;
    max_stock_level?: number | string;
    reorder_point: number | string;
    unit: string;
    shelf_location?: string;
}

export interface Category {
    category_id: number;
    name: string;
    description?: string;
    parent_category_id?: number;
    created_at: string;
}
