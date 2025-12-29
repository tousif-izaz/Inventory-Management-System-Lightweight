-- Migration: Add Square Terminal Integration
-- Version: 2.4.0
-- Date: 2025-11-12
-- Description: Adds Terminal device management and checkout tracking tables

-- ============================================================================
-- SQUARE TERMINAL INTEGRATION
-- ============================================================================

-- SquareTerminalDevices Table
-- Stores registered Terminal devices for merchants
CREATE TABLE IF NOT EXISTS SquareTerminalDevices (
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

CREATE INDEX IF NOT EXISTS idx_terminal_devices_merchant ON SquareTerminalDevices(MerchantID);
CREATE INDEX IF NOT EXISTS idx_terminal_devices_square_id ON SquareTerminalDevices(SquareDeviceID);

-- Trigger: Update SquareTerminalDevices UpdatedAt timestamp
CREATE TRIGGER IF NOT EXISTS trg_terminal_devices_update_timestamp
AFTER UPDATE ON SquareTerminalDevices
FOR EACH ROW
BEGIN
    UPDATE SquareTerminalDevices SET UpdatedAt = CURRENT_TIMESTAMP WHERE DeviceID = NEW.DeviceID;
END;

-- SquareTerminalCheckouts Table
-- Tracks Terminal checkout requests and payment status
CREATE TABLE IF NOT EXISTS SquareTerminalCheckouts (
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

CREATE INDEX IF NOT EXISTS idx_terminal_checkouts_sale ON SquareTerminalCheckouts(SaleID);
CREATE INDEX IF NOT EXISTS idx_terminal_checkouts_square_id ON SquareTerminalCheckouts(SquareCheckoutID);
CREATE INDEX IF NOT EXISTS idx_terminal_checkouts_status ON SquareTerminalCheckouts(Status);
CREATE INDEX IF NOT EXISTS idx_terminal_checkouts_merchant ON SquareTerminalCheckouts(MerchantID);

-- Trigger: Update SquareTerminalCheckouts UpdatedAt timestamp
CREATE TRIGGER IF NOT EXISTS trg_terminal_checkouts_update_timestamp
AFTER UPDATE ON SquareTerminalCheckouts
FOR EACH ROW
BEGIN
    UPDATE SquareTerminalCheckouts SET UpdatedAt = CURRENT_TIMESTAMP WHERE CheckoutID = NEW.CheckoutID;
END;

-- ============================================================================
-- UPDATE SALES TABLE
-- ============================================================================

-- Add Square payment tracking columns to Sales table
ALTER TABLE Sales ADD COLUMN SquarePaymentID TEXT;
ALTER TABLE Sales ADD COLUMN SquareCheckoutID TEXT;
ALTER TABLE Sales ADD COLUMN PaymentProcessedAt DATETIME;

CREATE INDEX IF NOT EXISTS idx_sales_square_payment ON Sales(SquarePaymentID);
CREATE INDEX IF NOT EXISTS idx_sales_square_checkout ON Sales(SquareCheckoutID);

-- Update schema version
INSERT OR REPLACE INTO SchemaVersion (Version, Description)
VALUES ('2.4.0', 'Added Square Terminal integration - Terminal devices and checkout tracking');
