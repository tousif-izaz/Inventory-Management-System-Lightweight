# Database Schema Documentation
## Inventory Management System - SQLite

**Schema Version:** 2.0.0
**Last Updated:** 2025-10-15
**Database Engine:** SQLite 3.x

---

## Table of Contents
1. [Overview](#overview)
2. [Entity Relationship](#entity-relationship)
3. [Table Specifications](#table-specifications)
4. [Indexes & Performance](#indexes--performance)
5. [Triggers & Automation](#triggers--automation)
6. [Views](#views)
7. [Data Types & Constraints](#data-types--constraints)

---

## Overview

This database schema supports a comprehensive inventory management system with the following capabilities:

- ✅ **Product Management**: Catalog with SKU, barcodes, categories, batch tracking, and expiry dates
- ✅ **Multi-Location Inventory**: Track stock across multiple physical locations
- ✅ **Purchase Management**: Complete purchase order processing with supplier tracking
- ✅ **Sales Management**: Point-of-sale transactions with optional customer tracking
- ✅ **Audit Trail**: Full transaction history for compliance and analysis
- ✅ **User Management**: Role-based access control
- ✅ **Automated Triggers**: Timestamp updates and data validation
- ✅ **Reporting Views**: Pre-built views for common queries

---

## Entity Relationship

```
Users ────┐
          │
          ├──> Purchases ──> PurchaseItems ──> Products ──> Categories
          │                                        │
          ├──> Sales ──> SalesItems ─────────────┘
          │                                        │
          └──> Transactions ────────────────────────┘
                    │                               │
                    └──> Locations <─── Inventory <─┘
                              │
Suppliers ────────────────────┘

Customers (optional) ──> Sales
```

---

## Table Specifications

### Core Tables

#### 1. **Users**
Manages system users with role-based access control.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **UserID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique user identifier |
| Username | TEXT | UNIQUE, NOT NULL | Login username |
| Name | TEXT | NOT NULL | Full name |
| Email | TEXT | UNIQUE | Email address |
| PasswordHash | TEXT | NOT NULL | Bcrypt hashed password |
| Role | TEXT | NOT NULL, CHECK | User role (admin, manager, staff, viewer) |
| IsActive | INTEGER | DEFAULT 1, CHECK | Active status (0=inactive, 1=active) |
| CreatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Account creation date |
| UpdatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Last modification date |
| LastLogin | DATETIME | | Last login timestamp |

**Roles:**
- `admin`: Full system access
- `manager`: Manage inventory, users, reports
- `staff`: Process sales/purchases
- `viewer`: Read-only access

---

#### 2. **Categories**
Product classification with hierarchical support.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **CategoryID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique category identifier |
| Name | TEXT | UNIQUE, NOT NULL | Category name |
| Description | TEXT | | Category description |
| ParentCategoryID | INTEGER | FOREIGN KEY | Parent category for hierarchy |
| CreatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Creation date |

**Example Hierarchy:**
```
Fragrances (parent)
├── Men's Fragrances (child)
├── Women's Fragrances (child)
└── Unisex Fragrances (child)
```

---

#### 3. **Products**
Core product catalog with comprehensive product information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **ProductID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique product identifier |
| Name | TEXT | NOT NULL | Product name |
| Description | TEXT | | Detailed description |
| SKU | TEXT | UNIQUE, NOT NULL | Stock Keeping Unit |
| Barcode | TEXT | UNIQUE | Product barcode/UPC |
| CategoryID | INTEGER | NOT NULL, FOREIGN KEY | Product category |
| BatchNo | TEXT | | Manufacturing batch number |
| ExpiryDate | DATE | | Product expiration date |
| CostPrice | REAL | NOT NULL, CHECK ≥ 0 | Purchase/cost price |
| SellingPrice | REAL | NOT NULL, CHECK ≥ 0 | Retail selling price |
| MinStockLevel | INTEGER | DEFAULT 0 | Minimum stock threshold |
| MaxStockLevel | INTEGER | | Maximum stock capacity |
| ReorderPoint | INTEGER | DEFAULT 0 | Reorder trigger level |
| Unit | TEXT | DEFAULT 'pcs', CHECK | Unit of measure |
| IsActive | INTEGER | DEFAULT 1, CHECK | Active status |
| CreatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Creation date |
| UpdatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Last update |

**Units of Measure:**
- `pcs`: Pieces (default)
- `kg`: Kilograms
- `liter`: Liters
- `box`: Boxes
- `carton`: Cartons

---

#### 4. **Suppliers**
Vendor and supplier information for procurement.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **SupplierID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique supplier identifier |
| Name | TEXT | UNIQUE, NOT NULL | Supplier company name |
| ContactPerson | TEXT | | Primary contact name |
| Email | TEXT | | Email address |
| Phone | TEXT | | Phone number |
| Address | TEXT | | Street address |
| City | TEXT | | City |
| Country | TEXT | | Country |
| TaxID | TEXT | | Tax/VAT ID number |
| PaymentTerms | TEXT | | Payment terms (e.g., Net 30) |
| IsActive | INTEGER | DEFAULT 1, CHECK | Active status |
| CreatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Creation date |
| UpdatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Last update |

---

#### 5. **Locations**
Physical storage locations for inventory tracking.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **LocationID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique location identifier |
| Name | TEXT | UNIQUE, NOT NULL | Location name |
| Type | TEXT | CHECK | Location type |
| Address | TEXT | | Physical address |
| IsActive | INTEGER | DEFAULT 1, CHECK | Active status |
| CreatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Creation date |

**Location Types:**
- `warehouse`: Main storage facility
- `store`: Retail store
- `shelf`: Specific shelf/bin
- `zone`: Storage zone/area

---

#### 6. **Customers** (Optional)
Customer information for sales tracking and loyalty programs.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **CustomerID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique customer identifier |
| Name | TEXT | NOT NULL | Customer name |
| Email | TEXT | UNIQUE | Email address |
| Phone | TEXT | | Phone number |
| Address | TEXT | | Address |
| City | TEXT | | City |
| Country | TEXT | | Country |
| LoyaltyPoints | INTEGER | DEFAULT 0 | Accumulated loyalty points |
| IsActive | INTEGER | DEFAULT 1, CHECK | Active status |
| CreatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Registration date |
| UpdatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Last update |

---

### Inventory Management

#### 7. **Inventory**
Real-time stock levels by product and location.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **InventoryID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique inventory record ID |
| ProductID | INTEGER | NOT NULL, FOREIGN KEY | Product reference |
| LocationID | INTEGER | NOT NULL, FOREIGN KEY | Storage location |
| Quantity | INTEGER | NOT NULL, DEFAULT 0, CHECK ≥ 0 | Total quantity |
| ReservedQuantity | INTEGER | DEFAULT 0, CHECK ≥ 0 | Reserved/allocated quantity |
| AvailableQuantity | INTEGER | VIRTUAL COMPUTED | Quantity - ReservedQuantity |
| LastUpdated | DATETIME | DEFAULT CURRENT_TIMESTAMP | Last stock update |

**Unique Constraint:** (ProductID, LocationID) - One record per product per location

---

### Purchase Management

#### 8. **Purchases**
Purchase order/receipt header information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **PurchaseID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique purchase ID |
| PurchaseDate | DATE | NOT NULL, DEFAULT CURRENT_DATE | Purchase date |
| SupplierID | INTEGER | NOT NULL, FOREIGN KEY | Supplier reference |
| InvoiceNumber | TEXT | UNIQUE, NOT NULL | Supplier invoice number |
| ReferenceNumber | TEXT | | Internal reference |
| TotalAmount | REAL | NOT NULL, DEFAULT 0, CHECK ≥ 0 | Subtotal amount |
| TaxAmount | REAL | DEFAULT 0, CHECK ≥ 0 | Total tax |
| DiscountAmount | REAL | DEFAULT 0, CHECK ≥ 0 | Total discount |
| NetAmount | REAL | NOT NULL, DEFAULT 0, CHECK ≥ 0 | Final amount |
| PaymentStatus | TEXT | DEFAULT 'pending', CHECK | Payment status |
| PaymentMethod | TEXT | | Payment method |
| Notes | TEXT | | Additional notes |
| ReceivedBy | INTEGER | FOREIGN KEY | User who received |
| CreatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Record creation |

**Payment Statuses:**
- `pending`: Payment not yet made
- `partial`: Partially paid
- `paid`: Fully paid

---

#### 9. **PurchaseItems**
Line items for purchases (for COGS calculation).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **PurchaseItemID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique line item ID |
| PurchaseID | INTEGER | NOT NULL, FOREIGN KEY CASCADE | Purchase reference |
| ProductID | INTEGER | NOT NULL, FOREIGN KEY | Product reference |
| Quantity | INTEGER | NOT NULL, CHECK > 0 | Quantity purchased |
| UnitCost | REAL | NOT NULL, CHECK ≥ 0 | Cost per unit |
| TaxRate | REAL | DEFAULT 0, CHECK 0-100 | Tax percentage |
| DiscountPercent | REAL | DEFAULT 0, CHECK 0-100 | Discount percentage |
| LineTotal | REAL | NOT NULL, CHECK ≥ 0 | Line total amount |
| LocationID | INTEGER | FOREIGN KEY | Receiving location |
| BatchNo | TEXT | | Batch number |
| ExpiryDate | DATE | | Expiry date |

---

### Sales Management

#### 10. **Sales**
Sales transaction header.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **SaleID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique sale ID |
| SaleDate | DATE | NOT NULL, DEFAULT CURRENT_DATE | Sale date |
| ReceiptNo | TEXT | UNIQUE, NOT NULL | Receipt/invoice number |
| CustomerID | INTEGER | FOREIGN KEY | Customer reference (optional) |
| TotalAmount | REAL | NOT NULL, DEFAULT 0, CHECK ≥ 0 | Subtotal amount |
| TaxAmount | REAL | DEFAULT 0, CHECK ≥ 0 | Total tax |
| DiscountAmount | REAL | DEFAULT 0, CHECK ≥ 0 | Total discount |
| NetAmount | REAL | NOT NULL, CHECK ≥ 0 | Final amount |
| PaymentStatus | TEXT | DEFAULT 'paid', CHECK | Payment status |
| PaymentMethod | TEXT | CHECK | Payment method |
| Notes | TEXT | | Additional notes |
| SoldBy | INTEGER | FOREIGN KEY | Sales person |
| CreatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Record creation |

**Payment Methods:**
- `cash`: Cash payment
- `card`: Credit/debit card
- `mobile`: Mobile payment
- `bank_transfer`: Bank transfer
- `credit`: Store credit

**Payment Statuses:**
- `pending`: Not yet paid
- `partial`: Partially paid
- `paid`: Fully paid
- `refunded`: Refunded

---

#### 11. **SalesItems**
Line items for sales transactions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **SaleItemID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique line item ID |
| SaleID | INTEGER | NOT NULL, FOREIGN KEY CASCADE | Sale reference |
| ProductID | INTEGER | NOT NULL, FOREIGN KEY | Product reference |
| Quantity | INTEGER | NOT NULL, CHECK > 0 | Quantity sold |
| UnitPrice | REAL | NOT NULL, CHECK ≥ 0 | Price per unit |
| TaxRate | REAL | DEFAULT 0, CHECK 0-100 | Tax percentage |
| DiscountPercent | REAL | DEFAULT 0, CHECK 0-100 | Discount percentage |
| LineTotal | REAL | NOT NULL, CHECK ≥ 0 | Line total amount |
| LocationID | INTEGER | FOREIGN KEY | Source location |

---

### Audit & Transactions

#### 12. **Transactions**
Complete audit trail for all inventory movements.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| **TransactionID** | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique transaction ID |
| TransactionType | TEXT | NOT NULL, CHECK | Type of transaction |
| TransactionDate | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Transaction date/time |
| ProductID | INTEGER | NOT NULL, FOREIGN KEY | Product reference |
| LocationID | INTEGER | FOREIGN KEY | Location reference |
| Quantity | INTEGER | NOT NULL | Quantity change (+ or -) |
| ReferenceType | TEXT | | Reference type (e.g., 'Sale', 'Purchase') |
| ReferenceID | INTEGER | | Reference ID |
| PreviousQuantity | INTEGER | | Stock before transaction |
| NewQuantity | INTEGER | | Stock after transaction |
| Reason | TEXT | | Reason for transaction |
| PerformedBy | INTEGER | FOREIGN KEY | User who performed |
| Notes | TEXT | | Additional notes |
| CreatedAt | DATETIME | DEFAULT CURRENT_TIMESTAMP | Record creation |

**Transaction Types:**
- `purchase`: Stock received from supplier
- `sale`: Stock sold to customer
- `adjustment`: Manual stock adjustment
- `transfer`: Stock moved between locations
- `return`: Product returned
- `damage`: Damaged goods write-off
- `theft`: Theft/loss write-off

---

## Indexes & Performance

### Automatic Indexes
- Primary keys automatically indexed
- Unique constraints automatically indexed

### Custom Indexes

| Index Name | Table | Columns | Purpose |
|------------|-------|---------|---------|
| idx_users_username | Users | Username | Fast login lookup |
| idx_users_email | Users | Email | Email validation |
| idx_users_role | Users | Role | Role-based queries |
| idx_products_sku | Products | SKU | Product lookup |
| idx_products_barcode | Products | Barcode | Barcode scanning |
| idx_products_category | Products | CategoryID | Category filtering |
| idx_products_name | Products | Name | Name search |
| idx_products_active | Products | IsActive | Active product queries |
| idx_inventory_product | Inventory | ProductID | Product stock lookup |
| idx_inventory_location | Inventory | LocationID | Location inventory |
| idx_purchases_date | Purchases | PurchaseDate | Date range queries |
| idx_purchases_supplier | Purchases | SupplierID | Supplier reports |
| idx_sales_date | Sales | SaleDate | Sales reports |
| idx_sales_customer | Sales | CustomerID | Customer history |
| idx_transactions_type | Transactions | TransactionType | Transaction filtering |
| idx_transactions_date | Transactions | TransactionDate | Audit reports |

---

## Triggers & Automation

### Update Timestamp Triggers

Automatically update `UpdatedAt` timestamp on record modification:

1. **trg_products_update_timestamp**: Updates Products.UpdatedAt
2. **trg_users_update_timestamp**: Updates Users.UpdatedAt
3. **trg_suppliers_update_timestamp**: Updates Suppliers.UpdatedAt
4. **trg_customers_update_timestamp**: Updates Customers.UpdatedAt
5. **trg_inventory_update_timestamp**: Updates Inventory.LastUpdated

---

## Views

### 1. vw_product_stock_summary
Comprehensive product stock overview with status indicators.

**Columns:**
- ProductID, Name, SKU, Barcode, Category
- CostPrice, SellingPrice
- TotalStock, ReservedStock, AvailableStock
- MinStockLevel, ReorderPoint
- StockStatus (OK, Low Stock, Critical)

**Usage:**
```sql
SELECT * FROM vw_product_stock_summary WHERE StockStatus = 'Critical';
```

---

### 2. vw_inventory_by_location
Inventory breakdown by physical location.

**Columns:**
- Location, LocationType
- Product, SKU, Category
- Quantity, ReservedQuantity, AvailableQuantity
- LastUpdated

**Usage:**
```sql
SELECT * FROM vw_inventory_by_location WHERE Location = 'Main Warehouse';
```

---

## Data Types & Constraints

### SQLite Data Types Used

| Type | Usage | Example |
|------|-------|---------|
| INTEGER | IDs, quantities, counts | ProductID, Quantity |
| REAL | Prices, amounts, rates | CostPrice, TaxRate |
| TEXT | Strings, names, descriptions | Name, SKU, Email |
| DATE | Date values | PurchaseDate, ExpiryDate |
| DATETIME | Timestamps | CreatedAt, UpdatedAt |

### Constraint Types

**CHECK Constraints:**
- Ensure data validity (e.g., `Quantity >= 0`)
- Enum-like constraints (e.g., `Role IN ('admin', 'manager')`)

**FOREIGN KEY Constraints:**
- Maintain referential integrity
- Actions: `ON DELETE RESTRICT`, `ON DELETE CASCADE`, `ON DELETE SET NULL`

**UNIQUE Constraints:**
- Prevent duplicates (e.g., Username, SKU, InvoiceNumber)

**NOT NULL Constraints:**
- Enforce required fields

---

## Schema Version Control

The `SchemaVersion` table tracks schema migrations:

| Version | Description | Applied Date |
|---------|-------------|--------------|
| 2.0.0 | Extended schema with comprehensive inventory management | 2025-10-15 |

---

## Notes

1. **Foreign Keys**: Enabled via `PRAGMA foreign_keys=ON`
2. **WAL Mode**: Enabled for better concurrency via `PRAGMA journal_mode=WAL`
3. **Backward Compatibility**: Falls back to basic schema if extended file not found
4. **Virtual Columns**: `AvailableQuantity` calculated automatically
5. **Cascading Deletes**: Line items deleted when parent record deleted

---

## Quick Reference

### Primary Keys
- All tables use INTEGER AUTOINCREMENT primary keys
- Naming pattern: `[TableName]ID` (e.g., ProductID, UserID)

### Timestamps
- `CreatedAt`: Record creation time (automatic)
- `UpdatedAt`: Last modification time (trigger-updated)
- `LastUpdated`: Specific last update time

### Status Fields
- `IsActive`: Boolean (0=inactive, 1=active)
- `PaymentStatus`: Text-based status tracking

---

**For implementation details, see:** `schema_sqlite_extended.sql`
