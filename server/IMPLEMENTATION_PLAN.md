# Extended IMS - Implementation Plan & Progress

## Overview
This document tracks the implementation of the extended Inventory Management System based on `schema_sqlite_extended.sql`.

---

## Project Structure

```
server/pkg/
├── domain/                  # Domain models (entities)
│   ├── category.go         ✅ Created
│   ├── supplier.go         ✅ Created
│   ├── location.go         ✅ Created
│   ├── customer.go         ✅ Created
│   ├── inventory.go        ✅ Created
│   ├── purchase.go         ✅ Created
│   ├── sale.go             ✅ Created
│   ├── transaction.go      ✅ Created
│   ├── product.go          ✅ Updated (major changes)
│   └── user.go             ✅ Updated (additional fields)
│
├── service/dto/            # Data Transfer Objects
│   └── model.go            ✅ Updated (all new DTOs added)
│
├── controller/request/     # Request models
│   ├── category_request.go     ✅ Created
│   ├── supplier_request.go     ✅ Created
│   ├── purchase_request.go     ✅ Created
│   ├── sale_request.go         ✅ Created
│   ├── inventory_request.go    ✅ Created
│   ├── product_request.go      ⏳ Needs update
│   └── user_request.go         ⏳ Needs update
│
├── controller/response/    # Response models
│   ├── productResponse.go      ⏳ Needs update
│   ├── user_response.go        ⏳ Needs update
│   └── [NEW responses]         ⏳ To be created
│
├── repository/             # Data access layer
│   ├── category_repository.go  ⏳ To be created
│   ├── supplier_repository.go  ⏳ To be created
│   ├── location_repository.go  ⏳ To be created
│   ├── customer_repository.go  ⏳ To be created
│   ├── inventory_repository.go ⏳ To be created
│   ├── purchase_repository.go  ⏳ To be created
│   ├── sale_repository.go      ⏳ To be created
│   ├── transaction_repository.go ⏳ To be created
│   ├── product_repository.go   ⏳ Needs major update
│   └── user_repository.go      ⏳ Needs update
│
├── service/                # Business logic layer
│   ├── category_service.go     ⏳ To be created
│   ├── supplier_service.go     ⏳ To be created
│   ├── location_service.go     ⏳ To be created
│   ├── customer_service.go     ⏳ To be created
│   ├── inventory_service.go    ⏳ To be created (critical - inventory logic)
│   ├── purchase_service.go     ⏳ To be created (critical - inventory updates)
│   ├── sale_service.go         ⏳ To be created (critical - inventory updates)
│   ├── transaction_service.go  ⏳ To be created
│   ├── report_service.go       ⏳ To be created (uses views)
│   ├── product_service.go      ⏳ Needs major update
│   └── user_service.go         ⏳ Needs update
│
└── controller/             # HTTP handlers
    ├── category_controller.go      ⏳ To be created
    ├── supplier_controller.go      ⏳ To be created
    ├── location_controller.go      ⏳ To be created
    ├── customer_controller.go      ⏳ To be created
    ├── inventory_controller.go     ⏳ To be created
    ├── purchase_controller.go      ⏳ To be created
    ├── sale_controller.go          ⏳ To be created
    ├── transaction_controller.go   ⏳ To be created
    ├── report_controller.go        ⏳ To be created
    ├── product_controller.go       ⏳ Needs major update
    └── user_controller.go          ⏳ Needs update
```

---

## New Entities & Their Relationships

### 1. **Categories** (New)
- Hierarchical structure with parent-child relationships
- Required for product classification
- Self-referencing foreign key (ParentCategoryID)

### 2. **Suppliers** (New)
- Vendor/supplier management
- Required for purchase tracking
- Soft delete via IsActive flag

### 3. **Locations** (New)
- Physical storage locations (warehouse, store, shelf, zone)
- Critical for multi-location inventory tracking
- Soft delete via IsActive flag

