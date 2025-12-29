package service

import (
	"database/sql"
	"fmt"
	"ims-intro/pkg/service/dto"
	"log"
	"time"
)

type SalesReportService struct {
	db *sql.DB
}

func NewSalesReportService(db *sql.DB) *SalesReportService {
	return &SalesReportService{db: db}
}

// GetSalesReport generates a comprehensive sales report for a given period
func (s *SalesReportService) GetSalesReport(req *dto.SalesReportRequest) (*dto.SalesReportDTO, error) {
	report := &dto.SalesReportDTO{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Period:    fmt.Sprintf("%s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02")),
	}

	// Build WHERE clause for filters
	whereClause, args := s.buildWhereClause(req)

	// Get overall metrics
	query := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(NetAmount), 0) as total_revenue,
			COALESCE(SUM(TotalAmount), 0) as gross_sales,
			COALESCE(SUM(TaxAmount), 0) as tax_collected,
			COALESCE(SUM(DiscountAmount), 0) as discounts_given,
			COUNT(DISTINCT SaleID) as transaction_count,
			COUNT(DISTINCT CustomerID) as unique_customers
		FROM Sales
		WHERE %s
	`, whereClause)

	var uniqueCustomers sql.NullInt64
	err := s.db.QueryRow(query, args...).Scan(
		&report.TotalRevenue,
		&report.GrossSales,
		&report.TaxCollected,
		&report.DiscountsGiven,
		&report.TransactionCount,
		&uniqueCustomers,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales metrics: %w", err)
	}

	report.UniqueCustomers = int(uniqueCustomers.Int64)

	// Calculate average transaction
	if report.TransactionCount > 0 {
		report.AvgTransaction = report.TotalRevenue / float64(report.TransactionCount)
	}

	// Calculate revenue per customer
	if report.UniqueCustomers > 0 {
		report.RevenuePerCustomer = report.TotalRevenue / float64(report.UniqueCustomers)
	}

	// Get total items sold
	itemsQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(si.Quantity), 0)
		FROM SalesItems si
		INNER JOIN Sales s ON si.SaleID = s.SaleID
		WHERE %s
	`, whereClause)

	err = s.db.QueryRow(itemsQuery, args...).Scan(&report.ItemsSold)
	if err != nil {
		return nil, fmt.Errorf("failed to get items sold: %w", err)
	}

	// Get daily breakdown
	dailyBreakdown, err := s.getDailyBreakdown(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily breakdown: %w", err)
	}
	report.DailyBreakdown = dailyBreakdown

	// Get top products
	topProducts, err := s.getTopProducts(req, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top products: %w", err)
	}
	report.TopProducts = topProducts

	// Get payment method breakdown
	paymentMethods, err := s.getPaymentMethodBreakdown(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	report.PaymentMethods = paymentMethods

	return report, nil
}

// getDailyBreakdown returns daily sales summary
func (s *SalesReportService) getDailyBreakdown(req *dto.SalesReportRequest) ([]dto.DailySalesDTO, error) {
	whereClause, args := s.buildWhereClause(req)

	query := fmt.Sprintf(`
		SELECT
			SUBSTR(SaleDate, 1, 10) as sale_day,
			COALESCE(SUM(NetAmount), 0) as revenue,
			COUNT(DISTINCT SaleID) as transaction_count,
			COALESCE(SUM((SELECT SUM(Quantity) FROM SalesItems WHERE SaleID = Sales.SaleID)), 0) as items_sold
		FROM Sales
		WHERE %s
		GROUP BY SUBSTR(SaleDate, 1, 10)
		ORDER BY sale_day
	`, whereClause)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breakdown []dto.DailySalesDTO
	for rows.Next() {
		var daily dto.DailySalesDTO
		var dateStr sql.NullString

		err := rows.Scan(&dateStr, &daily.Revenue, &daily.TransactionCount, &daily.ItemsSold)
		if err != nil {
			return nil, err
		}

		// Skip rows with NULL dates
		if !dateStr.Valid || dateStr.String == "" {
			continue
		}

		daily.Date, _ = time.Parse("2006-01-02", dateStr.String)
		breakdown = append(breakdown, daily)
	}

	return breakdown, nil
}

