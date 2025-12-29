package dto

import "time"

// ========== Sales Report DTOs ==========

// SalesReportRequest contains filters for sales reports
type SalesReportRequest struct {
	StartDate     time.Time  `json:"start_date"`
	EndDate       time.Time  `json:"end_date"`
	CustomerID    *int64     `json:"customer_id,omitempty"`
	ProductID     *int64     `json:"product_id,omitempty"`
	CategoryID    *int64     `json:"category_id,omitempty"`
	PaymentStatus *string    `json:"payment_status,omitempty"`
	PaymentMethod *string    `json:"payment_method,omitempty"`
	UserID        *int64     `json:"user_id,omitempty"`
}

// SalesReportDTO represents a comprehensive sales report
type SalesReportDTO struct {
	Period             string             `json:"period"`
	StartDate          time.Time          `json:"start_date"`
	EndDate            time.Time          `json:"end_date"`
	TotalRevenue       float64            `json:"total_revenue"`
	GrossSales         float64            `json:"gross_sales"`
	TaxCollected       float64            `json:"tax_collected"`
	DiscountsGiven     float64            `json:"discounts_given"`
	TransactionCount   int                `json:"transaction_count"`
	AvgTransaction     float64            `json:"avg_transaction"`
	ItemsSold          int64              `json:"items_sold"`
	UniqueCustomers    int                `json:"unique_customers"`
	RevenuePerCustomer float64            `json:"revenue_per_customer"`
	DailyBreakdown     []DailySalesDTO    `json:"daily_breakdown,omitempty"`
	TopProducts        []ProductSalesDTO  `json:"top_products,omitempty"`
	PaymentMethods     []PaymentMethodDTO `json:"payment_methods,omitempty"`
}

// DailySalesDTO represents daily sales summary
type DailySalesDTO struct {
	Date             time.Time `json:"date"`
	Revenue          float64   `json:"revenue"`
	TransactionCount int       `json:"transaction_count"`
	ItemsSold        int64     `json:"items_sold"`
}

// ProductSalesDTO represents product sales performance
type ProductSalesDTO struct {
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	SKU          string  `json:"sku"`
	CategoryName string  `json:"category_name"`
	QuantitySold int64   `json:"quantity_sold"`
	Revenue      float64 `json:"revenue"`
	Percentage   float64 `json:"percentage"`
}

// PaymentMethodDTO represents payment method breakdown
type PaymentMethodDTO struct {
	Method           string  `json:"method"`
	TransactionCount int     `json:"transaction_count"`
	TotalAmount      float64 `json:"total_amount"`
	Percentage       float64 `json:"percentage"`
}

// RecentSaleDTO represents a recent sale transaction
type RecentSaleDTO struct {
	SaleID        int64      `json:"sale_id"`
	ReceiptNo     string     `json:"receipt_no"`
	SaleDate      time.Time  `json:"sale_date"`
	CustomerID    *int64     `json:"customer_id,omitempty"`
	CustomerName  *string    `json:"customer_name,omitempty"`
	PaymentMethod string     `json:"payment_method"`
	PaymentStatus string     `json:"payment_status"`
	TotalAmount   float64    `json:"total_amount"`
	NetAmount     float64    `json:"net_amount"`
	ItemCount     int        `json:"item_count"`
}

// ========== Product Performance DTOs ==========

// ProductPerformanceRequest contains filters for product analytics
type ProductPerformanceRequest struct {
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	CategoryID *int64    `json:"category_id,omitempty"`
	SupplierID *int64    `json:"supplier_id,omitempty"`
	Limit      int       `json:"limit"` // For top/bottom N products
}

// ProductPerformanceDTO represents detailed product performance
type ProductPerformanceDTO struct {
	ProductID      int64   `json:"product_id"`
	ProductName    string  `json:"product_name"`
	SKU            string  `json:"sku"`
	CategoryName   string  `json:"category_name"`
	QuantitySold   int64   `json:"quantity_sold"`
	TotalRevenue   float64 `json:"total_revenue"`
	AvgUnitPrice   float64 `json:"avg_unit_price"`
	TotalCost      float64 `json:"total_cost"`
	GrossProfit    float64 `json:"gross_profit"`
	ProfitMargin   float64 `json:"profit_margin"`
	AvgDiscount    float64 `json:"avg_discount"`
	GrowthRate     float64 `json:"growth_rate,omitempty"`
	CurrentStock   int64   `json:"current_stock"`
	Rank           int     `json:"rank"`
}