### 4. **Customers** (New)
- Customer information for sales
- Optional on sales (can sell without customer)
- Loyalty points tracking

### 5. **Inventory** (New - Critical)
- Real-time stock levels by location
- Virtual column: AvailableQuantity = Quantity - ReservedQuantity
- Unique constraint on (ProductID, LocationID)

### 6. **Purchases** & **PurchaseItems** (New)
- Purchase order/receipt tracking
- Master-detail relationship
- Increases inventory on creation

### 7. **Sales** & **SalesItems** (New)
- Sales transaction tracking
- Master-detail relationship
- Decreases inventory on creation
- Updates customer loyalty points

### 8. **Transactions** (New - Audit)
- Complete audit trail for all inventory movements
- Types: purchase, sale, adjustment, transfer, return, damage, theft
- Links to reference documents via ReferenceType/ReferenceID

---

## Critical Business Logic

### Sale Flow (IMPORTANT)
```
1. Receive CreateSaleRequest with items
2. FOR EACH item:
   a. Validate sufficient inventory at specified location
   b. If insufficient → return 422 error
3. Calculate line totals (quantity × unit_price × (1-discount%) × (1+tax_rate%))
4. Calculate total_amount, tax_amount, discount_amount, net_amount
5. BEGIN TRANSACTION
6. Create Sale record
7. FOR EACH item:
   a. Create SaleItem record
   b. UPDATE Inventory SET Quantity = Quantity - item.quantity
      WHERE ProductID = item.product_id AND LocationID = item.location_id
   c. Create Transaction record (type='sale', quantity=negative)
8. IF customer_id provided:
   a. Calculate loyalty points (e.g., 1 point per dollar)
   b. UPDATE Customers SET LoyaltyPoints = LoyaltyPoints + points
9. COMMIT TRANSACTION
10. Return 201 Created with sale_id
```

### Purchase Flow (IMPORTANT)
```
1. Receive CreatePurchaseRequest with items
2. Calculate line totals
3. Calculate total_amount, tax_amount, discount_amount, net_amount
4. BEGIN TRANSACTION
5. Create Purchase record
6. FOR EACH item:
   a. Create PurchaseItem record
   b. UPDATE Inventory SET Quantity = Quantity + item.quantity
      WHERE ProductID = item.product_id AND LocationID = item.location_id
      (INSERT if not exists with ON CONFLICT/UPSERT logic)
   c. Create Transaction record (type='purchase', quantity=positive)
   d. OPTIONAL: Update Product.CostPrice = item.unit_cost (weighted average)
7. COMMIT TRANSACTION
8. Return 201 Created with purchase_id
```

### Inventory Adjustment Flow
```
1. Receive InventoryAdjustmentRequest
2. Get current inventory quantity
3. Calculate new quantity = current + adjustment (can be negative)
4. Validate new quantity >= 0
5. BEGIN TRANSACTION
6. UPDATE Inventory SET Quantity = new_quantity
7. Create Transaction record with PreviousQuantity and NewQuantity
8. COMMIT TRANSACTION
```

### Inventory Transfer Flow
```
1. Receive InventoryTransferRequest
2. Validate from_location has sufficient quantity
3. BEGIN TRANSACTION
4. UPDATE Inventory SET Quantity = Quantity - transfer_qty
   WHERE ProductID = X AND LocationID = from_location
5. UPDATE Inventory SET Quantity = Quantity + transfer_qty
   WHERE ProductID = X AND LocationID = to_location
   (INSERT if not exists)
6. Create Transaction record for source (type='transfer', quantity=negative)
7. Create Transaction record for destination (type='transfer', quantity=positive)
8. COMMIT TRANSACTION
```

---

## Database Views Implemented

### vw_product_stock_summary
```sql
-- Returns: ProductID, Name, SKU, Category, CostPrice, SellingPrice,
--          TotalStock, ReservedStock, AvailableStock, StockStatus
```
**API Endpoint:** `GET /reports/stock-summary`

