-- Migration: Add Square OAuth Integration
-- Version: 2.3.0
-- Date: 2025-11-11
-- Description: Adds SquareOAuthTokens table for storing Square merchant OAuth tokens

-- ============================================================================
-- SQUARE PAYMENT INTEGRATION
-- ============================================================================

-- SquareOAuthTokens Table
-- Stores OAuth tokens for Square merchant authorization
CREATE TABLE IF NOT EXISTS SquareOAuthTokens (
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
CREATE INDEX IF NOT EXISTS idx_square_tokens_merchant ON SquareOAuthTokens(MerchantID);

-- Trigger: Update SquareOAuthTokens UpdatedAt timestamp
CREATE TRIGGER IF NOT EXISTS trg_square_tokens_update_timestamp
AFTER UPDATE ON SquareOAuthTokens
FOR EACH ROW
BEGIN
    UPDATE SquareOAuthTokens SET UpdatedAt = CURRENT_TIMESTAMP WHERE TokenID = NEW.TokenID;
END;

-- Update schema version (if SchemaVersion table exists)
-- Uses INSERT OR REPLACE to handle both new and existing version entries
INSERT OR REPLACE INTO SchemaVersion (Version, Description)
VALUES ('2.3.0', 'Added Square OAuth integration - SquareOAuthTokens table');
