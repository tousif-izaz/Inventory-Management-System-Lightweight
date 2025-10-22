# Inventory Management System - API Specification

## Overview
This document outlines all API endpoints for the Inventory Management System based on the schema in `schema_sqlite.sql` (version 2.2.0).

## Schema Version
**Current Schema:** 2.2.0 - Simplified schema with inventory tracking directly in Products table

## Base URL
```
http://localhost:8080
```

## Authentication
All endpoints except `/login` and `/signup` require authentication via JWT token in Cookie.

## Schema Changes Summary (v2.2.0)

This version introduces significant simplifications to the database schema:

### What Changed:
1. **Removed Inventory Table** - Inventory is now tracked directly in the Products table
2. **Removed Barcode Field** - Simplified product identification to SKU only
3. **Added CurrentQuantity** - Direct quantity tracking on each product
4. **Removed Location References** - Simplified to single ShelfLocation field
5. **Simplified Views** - Stock views now query Products table directly

### Migration Impact:
- **Products API:** Request/response schemas updated (removed `barcode`, added `current_quantity`)
- **Inventory API:** Endpoints removed (use Products endpoints instead)
- **Purchases/Sales API:** `location_id` removed from items, quantities update Products directly
- **Frontend:** Types and forms updated to reflect new schema

### Benefits:
- Simpler database structure with fewer joins
- Easier to understand and maintain
- Better performance for single-location businesses
- Direct product quantity management

---

## 1. Categories API

### Create Category
- **Endpoint:** `POST /categories`
- **Auth Required:** Yes (admin, manager)
- **Request Body:**
```json
{
  "name": "Electronics",
  "description": "Electronic products",
  "parent_category_id": null
}
```
- **Response:** `201 Created`