### vw_inventory_by_location
```sql
-- Returns: Location, LocationType, Product, SKU, Category,
--          Quantity, ReservedQuantity, AvailableQuantity, LastUpdated
```
**API Endpoint:** `GET /reports/inventory-by-location`

---

## API Endpoints Summary

| Entity | Create | List | Get | Update | Delete | Special |
|--------|--------|------|-----|--------|--------|---------|
| **Categories** | POST /categories | GET /categories | GET /categories/:id | PUT /categories/:id | DELETE /categories/:id | GET /categories/:id/children |
| **Suppliers** | POST /suppliers | GET /suppliers | GET /suppliers/:id | PUT /suppliers/:id | DELETE /suppliers/:id (soft) | |
| **Locations** | POST /locations | GET /locations | GET /locations/:id | PUT /locations/:id | DELETE /locations/:id (soft) | |
| **Customers** | POST /customers | GET /customers | GET /customers/:id | PUT /customers/:id | DELETE /customers/:id (soft) | GET /customers/:id/sales, PUT /customers/:id/loyalty |
| **Products** | POST /products | GET /products | GET /products/:id | PUT /products/:id | DELETE /products/:id (soft) | GET /products/sku/:sku, GET /products/barcode/:barcode, GET /products/low-stock, GET /products/expiring |
| **Inventory** | - | GET /inventory | GET /inventory/product/:id | - | - | POST /inventory/adjust, POST /inventory/transfer, GET /inventory/location/:id, GET /inventory/view/by-location |
| **Purchases** | POST /purchases | GET /purchases | GET /purchases/:id | - | DELETE /purchases/:id | PUT /purchases/:id/payment-status |
| **Sales** | POST /sales | GET /sales | GET /sales/:id | - | DELETE /sales/:id | PUT /sales/:id/payment-status, POST /sales/:id/refund |
| **Transactions** | - | GET /transactions | - | - | - | GET /transactions/product/:id, GET /transactions/location/:id, GET /transactions/user/:id |
| **Reports** | - | - | - | - | - | GET /reports/stock-summary, GET /reports/inventory-by-location, GET /reports/sales-summary, GET /reports/purchase-summary, GET /reports/low-stock-alerts, GET /reports/expiring-products |

**Total New Endpoints:** ~60+

---

## Schema Changes from Original

### Products Table - Major Changes
**Old Fields:**
- id, name, price, quantity, category (string)

**New Fields:**
- ProductID, Name, Description, SKU (unique), Barcode (unique), CategoryID (FK)
- BatchNo, ExpiryDate
- CostPrice, SellingPrice (separate)
- MinStockLevel, MaxStockLevel, ReorderPoint
- Unit (enum), IsActive, CreatedAt, UpdatedAt

**Migration Impact:**
- `quantity` moved to Inventory table (per location)
- `category` string replaced with CategoryID foreign key
- `price` split into CostPrice and SellingPrice

### Users Table - Changes
**Added:**
- Name, Email, IsActive, CreatedAt, UpdatedAt, LastLogin
- Role now enum: admin, manager, staff, viewer

---

## Implementation Priority

### Phase 1: Foundation (Critical)
1. ✅ Domain models for all entities
2. ✅ DTOs for all operations
3. ✅ Request models for key operations
4. ⏳ Response models
5. ⏳ Category, Supplier, Location, Customer repositories (simple CRUD)

### Phase 2: Core Inventory (Critical)
1. ⏳ Inventory repository with complex queries
2. ⏳ Inventory service with adjustment/transfer logic
3. ⏳ Transaction repository (audit trail)
4. ⏳ Updated Product repository (new fields, SKU/barcode lookups)

### Phase 3: Purchases & Sales (Critical)
1. ⏳ Purchase repository & service (with inventory updates)
2. ⏳ Sale repository & service (with inventory updates & validation)
3. ⏳ Integration of Transaction logging in both
4. ⏳ Controllers for purchases & sales