// CategoryPerformanceDTO represents category sales performance
type CategoryPerformanceDTO struct {
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	ProductCount int     `json:"product_count"`
	QuantitySold int64   `json:"quantity_sold"`
	TotalRevenue float64 `json:"total_revenue"`
	AvgPrice     float64 `json:"avg_price"`
	Percentage   float64 `json:"percentage"`
	GrowthRate   float64 `json:"growth_rate,omitempty"`
}

// ========== ABC Analysis DTOs ==========

// ABCAnalysisRequest contains parameters for ABC analysis
type ABCAnalysisRequest struct {
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	CategoryID *int64    `json:"category_id,omitempty"`
	SupplierID *int64    `json:"supplier_id,omitempty"`
}

// ABCAnalysisDTO represents the complete ABC analysis
type ABCAnalysisDTO struct {
	AnalysisPeriod string         `json:"analysis_period"`
	StartDate      time.Time      `json:"start_date"`
	EndDate        time.Time      `json:"end_date"`
	TotalRevenue   float64        `json:"total_revenue"`
	ClassA         ABCClassDTO    `json:"class_a"`
	ClassB         ABCClassDTO    `json:"class_b"`
	ClassC         ABCClassDTO    `json:"class_c"`
	Products       []ABCProductDTO `json:"products"`
}

// ABCClassDTO represents a single ABC class summary
type ABCClassDTO struct {
	ProductCount   int     `json:"product_count"`
	TotalRevenue   float64 `json:"total_revenue"`
	RevenuePercent float64 `json:"revenue_percent"`
	Recommendation string  `json:"recommendation"`
}

// ABCProductDTO represents a product with ABC classification
type ABCProductDTO struct {
	ProductID         int64   `json:"product_id"`
	ProductName       string  `json:"product_name"`
	SKU               string  `json:"sku"`
	CategoryName      string  `json:"category_name"`
	TotalRevenue      float64 `json:"total_revenue"`
	QuantitySold      int64   `json:"quantity_sold"`
	RevenuePercent    float64 `json:"revenue_percent"`
	CumulativePercent float64 `json:"cumulative_percent"`
	Classification    string  `json:"classification"` // "A", "B", or "C"
	CurrentStock      int64   `json:"current_stock"`
	Recommendation    string  `json:"recommendation"`
}

// ========== Inventory Turnover DTOs ==========

// InventoryTurnoverDTO represents inventory turnover metrics
type InventoryTurnoverDTO struct {
	Period              string                  `json:"period"`
	TurnoverRate        float64                 `json:"turnover_rate"`
	DaysSalesInventory  float64                 `json:"days_sales_inventory"`
	StockToSalesRatio   float64                 `json:"stock_to_sales_ratio"`
	AvgInventoryValue   float64                 `json:"avg_inventory_value"`
	COGS                float64                 `json:"cogs"` // Cost of Goods Sold
	FastMovingCount     int                     `json:"fast_moving_count"`
	MediumMovingCount   int                     `json:"medium_moving_count"`
	SlowMovingCount     int                     `json:"slow_moving_count"`
	DeadStockCount      int                     `json:"dead_stock_count"`
	ProductMovement     []ProductMovementDTO    `json:"product_movement,omitempty"`
}