### List All Categories
- **Endpoint:** `GET /categories`
- **Auth Required:** Yes
- **Query Params:** None
- **Response:** `200 OK`
```json
[
  {
    "category_id": 1,
    "name": "Electronics",
    "description": "Electronic products",
    "parent_category_id": null,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### Get Category by ID
- **Endpoint:** `GET /categories/:id`
- **Auth Required:** Yes
- **Response:** `200 OK` (single category object)

### Get Child Categories
- **Endpoint:** `GET /categories/:id/children`
- **Auth Required:** Yes
- **Response:** `200 OK` (array of child categories)

### Update Category
- **Endpoint:** `PUT /categories/:id`
- **Auth Required:** Yes (admin, manager)
- **Request Body:** Same as Create
- **Response:** `200 OK`

### Delete Category
- **Endpoint:** `DELETE /categories/:id`
- **Auth Required:** Yes (admin)
- **Response:** `204 No Content`

---

## 2. Suppliers API

### Create Supplier
- **Endpoint:** `POST /suppliers`
- **Auth Required:** Yes (admin, manager)
- **Request Body:**
```json
{
  "name": "ABC Suppliers Ltd",
  "contact_person": "John Doe",
  "email": "contact@abc.com",
  "phone": "+1234567890",
  "address": "123 Main St",
  "city": "New York",
  "country": "USA",
  "tax_id": "TAX123456",
  "payment_terms": "Net 30"
}
```
- **Response:** `201 Created`

### List All Suppliers
- **Endpoint:** `GET /suppliers`
- **Auth Required:** Yes
- **Query Params:**
  - `is_active` (optional): `true` | `false`
- **Response:** `200 OK` (array of suppliers)

### Get Supplier by ID
- **Endpoint:** `GET /suppliers/:id`
- **Auth Required:** Yes
- **Response:** `200 OK`

### Update Supplier
- **Endpoint:** `PUT /suppliers/:id`
- **Auth Required:** Yes (admin, manager)
- **Request Body:** Same as Create
- **Response:** `200 OK`

### Deactivate Supplier
- **Endpoint:** `DELETE /suppliers/:id`
- **Auth Required:** Yes (admin)
- **Description:** Soft delete (sets IsActive=0)
- **Response:** `204 No Content`

---

## 3. Customers API

### Create Customer
- **Endpoint:** `POST /customers`
- **Auth Required:** Yes (staff+)
- **Request Body:**
```json
{
  "name": "Jane Smith",
  "email": "jane@example.com",
  "phone": "+1234567890",
  "address": "789 Customer St",
  "city": "Los Angeles",
  "country": "USA"
}
```
- **Response:** `201 Created`

### List All Customers
- **Endpoint:** `GET /customers`
- **Auth Required:** Yes
- **Query Params:**
  - `is_active` (optional): `true` | `false`
- **Response:** `200 OK`

### Get Customer by ID
- **Endpoint:** `GET /customers/:id`
- **Auth Required:** Yes
- **Response:** `200 OK`

### Get Customer Purchase History
- **Endpoint:** `GET /customers/:id/sales`
- **Auth Required:** Yes
- **Response:** `200 OK` (array of sales)

### Update Customer
- **Endpoint:** `PUT /customers/:id`
- **Auth Required:** Yes (staff+)
- **Response:** `200 OK`

### Update Customer Loyalty Points
- **Endpoint:** `PUT /customers/:id/loyalty`
- **Auth Required:** Yes (manager+)
- **Request Body:**
```json
{
  "loyalty_points": 100
}
```
- **Response:** `200 OK`

---

## 4. Products API

### Create Product
- **Endpoint:** `POST /products`
- **Auth Required:** Yes (manager+)
- **Request Body:**
```json
{
  "name": "Premium Perfume",
  "description": "Luxury fragrance",
  "sku": "PERF-001",
  "category_id": 1,
  "batch_no": "BATCH-2024-01",
  "expiry_date": "2026-01-01T00:00:00Z",
  "cost_price": 50.00,
  "selling_price": 89.99,
  "current_quantity": 0,
  "min_stock_level": 10,
  "max_stock_level": 100,
  "reorder_point": 20,
  "unit": "pcs",
  "shelf_location": "A1-B2"
}
```
- **Response:** `201 Created`
- **Notes:**
  - `barcode` field has been removed in schema v2.2.0
  - `current_quantity` is now managed directly in the product record
  - `shelf_location` is optional for tracking physical location

### List All Products
- **Endpoint:** `GET /products`
- **Auth Required:** Yes
- **Query Params:**
  - `category_id` (optional)
  - `is_active` (optional)
- **Response:** `200 OK`

### Get Product by ID
- **Endpoint:** `GET /products/:id`
- **Auth Required:** Yes
- **Response:** `200 OK`

### Get Product by SKU
- **Endpoint:** `GET /products/sku/:sku`
- **Auth Required:** Yes
- **Response:** `200 OK`

### Get Low Stock Products
- **Endpoint:** `GET /products/low-stock`
- **Auth Required:** Yes
- **Description:** Returns products where current_quantity <= reorder_point
- **Response:** `200 OK`

### Get Expiring Products
- **Endpoint:** `GET /products/expiring`
- **Auth Required:** Yes
- **Query Params:**
  - `days` (optional, default=30): Number of days to check
- **Description:** Returns products expiring within specified days
- **Response:** `200 OK`

### Update Product
- **Endpoint:** `PUT /products/:id`
- **Auth Required:** Yes (manager+)
- **Response:** `200 OK`

### Deactivate Product
- **Endpoint:** `DELETE /products/:id`
- **Auth Required:** Yes (admin)
- **Response:** `204 No Content`

---

## 5. Purchases API

### Create Purchase
- **Endpoint:** `POST /purchases`
- **Auth Required:** Yes (manager+)
- **Request Body:**
```json
{
  "purchase_date": "2024-01-01T00:00:00Z",
  "supplier_id": 1,
  "invoice_number": "INV-2024-001",
  "reference_number": "REF-001",
  "tax_amount": 15.00,
  "discount_amount": 5.00,
  "payment_status": "pending",
  "payment_method": "bank_transfer",
  "notes": "Bulk order",
  "received_by": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 100,
      "unit_cost": 50.00,
      "tax_rate": 10.0,
      "discount_percent": 5.0,
      "batch_no": "BATCH-2024-01",
      "expiry_date": "2026-01-01T00:00:00Z"
    }
  ]
}
```
- **Business Logic:**
  1. Calculate line totals for each item
  2. Calculate total_amount and net_amount
  3. Create Purchase record
  4. Create PurchaseItems records
  5. **Increase Product.CurrentQuantity** for each item
  6. Create Transaction records (type=purchase) for audit
  7. Optionally update Product.CostPrice and Product.BatchNo/ExpiryDate
- **Response:** `201 Created` with purchase ID
- **Notes:**
  - Schema v2.2.0: Inventory is now tracked directly in Products table
  - `location_id` removed from purchase items (use shelf_location in product)

### List All Purchases
- **Endpoint:** `GET /purchases`
- **Auth Required:** Yes
- **Query Params:**
  - `supplier_id` (optional)
  - `payment_status` (optional): `pending` | `partial` | `paid`
  - `from_date` (optional): ISO date
  - `to_date` (optional): ISO date
- **Response:** `200 OK`

### Get Purchase by ID
- **Endpoint:** `GET /purchases/:id`
- **Auth Required:** Yes
- **Description:** Returns purchase with all line items
- **Response:** `200 OK`

### Update Purchase Payment Status
- **Endpoint:** `PUT /purchases/:id/payment-status`
- **Auth Required:** Yes (manager+)
- **Request Body:**
```json
{
  "payment_status": "paid"
}
```
- **Response:** `200 OK`

---

## 6. Sales API

### Create Sale
- **Endpoint:** `POST /sales`
- **Auth Required:** Yes (staff+)
- **Request Body:**
```json
{
  "sale_date": "2024-01-01T00:00:00Z",
  "receipt_no": "RCP-2024-001",
  "customer_id": 1,
  "tax_amount": 8.10,
  "discount_amount": 5.00,
  "payment_status": "paid",
  "payment_method": "card",
  "notes": "Gift wrapped",
  "sold_by": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 2,
      "unit_price": 89.99,
      "tax_rate": 10.0,
      "discount_percent": 5.0
    }
  ]
}
```
- **Business Logic:**
  1. Validate sufficient inventory (check Product.CurrentQuantity)
  2. Calculate line totals for each item
  3. Calculate total_amount and net_amount
  4. Create Sale record
  5. Create SalesItems records
  6. **Reduce Product.CurrentQuantity** for each item
  7. Create Transaction records (type=sale) for audit trail
  8. Update Customer.LoyaltyPoints if customer specified
- **Response:** `201 Created` with sale ID
- **Notes:**
  - Schema v2.2.0: Inventory is now tracked directly in Products table
  - `location_id` removed from sale items

### List All Sales
- **Endpoint:** `GET /sales`
- **Auth Required:** Yes
- **Query Params:**
  - `customer_id` (optional)
  - `payment_status` (optional): `pending` | `partial` | `paid` | `refunded`
  - `from_date` (optional): ISO date
  - `to_date` (optional): ISO date
- **Response:** `200 OK`

### Get Sale by ID
- **Endpoint:** `GET /sales/:id`
- **Auth Required:** Yes
- **Description:** Returns sale with all line items
- **Response:** `200 OK`

### Update Sale Payment Status
- **Endpoint:** `PUT /sales/:id/payment-status`
- **Auth Required:** Yes (staff+)
- **Request Body:**
```json
{
  "payment_status": "paid"
}
```
- **Response:** `200 OK`

### Process Refund
- **Endpoint:** `POST /sales/:id/refund`
- **Auth Required:** Yes (manager+)
- **Description:** Process full refund and restore inventory
- **Business Logic:**
  1. Update Sale.PaymentStatus to "refunded"
  2. **Restore Product.CurrentQuantity** for all items
  3. Create Transaction records (type=return) for audit
  4. Adjust Customer.LoyaltyPoints if applicable
- **Response:** `200 OK`

### Generate Receipt PDF
- **Endpoint:** `GET /sales/:id/receipt`
- **Auth Required:** Yes (staff+)
- **Path Parameter:**
  - `id` (required): Sale ID
- **Query Parameters:**
  - `customer_name` (optional): Customer name to display on receipt
- **Description:** Generate and download PDF receipt for a sale
- **Response:** PDF file (application/pdf)
- **Response Headers:**
  - `Content-Type: application/pdf`
  - `Content-Disposition: attachment; filename="receipt-{receipt_no}.pdf"`
- **Example Requests:**
  ```
  GET /sales/123/receipt
  GET /sales/123/receipt?customer_name=Jane%20Smith
  ```
- **Receipt Contents:**
  - **Header:** Shop name, address, phone (Lux perfumes, 1950 E20th St # G-710, Chico, CA 95928, Phone: 530 717 3219)
  - **Details:** Receipt number, date/time, customer name (if provided)
  - **Items Table:** Product name/SKU, quantity, unit price, tax rate, discount, line total
  - **Totals:** Subtotal, tax amount, discount amount, net amount, payment method
  - **Footer:** Store policy (no refunds, exchanges within 14 days with original packaging)
- **Error Responses:**
  - `400 Bad Request` - Invalid sale ID
  - `404 Not Found` - Sale not found
  - `500 Internal Server Error` - PDF generation failed

---

## 7. Transactions API (Audit Trail)

### List All Transactions
- **Endpoint:** `GET /transactions`
- **Auth Required:** Yes (manager+)
- **Query Params:**
  - `transaction_type` (optional): `purchase` | `sale` | `adjustment` | `return` | `damage` | `theft`
  - `product_id` (optional)
  - `from_date` (optional): ISO date
  - `to_date` (optional): ISO date
- **Response:** `200 OK`
- **Notes:**
  - Schema v2.2.0: `location_id` removed from transactions

### Get Transactions for Product
- **Endpoint:** `GET /transactions/product/:productId`
- **Auth Required:** Yes (manager+)
- **Description:** Complete history of all movements for a product
- **Response:** `200 OK`

### Get Transactions by User
- **Endpoint:** `GET /transactions/user/:userId`
- **Auth Required:** Yes (admin)
- **Description:** All transactions performed by a specific user
- **Response:** `200 OK`

---

## 8. Reports API

### Stock Summary Report
- **Endpoint:** `GET /reports/stock-summary`
- **Auth Required:** Yes
- **Description:** Uses vw_product_stock_summary view
- **Response:** `200 OK`
```json
[
  {
    "product_id": 1,
    "name": "Premium Perfume",
    "sku": "PERF-001",
    "category": "Fragrances",
    "cost_price": 50.00,
    "selling_price": 89.99,
    "total_stock": 85,
    "min_stock_level": 10,
    "reorder_point": 20,
    "stock_status": "OK"
  }
]
```
- **Notes:**
  - Schema v2.2.0: `barcode` field removed
  - `total_stock` now comes from `Product.CurrentQuantity`
  - Reserved/available quantity tracking removed

### Inventory Details Report
- **Endpoint:** `GET /reports/inventory-details`
- **Auth Required:** Yes
- **Description:** Uses vw_inventory_details view for detailed product inventory
- **Response:** `200 OK`

### Sales Summary Report
- **Endpoint:** `GET /reports/sales-summary`
- **Auth Required:** Yes (manager+)
- **Query Params:**
  - `from_date` (optional)
  - `to_date` (optional)
  - `group_by` (optional): `day` | `week` | `month`
- **Response:** `200 OK`

### Purchase Summary Report
- **Endpoint:** `GET /reports/purchase-summary`
- **Auth Required:** Yes (manager+)
- **Query Params:**
  - `from_date` (optional)
  - `to_date` (optional)
  - `supplier_id` (optional)
- **Response:** `200 OK`

### Low Stock Alerts
- **Endpoint:** `GET /reports/low-stock-alerts`
- **Auth Required:** Yes
- **Description:** Products at or below reorder point
- **Response:** `200 OK`

### Expiring Products Report
- **Endpoint:** `GET /reports/expiring-products`
- **Auth Required:** Yes
- **Query Params:**
  - `days` (optional, default=30): Days until expiry
- **Response:** `200 OK`

---

## Error Responses

All endpoints return standard error responses:

```json
{
  "error_message": "Description of the error"
}
```

### HTTP Status Codes:
- `200 OK` - Success
- `201 Created` - Resource created
- `204 No Content` - Success with no response body
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation error
- `500 Internal Server Error` - Server error

---

## Business Logic Summary

### Critical Inventory Operations (Schema v2.2.0):

1. **Sale Transaction:**
   - Validates stock availability (checks Product.CurrentQuantity)
   - Reduces Product.CurrentQuantity directly
   - Creates audit trail (transactions table)
   - Updates customer loyalty points

2. **Purchase Transaction:**
   - Increases Product.CurrentQuantity directly
   - Creates audit trail
   - Updates Product.CostPrice, BatchNo, ExpiryDate (optional)

3. **Inventory Adjustment:**
   - Manual corrections for discrepancies
   - Updates Product.CurrentQuantity directly
   - Records reason and performer
   - Full audit trail

### Schema Simplifications (v2.2.0):
- **Removed:** Separate Inventory table
- **Removed:** Multi-location inventory tracking
- **Removed:** Barcode field from Products
- **Added:** CurrentQuantity field directly in Products table
- **Added:** ShelfLocation field for simple physical location tracking

---

## Implementation Notes

- All timestamps use RFC3339 format
- All monetary values use float64 (consider decimal library for production)
- Soft deletes used for suppliers, locations, customers (IsActive=0)
- Hard deletes used for categories (restricted if products exist)
- Cascading deletes for purchase/sale items when parent is deleted
- Foreign key constraints prevent orphaned records
- Indexes optimize query performance
- Triggers automatically update timestamps
- Views provide pre-formatted report data
