package service

import (
	"database/sql"
	"fmt"
	"ims-intro/pkg/service/dto"
	"sort"
)

type ABCAnalysisService struct {
	db *sql.DB
}

func NewABCAnalysisService(db *sql.DB) *ABCAnalysisService {
	return &ABCAnalysisService{db: db}
}

// PerformABCAnalysis performs ABC classification of Products based on revenue contribution
func (s *ABCAnalysisService) PerformABCAnalysis(req *dto.ABCAnalysisRequest) (*dto.ABCAnalysisDTO, error) {
	// Get all Products with their revenue
	Products, totalRevenue, err := s.getProductRevenue(req)
	if err != nil {
		return nil, err
	}

	if len(Products) == 0 {
		return &dto.ABCAnalysisDTO{
			AnalysisPeriod: fmt.Sprintf("%s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02")),
			StartDate:      req.StartDate,
			EndDate:        req.EndDate,
			TotalRevenue:   0,
			Products:       []dto.ABCProductDTO{},
		}, nil
	}

	// Sort by revenue (descending)
	sort.Slice(Products, func(i, j int) bool {
		return Products[i].TotalRevenue > Products[j].TotalRevenue
	})

	// Calculate cumulative percentages and classify
	var cumulativeRevenue float64
	for i := range Products {
		Products[i].RevenuePercent = (Products[i].TotalRevenue / totalRevenue) * 100
		cumulativeRevenue += Products[i].TotalRevenue
		Products[i].CumulativePercent = (cumulativeRevenue / totalRevenue) * 100

		// ABC Classification based on cumulative percentage
		if Products[i].CumulativePercent <= 80 {
			Products[i].Classification = "A"
			Products[i].Recommendation = "High priority: Maintain tight inventory control, frequent monitoring, and ensure stock availability"
		} else if Products[i].CumulativePercent <= 95 {
			Products[i].Classification = "B"
			Products[i].Recommendation = "Medium priority: Regular monitoring, moderate inventory control"
		} else {
			Products[i].Classification = "C"
			Products[i].Recommendation = "Low priority: Periodic review, minimal inventory control, consider discontinuation if Sales decline"
		}
	}

	// Calculate class summaries
	classA := s.calculateClassSummary(Products, "A")
	classB := s.calculateClassSummary(Products, "B")
	classC := s.calculateClassSummary(Products, "C")

	analysis := &dto.ABCAnalysisDTO{
		AnalysisPeriod: fmt.Sprintf("%s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02")),
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		TotalRevenue:   totalRevenue,
		ClassA:         classA,
		ClassB:         classB,
		ClassC:         classC,
		Products:       Products,
	}

	return analysis, nil
}

// GetProductsByClass returns Products filtered by ABC classification
func (s *ABCAnalysisService) GetProductsByClass(req *dto.ABCAnalysisRequest, class string) ([]dto.ABCProductDTO, error) {
	analysis, err := s.PerformABCAnalysis(req)
	if err != nil {
		return nil, err
	}

	var filtered []dto.ABCProductDTO
	for _, product := range analysis.Products {
		if product.Classification == class {
			filtered = append(filtered, product)
		}
	}

	return filtered, nil
}

// getProductRevenue fetches product revenue data for the analysis period
func (s *ABCAnalysisService) getProductRevenue(req *dto.ABCAnalysisRequest) ([]dto.ABCProductDTO, float64, error) {
	whereClause, args := s.buildWhereClause(req)

	query := fmt.Sprintf(`
		SELECT
			p.ProductID,
			p.Name as product_name,
			p.SKU,
			COALESCE(c.Name, 'Uncategorized') as category_name,
			COALESCE(SUM(si.LineTotal), 0) as total_revenue,
			COALESCE(SUM(si.Quantity), 0) as quantity_sold,
			p.CurrentQuantity as current_stock
		FROM SalesItems si
		INNER JOIN Sales s ON si.SaleID = s.SaleID
		INNER JOIN Products p ON si.ProductID = p.ProductID
		LEFT JOIN Categories c ON p.CategoryID = c.CategoryID
		WHERE %s
		GROUP BY p.ProductID, p.Name, p.SKU, c.Name, p.CurrentQuantity
		HAVING total_revenue > 0
	`, whereClause)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get product revenue: %w", err)
	}
	defer rows.Close()

	var Products []dto.ABCProductDTO
	var totalRevenue float64

	for rows.Next() {
		var product dto.ABCProductDTO
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.SKU,
			&product.CategoryName,
			&product.TotalRevenue,
			&product.QuantitySold,
			&product.CurrentStock,
		)
		if err != nil {
			return nil, 0, err
		}

		totalRevenue += product.TotalRevenue
		Products = append(Products, product)
	}

	return Products, totalRevenue, nil
}

// calculateClassSummary calculates summary statistics for a specific class
func (s *ABCAnalysisService) calculateClassSummary(Products []dto.ABCProductDTO, class string) dto.ABCClassDTO {
	var count int
	var revenue float64

	for _, product := range Products {
		if product.Classification == class {
			count++
			revenue += product.TotalRevenue
		}
	}

	var revenuePercent float64
	var totalRevenue float64
	for _, product := range Products {
		totalRevenue += product.TotalRevenue
	}
	if totalRevenue > 0 {
		revenuePercent = (revenue / totalRevenue) * 100
	}

	var recommendation string
	switch class {
	case "A":
		recommendation = "Focus on these high-value items: ensure consistent availability, negotiate better supplier terms, implement tight inventory controls"
	case "B":
		recommendation = "Monitor regularly: maintain adequate stock levels, review performance monthly, adjust pricing strategies as needed"
	case "C":
		recommendation = "Minimize attention: reduce stock levels, consider bulk ordering, evaluate discontinuation of very slow movers"
	}

	return dto.ABCClassDTO{
		ProductCount:   count,
		TotalRevenue:   revenue,
		RevenuePercent: revenuePercent,
		Recommendation: recommendation,
	}
}

// buildWhereClause constructs WHERE clause and args from request filters
func (s *ABCAnalysisService) buildWhereClause(req *dto.ABCAnalysisRequest) (string, []interface{}) {
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
