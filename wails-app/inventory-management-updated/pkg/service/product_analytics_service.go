package service

import (
	"database/sql"
	"fmt"
	"ims-intro/pkg/service/dto"
)

type ProductAnalyticsService struct {
	db *sql.DB
}

func NewProductAnalyticsService(db *sql.DB) *ProductAnalyticsService {
	return &ProductAnalyticsService{db: db}
}

// GetBestSellingProducts returns top performing Products
func (s *ProductAnalyticsService) GetBestSellingProducts(req *dto.ProductPerformanceRequest) ([]dto.ProductPerformanceDTO, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	whereClause, args := s.buildWhereClause(req)
	args = append(args, req.Limit)

	query := fmt.Sprintf(`
		SELECT
			p.ProductID,
			p.Name as product_name,
			p.SKU,
			COALESCE(c.Name, 'Uncategorized') as category_name,
			COALESCE(SUM(si.Quantity), 0) as quantity_sold,
			COALESCE(SUM(si.LineTotal), 0) as total_revenue,
			COALESCE(AVG(si.UnitPrice), 0) as avg_unit_price,
			COALESCE(SUM(si.Quantity * p.CostPrice), 0) as total_cost,
			COALESCE(AVG(si.DiscountPercent), 0) as avg_discount,
			p.CurrentQuantity as current_stock
		FROM SalesItems si
		INNER JOIN Sales s ON si.SaleID = s.SaleID
		INNER JOIN Products p ON si.ProductID = p.ProductID
		LEFT JOIN Categories c ON p.CategoryID = c.CategoryID
		WHERE %s
		GROUP BY p.ProductID, p.Name, p.SKU, c.Name, p.CurrentQuantity, p.CostPrice
		ORDER BY total_revenue DESC
		LIMIT ?
	`, whereClause)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get best selling Products: %w", err)
	}
	defer rows.Close()

	var Products []dto.ProductPerformanceDTO
	rank := 1

	for rows.Next() {
		var product dto.ProductPerformanceDTO
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.SKU,
			&product.CategoryName,
			&product.QuantitySold,
			&product.TotalRevenue,
			&product.AvgUnitPrice,
			&product.TotalCost,
			&product.AvgDiscount,
			&product.CurrentStock,
		)
		if err != nil {
			return nil, err
		}

		// Calculate gross profit and margin
		product.GrossProfit = product.TotalRevenue - product.TotalCost
		if product.TotalRevenue > 0 {
			product.ProfitMargin = (product.GrossProfit / product.TotalRevenue) * 100
		}

		product.Rank = rank
		rank++

		Products = append(Products, product)
	}

	return Products, nil
}

// GetWorstPerformingProducts returns underperforming Products
func (s *ProductAnalyticsService) GetWorstPerformingProducts(req *dto.ProductPerformanceRequest) ([]dto.ProductPerformanceDTO, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	whereClause, args := s.buildWhereClause(req)
	args = append(args, req.Limit)

	query := fmt.Sprintf(`
		SELECT
			p.ProductID,
			p.Name as product_name,
			p.SKU,
			COALESCE(c.Name, 'Uncategorized') as category_name,
			COALESCE(SUM(si.Quantity), 0) as quantity_sold,
			COALESCE(SUM(si.LineTotal), 0) as total_revenue,
			COALESCE(AVG(si.UnitPrice), 0) as avg_unit_price,
			COALESCE(SUM(si.Quantity * p.CostPrice), 0) as total_cost,
			COALESCE(AVG(si.DiscountPercent), 0) as avg_discount,
			p.CurrentQuantity as current_stock
		FROM SalesItems si
		INNER JOIN Sales s ON si.SaleID = s.SaleID
		INNER JOIN Products p ON si.ProductID = p.ProductID
		LEFT JOIN Categories c ON p.CategoryID = c.CategoryID
		WHERE %s
		GROUP BY p.ProductID, p.Name, p.SKU, c.Name, p.CurrentQuantity, p.CostPrice
		ORDER BY total_revenue ASC
		LIMIT ?
	`, whereClause)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get worst performing Products: %w", err)
	}
	defer rows.Close()

	var Products []dto.ProductPerformanceDTO
	rank := 1

	for rows.Next() {
		var product dto.ProductPerformanceDTO
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.SKU,
			&product.CategoryName,
			&product.QuantitySold,
			&product.TotalRevenue,
			&product.AvgUnitPrice,
			&product.TotalCost,
			&product.AvgDiscount,
			&product.CurrentStock,
		)
		if err != nil {
			return nil, err
		}

		// Calculate gross profit and margin
		product.GrossProfit = product.TotalRevenue - product.TotalCost
		if product.TotalRevenue > 0 {
			product.ProfitMargin = (product.GrossProfit / product.TotalRevenue) * 100
		}

		product.Rank = rank
		rank++

		Products = append(Products, product)
	}

	return Products, nil
}

// GetCategoryPerformance returns performance metrics by category
func (s *ProductAnalyticsService) GetCategoryPerformance(req *dto.ProductPerformanceRequest) ([]dto.CategoryPerformanceDTO, error) {
	whereClause, args := s.buildWhereClause(req)

	query := fmt.Sprintf(`
		SELECT
			c.CategoryID,
			c.Name as category_name,
			COUNT(DISTINCT p.ProductID) as product_count,
			COALESCE(SUM(si.Quantity), 0) as quantity_sold,
			COALESCE(SUM(si.LineTotal), 0) as total_revenue,
			COALESCE(AVG(si.UnitPrice), 0) as avg_price
		FROM SalesItems si
		INNER JOIN Sales s ON si.SaleID = s.SaleID
		INNER JOIN Products p ON si.ProductID = p.ProductID
		INNER JOIN Categories c ON p.CategoryID = c.CategoryID
		WHERE %s
		GROUP BY c.CategoryID, c.Name
		ORDER BY total_revenue DESC
	`, whereClause)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get category performance: %w", err)
	}
	defer rows.Close()

	var Categories []dto.CategoryPerformanceDTO
	var totalRevenue float64

	for rows.Next() {
		var category dto.CategoryPerformanceDTO
		err := rows.Scan(
			&category.CategoryID,
			&category.CategoryName,
			&category.ProductCount,
			&category.QuantitySold,
			&category.TotalRevenue,
			&category.AvgPrice,
		)
		if err != nil {
			return nil, err
		}

		totalRevenue += category.TotalRevenue
		Categories = append(Categories, category)
	}

	// Calculate percentages
	for i := range Categories {
		if totalRevenue > 0 {
			Categories[i].Percentage = (Categories[i].TotalRevenue / totalRevenue) * 100
		}
	}

	return Categories, nil
}

// GetProductsByCategory returns Products filtered by category with performance metrics
func (s *ProductAnalyticsService) GetProductsByCategory(categoryID int64, req *dto.ProductPerformanceRequest) ([]dto.ProductPerformanceDTO, error) {
	req.CategoryID = &categoryID
	return s.GetBestSellingProducts(req)
}

// buildWhereClause constructs WHERE clause and args from request filters
func (s *ProductAnalyticsService) buildWhereClause(req *dto.ProductPerformanceRequest) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	// Date range (required)
	conditions = append(conditions, "s.SaleDate >= ? AND s.SaleDate < ?")
	args = append(args, req.StartDate, req.EndDate)

	// Optional filters
	if req.CategoryID != nil {
		conditions = append(conditions, "p.CategoryID = ?")
		args = append(args, *req.CategoryID)
	}

	if req.SupplierID != nil {
		conditions = append(conditions, "p.supplier_id = ?")
		args = append(args, *req.SupplierID)
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

	return whereClause, args
}
