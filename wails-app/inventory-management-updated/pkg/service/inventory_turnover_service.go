package service

import (
	"database/sql"
	"fmt"
	"ims-intro/pkg/service/dto"
	"time"
)

type InventoryTurnoverService struct {
	db *sql.DB
}

func NewInventoryTurnoverService(db *sql.DB) *InventoryTurnoverService {
	return &InventoryTurnoverService{db: db}
}

// GetInventoryTurnoverMetrics calculates overall inventory turnover metrics
func (s *InventoryTurnoverService) GetInventoryTurnoverMetrics() (*dto.InventoryTurnoverDTO, error) {
	// Calculate for the last 365 days
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0)

	metrics := &dto.InventoryTurnoverDTO{
		Period: fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
	}

	// Calculate COGS (Cost of Goods Sold)
	cogsQuery := `
		SELECT COALESCE(SUM(si.Quantity * p.CostPrice), 0)
		FROM SalesItems si
		INNER JOIN Sales s ON si.SaleID = s.SaleID
		INNER JOIN Products p ON si.ProductID = p.ProductID
		WHERE s.SaleDate >= ? AND s.SaleDate < ?
	`
	err := s.db.QueryRow(cogsQuery, startDate, endDate).Scan(&metrics.COGS)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate COGS: %w", err)
	}

	// Calculate average inventory value
	avgInventoryQuery := `
		SELECT COALESCE(AVG(CurrentQuantity * CostPrice), 0)
		FROM Products
		WHERE IsActive = 1
	`
	err = s.db.QueryRow(avgInventoryQuery).Scan(&metrics.AvgInventoryValue)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate average inventory: %w", err)
	}

	// Calculate turnover rate
	if metrics.AvgInventoryValue > 0 {
		metrics.TurnoverRate = metrics.COGS / metrics.AvgInventoryValue
		metrics.DaysSalesInventory = 365 / metrics.TurnoverRate
	}

	// Calculate stock-to-Sales ratio (months of inventory)
	monthlySalesQuery := `
		SELECT COALESCE(SUM(si.Quantity * p.CostPrice), 0)
		FROM SalesItems si
		INNER JOIN Sales s ON si.SaleID = s.SaleID
		INNER JOIN Products p ON si.ProductID = p.ProductID
		WHERE s.SaleDate >= ? AND s.SaleDate < ?
	`
	thirtyDaysAgo := endDate.AddDate(0, 0, -30)
	var monthlyCOGS float64
	err = s.db.QueryRow(monthlySalesQuery, thirtyDaysAgo, endDate).Scan(&monthlyCOGS)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate monthly Sales: %w", err)
	}

	if monthlyCOGS > 0 {
		metrics.StockToSalesRatio = metrics.AvgInventoryValue / monthlyCOGS
	}

	// Get product movement analysis
	productMovement, err := s.GetProductMovementAnalysis()
	if err != nil {
		return nil, fmt.Errorf("failed to get product movement: %w", err)
	}

	metrics.ProductMovement = productMovement

	// Count Products by movement class
	for _, pm := range productMovement {
		switch pm.MovementClass {
		case "fast":
			metrics.FastMovingCount++
		case "medium":
			metrics.MediumMovingCount++
		case "slow":
			metrics.SlowMovingCount++
		case "dead":
			metrics.DeadStockCount++
		}
	}

	return metrics, nil
}

