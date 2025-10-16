-- ============================================================================
-- Inventory Management System - Extended Database Schema (SQLite)
-- ============================================================================
-- This schema supports comprehensive inventory management including:
-- - Product catalog with categories, SKU, barcodes, and batch tracking
-- - Multi-location inventory tracking
-- - Purchase management with supplier information
-- - Sales tracking with optional customer information
-- - Full transaction audit trail
-- - User management with role-based access
-- ============================================================================

-- Drop existing tables in reverse dependency order
DROP TABLE IF EXISTS SalesItems;
DROP TABLE IF EXISTS Sales;
DROP TABLE IF EXISTS PurchaseItems;
DROP TABLE IF EXISTS Purchases;
DROP TABLE IF EXISTS Transactions;
DROP TABLE IF EXISTS Inventory;
DROP TABLE IF EXISTS Products;
DROP TABLE IF EXISTS Categories;
DROP TABLE IF EXISTS Suppliers;
DROP TABLE IF EXISTS Locations;
DROP TABLE IF EXISTS Customers;
DROP TABLE IF EXISTS Users;

-- ============================================================================
-- CORE TABLES
-- ============================================================================

-- Users Table
-- Stores user information for authentication and authorization
CREATE TABLE Users (
    UserID INTEGER PRIMARY KEY AUTOINCREMENT,
    Username TEXT UNIQUE NOT NULL,
    Name TEXT NOT NULL,
    Email TEXT UNIQUE,
    PasswordHash TEXT NOT NULL,
    Role TEXT NOT NULL CHECK(Role IN ('admin', 'manager', 'staff', 'viewer')),
    IsActive INTEGER DEFAULT 1 CHECK(IsActive IN (0, 1)),
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    LastLogin DATETIME
);

-- Categories Table
-- Product categorization for organization and reporting
CREATE TABLE Categories (
    CategoryID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT UNIQUE NOT NULL,
    Description TEXT,
    ParentCategoryID INTEGER,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ParentCategoryID) REFERENCES Categories(CategoryID) ON DELETE SET NULL
);