// getTopProducts returns top selling products by revenue
func (s *SalesReportService) getTopProducts(req *dto.SalesReportRequest, limit int) ([]dto.ProductSalesDTO, error) {
	whereClause, args := s.buildWhereClause(req)
	args = append(args, limit)

	query := fmt.Sprintf(`
		SELECT
			p.ProductID,
			p.Name,
			p.SKU,
			COALESCE(c.Name, 'Uncategorized') as category_name,
			COALESCE(SUM(si.Quantity), 0) as quantity_sold,
			COALESCE(SUM(si.LineTotal), 0) as revenue
		FROM SalesItems si
		INNER JOIN Sales s ON si.SaleID = s.SaleID
		INNER JOIN Products p ON si.ProductID = p.ProductID
		LEFT JOIN Categories c ON p.CategoryID = c.CategoryID
		WHERE %s
		GROUP BY p.ProductID, p.Name, p.SKU, c.Name
		ORDER BY revenue DESC
		LIMIT ?
	`, whereClause)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []dto.ProductSalesDTO
	var totalRevenue float64

	for rows.Next() {
		var product dto.ProductSalesDTO
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.SKU,
			&product.CategoryName,
			&product.QuantitySold,
			&product.Revenue,
		)
		if err != nil {
			return nil, err
		}
		totalRevenue += product.Revenue
		products = append(products, product)
	}

	// Calculate percentages
	for i := range products {
		if totalRevenue > 0 {
			products[i].Percentage = (products[i].Revenue / totalRevenue) * 100
		}
	}

	return products, nil
}

// getPaymentMethodBreakdown returns sales breakdown by payment method
func (s *SalesReportService) getPaymentMethodBreakdown(req *dto.SalesReportRequest) ([]dto.PaymentMethodDTO, error) {
	whereClause, args := s.buildWhereClause(req)

	query := fmt.Sprintf(`
		SELECT
			COALESCE(PaymentMethod, 'Unknown') as method,
			COUNT(*) as transaction_count,
			COALESCE(SUM(NetAmount), 0) as total_amount
		FROM Sales
		WHERE %s
		GROUP BY PaymentMethod
		ORDER BY total_amount DESC
	`, whereClause)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []dto.PaymentMethodDTO
	var totalAmount float64

	for rows.Next() {
		var method dto.PaymentMethodDTO
		err := rows.Scan(&method.Method, &method.TransactionCount, &method.TotalAmount)
		if err != nil {
			return nil, err
		}
		totalAmount += method.TotalAmount
		methods = append(methods, method)
	}

	// Calculate percentages
	for i := range methods {
		if totalAmount > 0 {
			methods[i].Percentage = (methods[i].TotalAmount / totalAmount) * 100
		}
	}

	return methods, nil
}

// GetDailySalesReport generates a daily sales report
func (s *SalesReportService) GetDailySalesReport(date time.Time) (*dto.SalesReportDTO, error) {
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.Add(24 * time.Hour)

	req := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	return s.GetSalesReport(req)
}

// GetWeeklySalesReport generates a weekly sales report
func (s *SalesReportService) GetWeeklySalesReport(year int, week int) (*dto.SalesReportDTO, error) {
	// Calculate start of week (Monday)
	jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	daysToAdd := (week - 1) * 7
	startDate := jan1.AddDate(0, 0, daysToAdd)

	// Adjust to Monday
	weekday := int(startDate.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday
	}
	startDate = startDate.AddDate(0, 0, -(weekday - 1))

	endDate := startDate.AddDate(0, 0, 7)

	req := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	return s.GetSalesReport(req)
}

// GetMonthlySalesReport generates a monthly sales report
func (s *SalesReportService) GetMonthlySalesReport(year int, month int) (*dto.SalesReportDTO, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	req := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	return s.GetSalesReport(req)
}

// GetQuarterlySalesReport generates a quarterly sales report
func (s *SalesReportService) GetQuarterlySalesReport(year int, quarter int) (*dto.SalesReportDTO, error) {
	if quarter < 1 || quarter > 4 {
		return nil, fmt.Errorf("invalid quarter: %d (must be 1-4)", quarter)
	}

	startMonth := (quarter-1)*3 + 1
	startDate := time.Date(year, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 3, 0)

	req := &dto.SalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	return s.GetSalesReport(req)
}