// GetProductMovementAnalysis analyzes product movement patterns
func (s *InventoryTurnoverService) GetProductMovementAnalysis() ([]dto.ProductMovementDTO, error) {
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	ninetyDaysAgo := now.AddDate(0, 0, -90)

	query := `
		SELECT
			p.ProductID,
			p.Name as product_name,
			p.SKU,
			COALESCE(c.Name, 'Uncategorized') as category_name,
			p.CurrentQuantity,
			p.CurrentQuantity * p.CostPrice as stock_value,
			p.ReorderPoint,
			p.MinStockLevel,
			(SELECT MAX(s.SaleDate)
			 FROM SalesItems si
			 INNER JOIN Sales s ON si.SaleID = s.SaleID
			 WHERE si.ProductID = p.ProductID) as last_sale_date,
			COALESCE((SELECT SUM(si.Quantity)
			 FROM SalesItems si
			 INNER JOIN Sales s ON si.SaleID = s.SaleID
			 WHERE si.ProductID = p.ProductID
			 AND s.SaleDate >= ?), 0) as sold_30_days,
			COALESCE((SELECT SUM(si.Quantity)
			 FROM SalesItems si
			 INNER JOIN Sales s ON si.SaleID = s.SaleID
			 WHERE si.ProductID = p.ProductID
			 AND s.SaleDate >= ?), 0) as sold_90_days
		FROM Products p
		LEFT JOIN Categories c ON p.CategoryID = c.CategoryID
		WHERE p.IsActive = 1
		ORDER BY p.ProductID
	`

	rows, err := s.db.Query(query, thirtyDaysAgo, ninetyDaysAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to get product movement: %w", err)
	}
	defer rows.Close()

	var movements []dto.ProductMovementDTO

	for rows.Next() {
		var pm dto.ProductMovementDTO
		var lastSaleDateStr sql.NullString
		var reorderPoint int64
		var minStockLevel int64

		err := rows.Scan(
			&pm.ProductID,
			&pm.ProductName,
			&pm.SKU,
			&pm.CategoryName,
			&pm.CurrentStock,
			&pm.StockValue,
			&reorderPoint,
			&minStockLevel,
			&lastSaleDateStr,
			&pm.TotalSold30Days,
			&pm.TotalSold90Days,
		)
		if err != nil {
			return nil, err
		}

		// Parse last sale date
		if lastSaleDateStr.Valid {
			lastSaleDate, _ := time.Parse("2006-01-02 15:04:05", lastSaleDateStr.String)
			pm.LastSaleDate = &lastSaleDate

			daysSince := int(now.Sub(lastSaleDate).Hours() / 24)
			pm.DaysSinceLastSale = &daysSince
		}

		// Calculate average days between Sales (based on 90-day Sales)
		if pm.TotalSold90Days > 0 {
			avgDays := 90.0 / float64(pm.TotalSold90Days)
			pm.AvgDaysBetweenSales = &avgDays
		}

		// Calculate days of stock left (based on 30-day Sales rate)
		if pm.TotalSold30Days > 0 {
			dailySalesRate := float64(pm.TotalSold30Days) / 30.0
			if dailySalesRate > 0 {
				daysLeft := float64(pm.CurrentStock) / dailySalesRate
				pm.DaysOfStockLeft = &daysLeft
			}
		}

		// Classify movement
		pm.MovementClass = s.classifyMovement(pm.DaysSinceLastSale, pm.AvgDaysBetweenSales, pm.TotalSold90Days)

		// Reorder recommendation
		pm.ReorderRecommended = pm.CurrentStock <= reorderPoint && pm.MovementClass != "dead"

		movements = append(movements, pm)
	}

	return movements, nil
}

