-- ============================================================================
-- Inventory Management System - Extended Database Schema (SQLite)
-- ============================================================================
-- This schema supports comprehensive inventory management including:
-- - Product catalog with categories, SKU, batch tracking, and shelf locations
-- - Direct inventory tracking in product table
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
DROP TABLE IF EXISTS Products;
DROP TABLE IF EXISTS Categories;
DROP TABLE IF EXISTS Suppliers;
DROP TABLE IF EXISTS Customers;
DROP TABLE IF EXISTS Settings;
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

-- Settings Table
-- Stores global application settings and business configuration
CREATE TABLE Settings (
    SettingID INTEGER PRIMARY KEY AUTOINCREMENT,
    SettingKey TEXT UNIQUE NOT NULL,
    SettingValue TEXT NOT NULL,
    Description TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
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
-- Core product catalog with batch and expiry tracking and direct inventory management
CREATE TABLE Products (
    ProductID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL,
    Description TEXT,
    SKU TEXT UNIQUE NOT NULL,
    CategoryID INTEGER NOT NULL,
    BatchNo TEXT,
    ExpiryDate DATE,
    CostPrice REAL NOT NULL CHECK(CostPrice >= 0),
    SellingPrice REAL NOT NULL CHECK(SellingPrice >= 0),
    CurrentQuantity INTEGER NOT NULL DEFAULT 0 CHECK(CurrentQuantity >= 0),
    MinStockLevel INTEGER DEFAULT 0,
    MaxStockLevel INTEGER,
    ReorderPoint INTEGER DEFAULT 0,
    Unit TEXT DEFAULT 'pcs' CHECK(Unit IN ('pcs', 'kg', 'liter', 'box', 'carton')),
    ShelfLocation TEXT,
    IsActive INTEGER DEFAULT 1 CHECK(IsActive IN (0, 1)),
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (CategoryID) REFERENCES Categories(CategoryID) ON DELETE RESTRICT
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
    BatchNo TEXT,
    ExpiryDate DATE,
    FOREIGN KEY (PurchaseID) REFERENCES Purchases(PurchaseID) ON DELETE CASCADE,
    FOREIGN KEY (ProductID) REFERENCES Products(ProductID) ON DELETE RESTRICT
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
    PaymentMethod TEXT CHECK(PaymentMethod IN ('cash', 'card', 'mobile', 'bank_transfer', 'credit', 'square_terminal')),
    Notes TEXT,
    SoldBy INTEGER,
    SquarePaymentID TEXT,
    SquareCheckoutID TEXT,
    PaymentProcessedAt DATETIME,
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
    FOREIGN KEY (SaleID) REFERENCES Sales(SaleID) ON DELETE CASCADE,
    FOREIGN KEY (ProductID) REFERENCES Products(ProductID) ON DELETE RESTRICT
);

-- ============================================================================
-- AUDIT & TRANSACTIONS
-- ============================================================================

-- Transactions Table
-- Audit trail for all inventory movements
CREATE TABLE Transactions (
    TransactionID INTEGER PRIMARY KEY AUTOINCREMENT,
    TransactionType TEXT NOT NULL CHECK(TransactionType IN ('purchase', 'sale', 'adjustment', 'return', 'damage', 'theft')),
    TransactionDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ProductID INTEGER NOT NULL,
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
    FOREIGN KEY (PerformedBy) REFERENCES Users(UserID) ON DELETE SET NULL
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

-- User indexes
CREATE INDEX idx_users_username ON Users(Username);
CREATE INDEX idx_users_email ON Users(Email);
CREATE INDEX idx_users_role ON Users(Role);

-- Settings indexes
CREATE INDEX idx_settings_key ON Settings(SettingKey);

-- Product indexes
CREATE INDEX idx_products_sku ON Products(SKU);
CREATE INDEX idx_products_category ON Products(CategoryID);
CREATE INDEX idx_products_name ON Products(Name);
CREATE INDEX idx_products_active ON Products(IsActive);
CREATE INDEX idx_products_quantity ON Products(CurrentQuantity);

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
CREATE INDEX idx_sales_square_payment ON Sales(SquarePaymentID);
CREATE INDEX idx_sales_square_checkout ON Sales(SquareCheckoutID);
CREATE INDEX idx_sale_items_sale ON SalesItems(SaleID);
CREATE INDEX idx_sale_items_product ON SalesItems(ProductID);

-- Transaction indexes
CREATE INDEX idx_transactions_type ON Transactions(TransactionType);
CREATE INDEX idx_transactions_date ON Transactions(TransactionDate);
CREATE INDEX idx_transactions_product ON Transactions(ProductID);

-- Category indexes
CREATE INDEX idx_categories_name ON Categories(Name);
CREATE INDEX idx_categories_parent ON Categories(ParentCategoryID);

-- Supplier indexes
CREATE INDEX idx_suppliers_name ON Suppliers(Name);
CREATE INDEX idx_suppliers_active ON Suppliers(IsActive);

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

-- Trigger: Update settings UpdatedAt timestamp
CREATE TRIGGER trg_settings_update_timestamp
AFTER UPDATE ON Settings
FOR EACH ROW
BEGIN
    UPDATE Settings SET UpdatedAt = CURRENT_TIMESTAMP WHERE SettingID = NEW.SettingID;
END;

-- ============================================================================
-- INITIAL DATA (Optional)
-- ============================================================================

-- Default admin user (password: admin123 - CHANGE IN PRODUCTION!)
-- Password hash should be generated using bcrypt in production

-- Default Users


-- Default categories
INSERT INTO Categories (Name, Description) VALUES
('Fragrances - Men', 'Perfumes and fragrances men'),
('Fragrances - Women', 'Perfumes and fragrances for women'),
('Fragrances - Children', 'Perfumes and fragrances for children'),
('Fragrances - Misc', 'Miscellaneous perfume'),
('Accessories', 'Perfume accessories and gift sets'),
('Samples', 'Sample and tester bottles'),
('Gift Sets', 'Packaged gift sets');

-- Default settings
INSERT INTO Settings (SettingKey, SettingValue, Description) VALUES
('store_name', 'Lux perfumes', 'Name of the store displayed on receipts and documents'),
('store_address', '1950 E20th St # G-710, Chico, CA 95928', 'Store physical address'),
('store_phone', '530 717 3219', 'Store contact phone number'),
('store_policy', 'No refunds. Exchanges within 14 days with original packaging.', 'Store return/exchange policy'),
('sales_tax_rate', '7.25', 'Default sales tax rate (percentage)');

-- ============================================================================
-- VIEWS FOR COMMON QUERIES
-- ============================================================================

-- View: Product Stock Summary
CREATE VIEW vw_product_stock_summary AS
SELECT
    p.ProductID,
    p.Name,
    p.SKU,
    c.Name AS Category,
    p.CostPrice,
    p.SellingPrice,
    p.CurrentQuantity AS TotalStock,
    p.MinStockLevel,
    p.ReorderPoint,
    CASE
        WHEN p.CurrentQuantity <= p.ReorderPoint THEN 'Low Stock'
        WHEN p.CurrentQuantity <= p.MinStockLevel THEN 'Critical'
        ELSE 'OK'
    END AS StockStatus
FROM Products p
LEFT JOIN Categories c ON p.CategoryID = c.CategoryID
WHERE p.IsActive = 1;

-- View: Complete Inventory Details
CREATE VIEW vw_inventory_details AS
SELECT
    p.ProductID,
    p.Name AS Product,
    p.SKU,
    p.ShelfLocation,
    c.Name AS Category,
    p.CostPrice,
    p.SellingPrice,
    p.CurrentQuantity AS Quantity,
    p.MinStockLevel,
    p.ReorderPoint,
    p.UpdatedAt AS LastUpdated
FROM Products p
INNER JOIN Categories c ON p.CategoryID = c.CategoryID
WHERE p.IsActive = 1;

-- ============================================================================
-- SQUARE PAYMENT INTEGRATION
-- ============================================================================

-- SquareOAuthTokens Table
-- Stores OAuth tokens for Square merchant authorization
CREATE TABLE SquareOAuthTokens (
    TokenID INTEGER PRIMARY KEY AUTOINCREMENT,
    MerchantID TEXT UNIQUE NOT NULL,
    AccessToken TEXT NOT NULL,
    RefreshToken TEXT,
    ExpiresAt DATETIME,
    TokenType TEXT DEFAULT 'Bearer',
    Scopes TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for merchant lookup
CREATE INDEX idx_square_tokens_merchant ON SquareOAuthTokens(MerchantID);

-- Trigger: Update SquareOAuthTokens UpdatedAt timestamp
CREATE TRIGGER trg_square_tokens_update_timestamp
AFTER UPDATE ON SquareOAuthTokens
FOR EACH ROW
BEGIN
    UPDATE SquareOAuthTokens SET UpdatedAt = CURRENT_TIMESTAMP WHERE TokenID = NEW.TokenID;
END;

-- ============================================================================
-- SQUARE TERMINAL INTEGRATION
-- ============================================================================

-- SquareTerminalDevices Table
-- Stores registered Terminal devices for merchants
CREATE TABLE SquareTerminalDevices (
    DeviceID INTEGER PRIMARY KEY AUTOINCREMENT,
    MerchantID TEXT NOT NULL,
    SquareDeviceID TEXT UNIQUE NOT NULL,
    DeviceName TEXT NOT NULL,
    DeviceCode TEXT,
    LocationID TEXT,
    Status TEXT DEFAULT 'active' CHECK(Status IN ('active', 'inactive', 'paired', 'unpaired')),
    IsDefault INTEGER DEFAULT 0 CHECK(IsDefault IN (0, 1)),
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (MerchantID) REFERENCES SquareOAuthTokens(MerchantID) ON DELETE CASCADE
);

CREATE INDEX idx_terminal_devices_merchant ON SquareTerminalDevices(MerchantID);
CREATE INDEX idx_terminal_devices_square_id ON SquareTerminalDevices(SquareDeviceID);

-- Trigger: Update SquareTerminalDevices UpdatedAt timestamp
CREATE TRIGGER trg_terminal_devices_update_timestamp
AFTER UPDATE ON SquareTerminalDevices
FOR EACH ROW
BEGIN
    UPDATE SquareTerminalDevices SET UpdatedAt = CURRENT_TIMESTAMP WHERE DeviceID = NEW.DeviceID;
END;

-- SquareTerminalCheckouts Table
-- Tracks Terminal checkout requests and payment status
CREATE TABLE SquareTerminalCheckouts (
    CheckoutID INTEGER PRIMARY KEY AUTOINCREMENT,
    SaleID INTEGER NOT NULL,
    MerchantID TEXT NOT NULL,
    DeviceID INTEGER NOT NULL,
    SquareCheckoutID TEXT UNIQUE NOT NULL,
    AmountMoney INTEGER NOT NULL,
    Currency TEXT DEFAULT 'USD',
    Status TEXT DEFAULT 'PENDING' CHECK(Status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'CANCELED', 'FAILED')),
    SquarePaymentID TEXT,
    ErrorCode TEXT,
    ErrorMessage TEXT,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    CompletedAt DATETIME,
    FOREIGN KEY (SaleID) REFERENCES Sales(SaleID) ON DELETE CASCADE,
    FOREIGN KEY (DeviceID) REFERENCES SquareTerminalDevices(DeviceID) ON DELETE RESTRICT,
    FOREIGN KEY (MerchantID) REFERENCES SquareOAuthTokens(MerchantID) ON DELETE CASCADE
);

CREATE INDEX idx_terminal_checkouts_sale ON SquareTerminalCheckouts(SaleID);
CREATE INDEX idx_terminal_checkouts_square_id ON SquareTerminalCheckouts(SquareCheckoutID);
CREATE INDEX idx_terminal_checkouts_status ON SquareTerminalCheckouts(Status);
CREATE INDEX idx_terminal_checkouts_merchant ON SquareTerminalCheckouts(MerchantID);

-- Trigger: Update SquareTerminalCheckouts UpdatedAt timestamp
CREATE TRIGGER trg_terminal_checkouts_update_timestamp
AFTER UPDATE ON SquareTerminalCheckouts
FOR EACH ROW
BEGIN
    UPDATE SquareTerminalCheckouts SET UpdatedAt = CURRENT_TIMESTAMP WHERE CheckoutID = NEW.CheckoutID;
END;

-- ============================================================================
-- SCHEMA VERSION
-- ============================================================================

CREATE TABLE SchemaVersion (
    Version TEXT PRIMARY KEY,
    AppliedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    Description TEXT
);

INSERT INTO SchemaVersion (Version, Description)
VALUES ('2.4.0', 'Added Square Terminal integration - Terminal devices and checkout tracking');