// GetRecentSales returns the most recent sales transactions
func (s *SalesReportService) GetRecentSales(limit int, filters *dto.SalesReportRequest) ([]dto.RecentSaleDTO, error) {
	log.Printf("🔵 [GetRecentSales Service] Starting - limit: %d", limit)

	whereClause := "1=1"
	var args []interface{}

	// Apply filters if provided
	if filters != nil {
		log.Printf("🔵 [GetRecentSales Service] Applying filters - StartDate: %s, EndDate: %s",
			filters.StartDate.Format("2006-01-02 15:04:05"),
			filters.EndDate.Format("2006-01-02 15:04:05"))
		whereClause, args = s.buildWhereClause(filters)
	}

	args = append(args, limit)

	query := fmt.Sprintf(`
		SELECT
			s.SaleID,
			s.ReceiptNo,
			s.SaleDate,
			s.CustomerID,
			c.Name as customer_name,
			COALESCE(s.PaymentMethod, 'Unknown') as payment_method,
			COALESCE(s.PaymentStatus, 'Unknown') as payment_status,
			s.TotalAmount,
			s.NetAmount,
			(SELECT COUNT(*) FROM SalesItems WHERE SaleID = s.SaleID) as item_count
		FROM Sales s
		LEFT JOIN Customers c ON s.CustomerID = c.CustomerID
		WHERE %s
		ORDER BY s.SaleDate DESC, s.SaleID DESC
		LIMIT ?
	`, whereClause)

	log.Printf("🔵 [GetRecentSales Service] Executing SQL:\n%s", query)
	log.Printf("🔵 [GetRecentSales Service] SQL Args: %v", args)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Printf("❌ [GetRecentSales Service] SQL Error: %v", err)
		return nil, fmt.Errorf("failed to get recent sales: %w", err)
	}
	defer rows.Close()

	var sales []dto.RecentSaleDTO
	for rows.Next() {
		var sale dto.RecentSaleDTO
		var customerID sql.NullInt64
		var customerName sql.NullString

		err := rows.Scan(
			&sale.SaleID,
			&sale.ReceiptNo,
			&sale.SaleDate,
			&customerID,
			&customerName,
			&sale.PaymentMethod,
			&sale.PaymentStatus,
			&sale.TotalAmount,
			&sale.NetAmount,
			&sale.ItemCount,
		)
		if err != nil {
			log.Printf("❌ [GetRecentSales Service] Row scan error: %v", err)
			return nil, err
		}

		if customerID.Valid {
			sale.CustomerID = &customerID.Int64
		}
		if customerName.Valid {
			sale.CustomerName = &customerName.String
		}

		sales = append(sales, sale)
	}

	log.Printf("✅ [GetRecentSales Service] Retrieved %d sales", len(sales))
	if len(sales) > 0 {
		log.Printf("🔵 [GetRecentSales Service] First sale: ID=%d, Receipt=%s, Date=%s, Amount=%.2f",
			sales[0].SaleID, sales[0].ReceiptNo, sales[0].SaleDate.Format("2006-01-02 15:04:05"), sales[0].NetAmount)
	}

	return sales, nil
}

// buildWhereClause constructs WHERE clause and args from request filters
func (s *SalesReportService) buildWhereClause(req *dto.SalesReportRequest) (string, []interface{}) {
	log.Printf("🟢 [buildWhereClause] Building WHERE clause for date range: %s to %s",
		req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02"))

	var conditions []string
	var args []interface{}

	// Always filter out NULL dates
	conditions = append(conditions, "SaleDate IS NOT NULL")

	// Date range (required) - Use SUBSTR because SQLite DATE() doesn't parse Go's time format
	conditions = append(conditions, "SUBSTR(SaleDate, 1, 10) >= ?")
	conditions = append(conditions, "SUBSTR(SaleDate, 1, 10) <= ?")
	args = append(args, req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02"))

	// Optional filters
	if req.CustomerID != nil {
		conditions = append(conditions, "CustomerID = ?")
		args = append(args, *req.CustomerID)
		log.Printf("🟢 [buildWhereClause] Adding CustomerID filter: %d", *req.CustomerID)
	}

	if req.PaymentStatus != nil {
		conditions = append(conditions, "PaymentStatus = ?")
		args = append(args, *req.PaymentStatus)
		log.Printf("🟢 [buildWhereClause] Adding PaymentStatus filter: %s", *req.PaymentStatus)
	}

	if req.PaymentMethod != nil {
		conditions = append(conditions, "PaymentMethod = ?")
		args = append(args, *req.PaymentMethod)
		log.Printf("🟢 [buildWhereClause] Adding PaymentMethod filter: %s", *req.PaymentMethod)
	}

	if req.UserID != nil {
		conditions = append(conditions, "SoldBy = ?")
		args = append(args, *req.UserID)
		log.Printf("🟢 [buildWhereClause] Adding UserID filter: %d", *req.UserID)
	}

	whereClause := "1=1"
	if len(conditions) > 0 {
		whereClause = ""
		for i, cond := range conditions {
			if i > 0 {
				whereClause += " AND "
			}
			whereClause += cond
		}
	}

	log.Printf("🟢 [buildWhereClause] Final WHERE clause: %s", whereClause)
	log.Printf("🟢 [buildWhereClause] Final args: %v", args)

	return whereClause, args
}