### Phase 4: APIs & Controllers
1. ⏳ Controllers for categories, suppliers, locations, customers
2. ⏳ Updated product controller with new endpoints
3. ⏳ Inventory controller (adjust, transfer)
4. ⏳ Purchase & sale controllers

### Phase 5: Reports
1. ⏳ Report service (query views)
2. ⏳ Report controller
3. ⏳ Custom report endpoints

### Phase 6: Testing & Refinement
1. ⏳ Integration testing
2. ⏳ Database migrations
3. ⏳ API documentation (Swagger)
4. ⏳ Error handling improvements

---

## Files Created/Modified So Far

### ✅ Created (New Files)
- `pkg/domain/category.go`
- `pkg/domain/supplier.go`
- `pkg/domain/location.go`
- `pkg/domain/customer.go`
- `pkg/domain/inventory.go`
- `pkg/domain/purchase.go`
- `pkg/domain/sale.go`
- `pkg/domain/transaction.go`
- `pkg/controller/request/category_request.go`
- `pkg/controller/request/supplier_request.go`
- `pkg/controller/request/purchase_request.go`
- `pkg/controller/request/sale_request.go`
- `pkg/controller/request/inventory_request.go`
- `API_SPECIFICATION.md` (complete API documentation)
- `IMPLEMENTATION_PLAN.md` (this file)

### ✅ Modified (Updated Files)
- `pkg/domain/product.go` - Major restructure to match new schema
- `pkg/domain/user.go` - Added new fields
- `pkg/service/dto/model.go` - Added all new DTOs

### ⏳ To Be Modified
- `pkg/controller/request/product_request.go` - Update for new fields
- `pkg/controller/request/user_request.go` - Update for new fields
- `pkg/controller/response/productResponse.go` - Update for new fields
- `pkg/controller/response/user_response.go` - Update for new fields
- `pkg/repository/product_repository.go` - Complete rewrite
- `pkg/repository/user_repository.go` - Update queries
- `pkg/service/product_service.go` - Update logic
- `pkg/service/user_service.go` - Update logic
- `pkg/controller/product_controller.go` - Add new endpoints
- `pkg/controller/user_controller.go` - Update logic
- `cmd/imsapi/main.go` - Register all new routes

---

## Next Steps

1. **Create response models** for all new entities
2. **Implement repositories** starting with simple CRUD (categories, suppliers, locations, customers)
3. **Implement critical repositories** (inventory, transactions) with complex logic
4. **Implement services** with business logic, especially:
   - SaleService (inventory reduction logic)
   - PurchaseService (inventory increase logic)
   - InventoryService (adjustment & transfer logic)
5. **Implement controllers** and register routes
6. **Update existing** product/user repositories, services, controllers
7. **Test the complete flow** with database

---

## Testing Strategy

### Critical Test Scenarios

1. **Sale with inventory reduction:**
   - Create product with inventory
   - Make a sale
   - Verify inventory reduced
   - Verify transaction created

2. **Purchase with inventory increase:**
   - Create purchase
   - Verify inventory increased
   - Verify transaction created

3. **Insufficient inventory:**
   - Attempt sale with quantity > available
   - Should return 422 error
   - Inventory should not change

4. **Inventory transfer:**
   - Transfer between locations
   - Verify source reduced, destination increased
   - Verify two transactions created

5. **Refund:**
   - Create sale
   - Process refund
   - Verify inventory restored
   - Verify transactions logged

---

## Notes

- All monetary calculations should use `float64` (consider decimal library for production)
- Soft deletes used for: Suppliers, Locations, Customers, Products
- Hard deletes with CASCADE for: SalesItems, PurchaseItems
- Triggers automatically update UpdatedAt timestamps
- Use database transactions for all multi-step operations
- Complete audit trail via Transactions table
- Views optimize common reporting queries

---

## References

- Schema: `schema_sqlite_extended.sql`
- API Spec: `API_SPECIFICATION.md`
- Current project structure: `server/pkg/`