// ProductMovementDTO represents individual product movement analysis
type ProductMovementDTO struct {
	ProductID           int64      `json:"product_id"`
	ProductName         string     `json:"product_name"`
	SKU                 string     `json:"sku"`
	CategoryName        string     `json:"category_name"`
	CurrentStock        int64      `json:"current_stock"`
	StockValue          float64    `json:"stock_value"`
	LastSaleDate        *time.Time `json:"last_sale_date,omitempty"`
	DaysSinceLastSale   *int       `json:"days_since_last_sale,omitempty"`
	AvgDaysBetweenSales *float64   `json:"avg_days_between_sales,omitempty"`
	TotalSold30Days     int64      `json:"total_sold_30_days"`
	TotalSold90Days     int64      `json:"total_sold_90_days"`
	DaysOfStockLeft     *float64   `json:"days_of_stock_left,omitempty"`
	MovementClass       string     `json:"movement_class"` // "fast", "medium", "slow", "dead"
	ReorderRecommended  bool       `json:"reorder_recommended"`
}

// StockHealthDTO represents overall stock health snapshot
type StockHealthDTO struct {
	TotalProducts       int                   `json:"total_products"`
	TotalStockValue     float64               `json:"total_stock_value"`
	BelowReorderPoint   []LowStockProductDTO  `json:"below_reorder_point"`
	OverstockItems      []OverstockProductDTO `json:"overstock_items"`
	OutOfStock          []OutOfStockProductDTO `json:"out_of_stock"`
	ExpiringProducts    []ExpiringProductDTO  `json:"expiring_products,omitempty"`
}

// LowStockProductDTO represents products below reorder point
type LowStockProductDTO struct {
	ProductID       int64  `json:"product_id"`
	ProductName     string `json:"product_name"`
	SKU             string `json:"sku"`
	CurrentStock    int64  `json:"current_stock"`
	ReorderPoint    int64  `json:"reorder_point"`
	RecommendedQty  int64  `json:"recommended_qty"`
}

// OverstockProductDTO represents overstocked items
type OverstockProductDTO struct {
	ProductID       int64   `json:"product_id"`
	ProductName     string  `json:"product_name"`
	SKU             string  `json:"sku"`
	CurrentStock    int64   `json:"current_stock"`
	MaxStockLevel   int64   `json:"max_stock_level"`
	DaysOfSupply    float64 `json:"days_of_supply"`
	StockValue      float64 `json:"stock_value"`
}

// OutOfStockProductDTO represents out of stock items
type OutOfStockProductDTO struct {
	ProductID      int64      `json:"product_id"`
	ProductName    string     `json:"product_name"`
	SKU            string     `json:"sku"`
	LastSaleDate   *time.Time `json:"last_sale_date,omitempty"`
	AvgDailySales  float64    `json:"avg_daily_sales"`
	DaysOutOfStock int        `json:"days_out_of_stock"`
}

// ExpiringProductDTO represents products nearing expiry
type ExpiringProductDTO struct {
	ProductID    int64     `json:"product_id"`
	ProductName  string    `json:"product_name"`
	SKU          string    `json:"sku"`
	BatchNo      string    `json:"batch_no"`
	ExpiryDate   time.Time `json:"expiry_date"`
	DaysToExpiry int       `json:"days_to_expiry"`
	Quantity     int64     `json:"quantity"`
	StockValue   float64   `json:"stock_value"`
}

// ========== Trend Analysis DTOs ==========

// TrendAnalysisRequest contains parameters for trend analysis
type TrendAnalysisRequest struct {
	Metric     string    `json:"metric"` // "revenue", "quantity", "transactions"
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Granularity string   `json:"granularity"` // "daily", "weekly", "monthly"
	ProductID  *int64    `json:"product_id,omitempty"`
	CategoryID *int64    `json:"category_id,omitempty"`
}

// TrendDataDTO represents trend data over time
type TrendDataDTO struct {
	Metric      string           `json:"metric"`
	Period      string           `json:"period"`
	Granularity string           `json:"granularity"`
	DataPoints  []TrendPointDTO  `json:"data_points"`
	MovingAvg7  []float64        `json:"moving_avg_7,omitempty"`
	MovingAvg30 []float64        `json:"moving_avg_30,omitempty"`
	GrowthRate  float64          `json:"growth_rate"`
	Trend       string           `json:"trend"` // "increasing", "decreasing", "stable"
}

// TrendPointDTO represents a single data point in a trend
type TrendPointDTO struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
	Label string    `json:"label,omitempty"`
}