// GetStockHealthReport generates a comprehensive stock health report
func (s *InventoryTurnoverService) GetStockHealthReport() (*dto.StockHealthDTO, error) {
	health := &dto.StockHealthDTO{}

	// Get total Products and stock value
	totalQuery := `
		SELECT
			COUNT(*) as total_products,
			COALESCE(SUM(CurrentQuantity * CostPrice), 0) as total_stock_value
		FROM Products
		WHERE IsActive = 1
	`
	err := s.db.QueryRow(totalQuery).Scan(&health.TotalProducts, &health.TotalStockValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get totals: %w", err)
	}

	// Get Products below reorder point
	lowStockQuery := `
		SELECT
			ProductID,
			Name,
			SKU,
			CurrentQuantity,
			ReorderPoint,
			MaxStockLevel - CurrentQuantity as recommended_qty
		FROM Products
		WHERE IsActive = 1
		AND CurrentQuantity <= ReorderPoint
		AND CurrentQuantity > 0
		ORDER BY (ReorderPoint - CurrentQuantity) DESC
		LIMIT 50
	`
	lowStockRows, err := s.db.Query(lowStockQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock items: %w", err)
	}
	defer lowStockRows.Close()

	for lowStockRows.Next() {
		var item dto.LowStockProductDTO
		err := lowStockRows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.SKU,
			&item.CurrentStock,
			&item.ReorderPoint,
			&item.RecommendedQty,
		)
		if err != nil {
			continue
		}
		health.BelowReorderPoint = append(health.BelowReorderPoint, item)
	}

	// Get overstock items (above max stock level)
	overstockQuery := `
		SELECT
			p.ProductID,
			p.Name,
			p.SKU,
			p.CurrentQuantity,
			COALESCE(p.MaxStockLevel, p.MinStockLevel * 3),
			p.CurrentQuantity * p.CostPrice as stock_value
		FROM Products p
		WHERE p.IsActive = 1
		AND p.MaxStockLevel IS NOT NULL
		AND p.CurrentQuantity > p.MaxStockLevel
		ORDER BY (p.CurrentQuantity - p.MaxStockLevel) DESC
		LIMIT 50
	`
	overstockRows, err := s.db.Query(overstockQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get overstock items: %w", err)
	}
	defer overstockRows.Close()

	for overstockRows.Next() {
		var item dto.OverstockProductDTO
		err := overstockRows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.SKU,
			&item.CurrentStock,
			&item.MaxStockLevel,
			&item.StockValue,
		)
		if err != nil {
			continue
		}
		// Calculate days of supply (simplified estimate)
		item.DaysOfSupply = 90 // Placeholder - would need Sales rate calculation
		health.OverstockItems = append(health.OverstockItems, item)
	}

	// Get out of stock items
	outOfStockQuery := `
		SELECT
			p.ProductID,
			p.Name,
			p.SKU,
			(SELECT MAX(s.SaleDate)
			 FROM SalesItems si
			 INNER JOIN Sales s ON si.SaleID = s.SaleID
			 WHERE si.ProductID = p.ProductID) as last_sale_date
		FROM Products p
		WHERE p.IsActive = 1
		AND p.CurrentQuantity = 0
		ORDER BY p.ProductID DESC
		LIMIT 50
	`
	outOfStockRows, err := s.db.Query(outOfStockQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get out of stock items: %w", err)
	}
	defer outOfStockRows.Close()

	for outOfStockRows.Next() {
		var item dto.OutOfStockProductDTO
		var lastSaleDateStr sql.NullString

		err := outOfStockRows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.SKU,
			&lastSaleDateStr,
		)
		if err != nil {
			continue
		}

		if lastSaleDateStr.Valid {
			lastSaleDate, _ := time.Parse("2006-01-02 15:04:05", lastSaleDateStr.String)
			item.LastSaleDate = &lastSaleDate
			item.DaysOutOfStock = int(time.Since(lastSaleDate).Hours() / 24)
		}

		health.OutOfStock = append(health.OutOfStock, item)
	}

	return health, nil
}

// classifyMovement determines the movement class of a product
func (s *InventoryTurnoverService) classifyMovement(daysSinceLastSale *int, avgDaysBetweenSales *float64, sold90Days int64) string {
	// No Sales in 180+ days = Dead stock
	if daysSinceLastSale != nil && *daysSinceLastSale > 180 {
		return "dead"
	}

	// No Sales in 90 days period = Dead stock
	if sold90Days == 0 {
		return "dead"
	}

	// Based on average days between Sales
	if avgDaysBetweenSales != nil {
		if *avgDaysBetweenSales < 30 {
			return "fast"
		} else if *avgDaysBetweenSales < 90 {
			return "medium"
		} else {
			return "slow"
		}
	}

	// Fallback to simple classification based on 90-day Sales
	if sold90Days >= 10 {
		return "fast"
	} else if sold90Days >= 3 {
		return "medium"
	} else {
		return "slow"
	}
}
