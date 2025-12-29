package service

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
)

type ReceiptService struct {
	productRepo repository.IProductRepository
	settingRepo repository.ISettingRepository
}

func NewReceiptService(productRepo repository.IProductRepository, settingRepo repository.ISettingRepository) *ReceiptService {
	return &ReceiptService{
		productRepo: productRepo,
		settingRepo: settingRepo,
	}
}

// ReceiptData contains all information needed to generate a receipt
type ReceiptData struct {
	SaleWithItems *domain.SaleWithItems
	CustomerName  string // Optional
}

// ReceiptItem combines sale item with product details
type ReceiptItem struct {
	ProductName     string
	ProductSKU      string
	Quantity        int64
	UnitPrice       float64
	TaxRate         float64
	DiscountPercent float64
	LineTotal       float64
}

// GenerateReceipt generates a PDF receipt and returns it as a byte buffer
func (s *ReceiptService) GenerateReceipt(data *ReceiptData) (*bytes.Buffer, error) {
	// Fetch product details for all items
	receiptItems, err := s.enrichItemsWithProductDetails(data.SaleWithItems.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich items with product details: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Shop header
	s.addShopHeader(pdf)

	// Receipt details
	s.addReceiptDetails(pdf, data)

	// Items table
	s.addItemsTable(pdf, receiptItems)

	// Totals
	s.addTotals(pdf, data.SaleWithItems)

	// Store policy footer
	s.addStorePolicy(pdf)

	// Generate PDF to buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return &buf, nil
}

// enrichItemsWithProductDetails fetches product information for sale items
func (s *ReceiptService) enrichItemsWithProductDetails(items []domain.SaleItem) ([]ReceiptItem, error) {
	receiptItems := make([]ReceiptItem, len(items))

	for i, item := range items {
		product, err := s.productRepo.GetProductByID(item.ProductID)
		if err != nil {
			// If product not found, use placeholder
			receiptItems[i] = ReceiptItem{
				ProductName:     fmt.Sprintf("Product ID %d", item.ProductID),
				ProductSKU:      "N/A",
				Quantity:        item.Quantity,
				UnitPrice:       item.UnitPrice,
				TaxRate:         item.TaxRate,
				DiscountPercent: item.DiscountPercent,
				LineTotal:       item.LineTotal,
			}
		} else {
			receiptItems[i] = ReceiptItem{
				ProductName:     product.Name,
				ProductSKU:      product.SKU,
				Quantity:        item.Quantity,
				UnitPrice:       item.UnitPrice,
				TaxRate:         item.TaxRate,
				DiscountPercent: item.DiscountPercent,
				LineTotal:       item.LineTotal,
			}
		}
	}

	return receiptItems, nil
}

// addShopHeader adds the shop information at the top
func (s *ReceiptService) addShopHeader(pdf *gofpdf.Fpdf) {
	// Fetch settings from database
	storeName := "Store Name"
	storeAddress := "Store Address"
	storePhone := "Store Phone"

	if setting, err := s.settingRepo.GetSettingByKey("store_name"); err == nil {
		storeName = setting.SettingValue
	}
	if setting, err := s.settingRepo.GetSettingByKey("store_address"); err == nil {
		storeAddress = setting.SettingValue
	}
	if setting, err := s.settingRepo.GetSettingByKey("store_phone"); err == nil {
		storePhone = setting.SettingValue
	}

	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 10, storeName, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, storeAddress, "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Phone: %s", storePhone), "", 1, "C", false, 0, "")

	// Separator line
	pdf.Ln(5)
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(5)
}