-- Suppliers Table
-- Vendor/supplier information for purchase tracking
CREATE TABLE Suppliers (
    SupplierID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT UNIQUE NOT NULL,
    ContactPerson TEXT,
    Email TEXT,
    Phone TEXT,
    Address TEXT,
    City TEXT,
    Country TEXT,
    TaxID TEXT,
    PaymentTerms TEXT,
    IsActive INTEGER DEFAULT 1 CHECK(IsActive IN (0, 1)),
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Locations Table
-- Physical locations where inventory is stored
CREATE TABLE Locations (
    LocationID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT UNIQUE NOT NULL,
    Type TEXT CHECK(Type IN ('warehouse', 'store', 'shelf', 'zone')),
    Address TEXT,
    IsActive INTEGER DEFAULT 1 CHECK(IsActive IN (0, 1)),
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Customers Table (Optional)
-- Customer information for sales tracking
CREATE TABLE Customers (
    CustomerID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL,
    Email TEXT UNIQUE,
    Phone TEXT,
    Address TEXT,
    City TEXT,
    Country TEXT,
    LoyaltyPoints INTEGER DEFAULT 0,
    IsActive INTEGER DEFAULT 1 CHECK(IsActive IN (0, 1)),
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Products Table
-- Core product catalog with batch and expiry tracking
CREATE TABLE Products (
    ProductID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL,
    Description TEXT,
    SKU TEXT UNIQUE NOT NULL,
    Barcode TEXT UNIQUE,
    CategoryID INTEGER NOT NULL,
    BatchNo TEXT,
    ExpiryDate DATE,
    CostPrice REAL NOT NULL CHECK(CostPrice >= 0),
    SellingPrice REAL NOT NULL CHECK(SellingPrice >= 0),
    MinStockLevel INTEGER DEFAULT 0,
    MaxStockLevel INTEGER,
    ReorderPoint INTEGER DEFAULT 0,
    Unit TEXT DEFAULT 'pcs' CHECK(Unit IN ('pcs', 'kg', 'liter', 'box', 'carton')),
    IsActive INTEGER DEFAULT 1 CHECK(IsActive IN (0, 1)),
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (CategoryID) REFERENCES Categories(CategoryID) ON DELETE RESTRICT
);

-- ============================================================================
-- INVENTORY TRACKING
-- ============================================================================

-- Inventory Table
-- Real-time inventory levels by location
CREATE TABLE Inventory (
    InventoryID INTEGER PRIMARY KEY AUTOINCREMENT,
    ProductID INTEGER NOT NULL,
    LocationID INTEGER NOT NULL,
    Quantity INTEGER NOT NULL DEFAULT 0 CHECK(Quantity >= 0),
    ReservedQuantity INTEGER DEFAULT 0 CHECK(ReservedQuantity >= 0),
    AvailableQuantity INTEGER GENERATED ALWAYS AS (Quantity - ReservedQuantity) VIRTUAL,
    LastUpdated DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ProductID) REFERENCES Products(ProductID) ON DELETE RESTRICT,
    FOREIGN KEY (LocationID) REFERENCES Locations(LocationID) ON DELETE RESTRICT,
    UNIQUE(ProductID, LocationID)
);

-- ============================================================================
-- PURCHASE MANAGEMENT
-- ============================================================================

-- Purchases Table
-- Purchase order/receipt header
CREATE TABLE Purchases (
    PurchaseID INTEGER PRIMARY KEY AUTOINCREMENT,
    PurchaseDate DATE NOT NULL DEFAULT CURRENT_DATE,
    SupplierID INTEGER NOT NULL,
    InvoiceNumber TEXT UNIQUE NOT NULL,
    ReferenceNumber TEXT,
    TotalAmount REAL NOT NULL DEFAULT 0 CHECK(TotalAmount >= 0),
    TaxAmount REAL DEFAULT 0 CHECK(TaxAmount >= 0),
    DiscountAmount REAL DEFAULT 0 CHECK(DiscountAmount >= 0),
    NetAmount REAL NOT NULL DEFAULT 0 CHECK(NetAmount >= 0),
    PaymentStatus TEXT DEFAULT 'pending' CHECK(PaymentStatus IN ('pending', 'partial', 'paid')),
    PaymentMethod TEXT,
    Notes TEXT,
    ReceivedBy INTEGER,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (SupplierID) REFERENCES Suppliers(SupplierID) ON DELETE RESTRICT,
    FOREIGN KEY (ReceivedBy) REFERENCES Users(UserID) ON DELETE SET NULL
);

-- PurchaseItems Table
-- Line items for purchases (for COGS calculation)
CREATE TABLE PurchaseItems (
    PurchaseItemID INTEGER PRIMARY KEY AUTOINCREMENT,
    PurchaseID INTEGER NOT NULL,
    ProductID INTEGER NOT NULL,
    Quantity INTEGER NOT NULL CHECK(Quantity > 0),
    UnitCost REAL NOT NULL CHECK(UnitCost >= 0),
    TaxRate REAL DEFAULT 0 CHECK(TaxRate >= 0 AND TaxRate <= 100),
    DiscountPercent REAL DEFAULT 0 CHECK(DiscountPercent >= 0 AND DiscountPercent <= 100),
    LineTotal REAL NOT NULL CHECK(LineTotal >= 0),
    LocationID INTEGER,
    BatchNo TEXT,
    ExpiryDate DATE,
    FOREIGN KEY (PurchaseID) REFERENCES Purchases(PurchaseID) ON DELETE CASCADE,
    FOREIGN KEY (ProductID) REFERENCES Products(ProductID) ON DELETE RESTRICT,
    FOREIGN KEY (LocationID) REFERENCES Locations(LocationID) ON DELETE SET NULL
);

-- ============================================================================
-- SALES MANAGEMENT
-- ============================================================================

-- Sales Table
-- Sales transaction header
CREATE TABLE Sales (
    SaleID INTEGER PRIMARY KEY AUTOINCREMENT,
    SaleDate DATE NOT NULL DEFAULT CURRENT_DATE,
    ReceiptNo TEXT UNIQUE NOT NULL,
    CustomerID INTEGER,
    TotalAmount REAL NOT NULL DEFAULT 0 CHECK(TotalAmount >= 0),
    TaxAmount REAL DEFAULT 0 CHECK(TaxAmount >= 0),
    DiscountAmount REAL DEFAULT 0 CHECK(DiscountAmount >= 0),
    NetAmount REAL NOT NULL DEFAULT 0 CHECK(NetAmount >= 0),
    PaymentStatus TEXT DEFAULT 'paid' CHECK(PaymentStatus IN ('pending', 'partial', 'paid', 'refunded')),
    PaymentMethod TEXT CHECK(PaymentMethod IN ('cash', 'card', 'mobile', 'bank_transfer', 'credit')),
    Notes TEXT,
    SoldBy INTEGER,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (CustomerID) REFERENCES Customers(CustomerID) ON DELETE SET NULL,
    FOREIGN KEY (SoldBy) REFERENCES Users(UserID) ON DELETE SET NULL
);

-- SalesItems Table
-- Line items for sales
CREATE TABLE SalesItems (
    SaleItemID INTEGER PRIMARY KEY AUTOINCREMENT,
    SaleID INTEGER NOT NULL,
    ProductID INTEGER NOT NULL,
    Quantity INTEGER NOT NULL CHECK(Quantity > 0),
    UnitPrice REAL NOT NULL CHECK(UnitPrice >= 0),
    TaxRate REAL DEFAULT 0 CHECK(TaxRate >= 0 AND TaxRate <= 100),
    DiscountPercent REAL DEFAULT 0 CHECK(DiscountPercent >= 0 AND DiscountPercent <= 100),
    LineTotal REAL NOT NULL CHECK(LineTotal >= 0),
    LocationID INTEGER,
    FOREIGN KEY (SaleID) REFERENCES Sales(SaleID) ON DELETE CASCADE,
    FOREIGN KEY (ProductID) REFERENCES Products(ProductID) ON DELETE RESTRICT,
    FOREIGN KEY (LocationID) REFERENCES Locations(LocationID) ON DELETE SET NULL
);

-- ============================================================================
-- AUDIT & TRANSACTIONS
-- ============================================================================

-- Transactions Table
-- Audit trail for all inventory movements
CREATE TABLE Transactions (
    TransactionID INTEGER PRIMARY KEY AUTOINCREMENT,
    TransactionType TEXT NOT NULL CHECK(TransactionType IN ('purchase', 'sale', 'adjustment', 'transfer', 'return', 'damage', 'theft')),
    TransactionDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ProductID INTEGER NOT NULL,
    LocationID INTEGER,
    Quantity INTEGER NOT NULL,
    ReferenceType TEXT,
    ReferenceID INTEGER,
    PreviousQuantity INTEGER,
    NewQuantity INTEGER,
    Reason TEXT,
    PerformedBy INTEGER,
    Notes TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ProductID) REFERENCES Products(ProductID) ON DELETE RESTRICT,
    FOREIGN KEY (LocationID) REFERENCES Locations(LocationID) ON DELETE SET NULL,
    FOREIGN KEY (PerformedBy) REFERENCES Users(UserID) ON DELETE SET NULL
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

-- User indexes
CREATE INDEX idx_users_username ON Users(Username);
CREATE INDEX idx_users_email ON Users(Email);
CREATE INDEX idx_users_role ON Users(Role);

-- Product indexes
CREATE INDEX idx_products_sku ON Products(SKU);
CREATE INDEX idx_products_barcode ON Products(Barcode);
CREATE INDEX idx_products_category ON Products(CategoryID);
CREATE INDEX idx_products_name ON Products(Name);
CREATE INDEX idx_products_active ON Products(IsActive);

-- Inventory indexes
CREATE INDEX idx_inventory_product ON Inventory(ProductID);
CREATE INDEX idx_inventory_location ON Inventory(LocationID);
CREATE INDEX idx_inventory_quantity ON Inventory(Quantity);

-- Purchase indexes
CREATE INDEX idx_purchases_date ON Purchases(PurchaseDate);
CREATE INDEX idx_purchases_supplier ON Purchases(SupplierID);
CREATE INDEX idx_purchases_invoice ON Purchases(InvoiceNumber);
CREATE INDEX idx_purchase_items_purchase ON PurchaseItems(PurchaseID);
CREATE INDEX idx_purchase_items_product ON PurchaseItems(ProductID);

-- Sales indexes
CREATE INDEX idx_sales_date ON Sales(SaleDate);
CREATE INDEX idx_sales_customer ON Sales(CustomerID);
CREATE INDEX idx_sales_receipt ON Sales(ReceiptNo);
CREATE INDEX idx_sale_items_sale ON SalesItems(SaleID);
CREATE INDEX idx_sale_items_product ON SalesItems(ProductID);

-- Transaction indexes
CREATE INDEX idx_transactions_type ON Transactions(TransactionType);
CREATE INDEX idx_transactions_date ON Transactions(TransactionDate);
CREATE INDEX idx_transactions_product ON Transactions(ProductID);
CREATE INDEX idx_transactions_location ON Transactions(LocationID);

-- Category indexes
CREATE INDEX idx_categories_name ON Categories(Name);
CREATE INDEX idx_categories_parent ON Categories(ParentCategoryID);

-- Supplier indexes
CREATE INDEX idx_suppliers_name ON Suppliers(Name);
CREATE INDEX idx_suppliers_active ON Suppliers(IsActive);

-- Location indexes
CREATE INDEX idx_locations_name ON Locations(Name);
CREATE INDEX idx_locations_type ON Locations(Type);

-- Customer indexes
CREATE INDEX idx_customers_name ON Customers(Name);
CREATE INDEX idx_customers_email ON Customers(Email);

-- ============================================================================
-- TRIGGERS FOR AUTOMATION
-- ============================================================================

-- Trigger: Update product UpdatedAt timestamp
CREATE TRIGGER trg_products_update_timestamp
AFTER UPDATE ON Products
FOR EACH ROW
BEGIN
    UPDATE Products SET UpdatedAt = CURRENT_TIMESTAMP WHERE ProductID = NEW.ProductID;
END;

-- Trigger: Update inventory LastUpdated timestamp
CREATE TRIGGER trg_inventory_update_timestamp
AFTER UPDATE ON Inventory
FOR EACH ROW
BEGIN
    UPDATE Inventory SET LastUpdated = CURRENT_TIMESTAMP WHERE InventoryID = NEW.InventoryID;
END;

-- Trigger: Update user UpdatedAt timestamp
CREATE TRIGGER trg_users_update_timestamp
AFTER UPDATE ON Users
FOR EACH ROW
BEGIN
    UPDATE Users SET UpdatedAt = CURRENT_TIMESTAMP WHERE UserID = NEW.UserID;
END;

-- Trigger: Update supplier UpdatedAt timestamp
CREATE TRIGGER trg_suppliers_update_timestamp
AFTER UPDATE ON Suppliers
FOR EACH ROW
BEGIN
    UPDATE Suppliers SET UpdatedAt = CURRENT_TIMESTAMP WHERE SupplierID = NEW.SupplierID;
END;

-- Trigger: Update customer UpdatedAt timestamp
CREATE TRIGGER trg_customers_update_timestamp
AFTER UPDATE ON Customers
FOR EACH ROW
BEGIN
    UPDATE Customers SET UpdatedAt = CURRENT_TIMESTAMP WHERE CustomerID = NEW.CustomerID;
END;

-- ============================================================================
-- INITIAL DATA (Optional)
-- ============================================================================

-- Default admin user (password: admin123 - CHANGE IN PRODUCTION!)
-- Password hash should be generated using bcrypt in production
-- INSERT INTO Users (Username, Name, Email, PasswordHash, Role)
-- VALUES ('admin', 'System Administrator', 'admin@example.com', '$2a$10$...', 'admin');

-- Default location
INSERT INTO Locations (Name, Type, Address)
VALUES ('Main Warehouse', 'warehouse', 'Default storage location');

-- Default categories
INSERT INTO Categories (Name, Description) VALUES
('Fragrances', 'Perfumes and fragrances'),
('Accessories', 'Perfume accessories and gift sets'),
('Samples', 'Sample and tester bottles'),
('Gift Sets', 'Packaged gift sets');

-- ============================================================================
-- VIEWS FOR COMMON QUERIES
-- ============================================================================

-- View: Product Stock Summary
CREATE VIEW vw_product_stock_summary AS
SELECT
    p.ProductID,
    p.Name,
    p.SKU,
    p.Barcode,
    c.Name AS Category,
    p.CostPrice,
    p.SellingPrice,
    COALESCE(SUM(i.Quantity), 0) AS TotalStock,
    COALESCE(SUM(i.ReservedQuantity), 0) AS ReservedStock,
    COALESCE(SUM(i.AvailableQuantity), 0) AS AvailableStock,
    p.MinStockLevel,
    p.ReorderPoint,
    CASE
        WHEN COALESCE(SUM(i.Quantity), 0) <= p.ReorderPoint THEN 'Low Stock'
        WHEN COALESCE(SUM(i.Quantity), 0) <= p.MinStockLevel THEN 'Critical'
        ELSE 'OK'
    END AS StockStatus
FROM Products p
LEFT JOIN Categories c ON p.CategoryID = c.CategoryID
LEFT JOIN Inventory i ON p.ProductID = i.ProductID
WHERE p.IsActive = 1
GROUP BY p.ProductID;

-- View: Inventory by Location
CREATE VIEW vw_inventory_by_location AS
SELECT
    l.Name AS Location,
    l.Type AS LocationType,
    p.Name AS Product,
    p.SKU,
    c.Name AS Category,
    i.Quantity,
    i.ReservedQuantity,
    i.AvailableQuantity,
    i.LastUpdated
FROM Inventory i
INNER JOIN Products p ON i.ProductID = p.ProductID
INNER JOIN Locations l ON i.LocationID = l.LocationID
INNER JOIN Categories c ON p.CategoryID = c.CategoryID
WHERE p.IsActive = 1 AND l.IsActive = 1;

-- ============================================================================
-- SCHEMA VERSION
-- ============================================================================

CREATE TABLE SchemaVersion (
    Version TEXT PRIMARY KEY,
    AppliedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    Description TEXT
);

INSERT INTO SchemaVersion (Version, Description)
VALUES ('2.0.0', 'Extended schema with comprehensive inventory management');