// addReceiptDetails adds receipt number, date, and customer name
func (s *ReceiptService) addReceiptDetails(pdf *gofpdf.Fpdf, data *ReceiptData) {
	pdf.SetFont("Arial", "", 11)

	// Receipt number
	pdf.CellFormat(50, 6, "Receipt #:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 6, data.SaleWithItems.ReceiptNo, "", 1, "L", false, 0, "")

	// Date and time
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(50, 6, "Date:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	formattedDate := data.SaleWithItems.SaleDate.Format("January 02, 2006 03:04 PM")
	pdf.CellFormat(0, 6, formattedDate, "", 1, "L", false, 0, "")

	// Customer name (optional)
	if data.CustomerName != "" {
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(50, 6, "Customer:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "B", 11)
		pdf.CellFormat(0, 6, data.CustomerName, "", 1, "L", false, 0, "")
	}

	pdf.Ln(5)
}

// addItemsTable adds the table of sold items
func (s *ReceiptService) addItemsTable(pdf *gofpdf.Fpdf, items []ReceiptItem) {
	// Table header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)

	pdf.CellFormat(70, 7, "Item", "1", 0, "L", true, 0, "")
	pdf.CellFormat(15, 7, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 7, "Price", "1", 0, "R", true, 0, "")
	pdf.CellFormat(20, 7, "Tax%", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 7, "Disc%", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 7, "Total", "1", 1, "R", true, 0, "")

	// Table rows
	pdf.SetFont("Arial", "", 10)
	for _, item := range items {
		pdf.CellFormat(70, 6, item.ProductName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(15, 6, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("$%.2f", item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.1f%%", item.TaxRate), "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.1f%%", item.DiscountPercent), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 6, fmt.Sprintf("$%.2f", item.LineTotal), "1", 1, "R", false, 0, "")

		// Product SKU
		if item.ProductSKU != "N/A" && item.ProductSKU != "" {
			pdf.SetFont("Arial", "I", 8)
			pdf.CellFormat(180, 4, fmt.Sprintf("  (%s)", item.ProductSKU), "LR", 1, "L", false, 0, "")
			pdf.SetFont("Arial", "", 10)
		}
	}

	pdf.Ln(3)
}

// addTotals adds the totals section
func (s *ReceiptService) addTotals(pdf *gofpdf.Fpdf, sale *domain.SaleWithItems) {
	// Right-aligned totals
	rightColX := float64(150)
	labelWidth := float64(30)
	valueWidth := float64(30)

	pdf.SetFont("Arial", "", 11)

	// Subtotal (Total - Tax + Discount)
	subtotal := sale.TotalAmount
	pdf.SetX(rightColX)
	pdf.CellFormat(labelWidth, 6, "Subtotal:", "", 0, "L", false, 0, "")
	pdf.CellFormat(valueWidth, 6, fmt.Sprintf("$%.2f", subtotal), "", 1, "R", false, 0, "")

	// Tax
	pdf.SetX(rightColX)
	pdf.CellFormat(labelWidth, 6, "Tax:", "", 0, "L", false, 0, "")
	pdf.CellFormat(valueWidth, 6, fmt.Sprintf("$%.2f", sale.TaxAmount), "", 1, "R", false, 0, "")

	// Discount
	if sale.DiscountAmount > 0 {
		pdf.SetX(rightColX)
		pdf.CellFormat(labelWidth, 6, "Discount:", "", 0, "L", false, 0, "")
		pdf.CellFormat(valueWidth, 6, fmt.Sprintf("-$%.2f", sale.DiscountAmount), "", 1, "R", false, 0, "")
	}

	// Separator line
	pdf.SetX(rightColX)
	pdf.SetDrawColor(0, 0, 0)
	lineY := pdf.GetY()
	pdf.Line(rightColX, lineY, rightColX+labelWidth+valueWidth, lineY)
	pdf.Ln(2)

	// Total (Net Amount)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetX(rightColX)
	pdf.CellFormat(labelWidth, 8, "TOTAL:", "", 0, "L", false, 0, "")
	pdf.CellFormat(valueWidth, 8, fmt.Sprintf("$%.2f", sale.NetAmount), "", 1, "R", false, 0, "")

	// Payment method
	if sale.PaymentMethod != nil {
		pdf.SetFont("Arial", "", 10)
		pdf.Ln(2)
		pdf.SetX(rightColX)
		pdf.CellFormat(labelWidth, 6, "Paid:", "", 0, "L", false, 0, "")
		paymentMethod := *sale.PaymentMethod
		// Capitalize first letter
		if len(paymentMethod) > 0 {
			paymentMethod = fmt.Sprintf("%c%s", paymentMethod[0]-32, paymentMethod[1:])
		}
		pdf.CellFormat(valueWidth, 6, paymentMethod, "", 1, "R", false, 0, "")
	}

	pdf.Ln(10)
}

// addStorePolicy adds the store policy footer
func (s *ReceiptService) addStorePolicy(pdf *gofpdf.Fpdf) {
	// Separator line
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(5)

	// Policy header
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 6, "Store Policy", "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Fetch policy from settings
	policyText := "Thank you for shopping with us!"
	if setting, err := s.settingRepo.GetSettingByKey("store_policy"); err == nil {
		policyText = setting.SettingValue
	}

	// Policy text
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 5, policyText, "", "C", false)
}
