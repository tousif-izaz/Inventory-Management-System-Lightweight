-- Migration: Add Square Refunds Support
-- Date: 2025-11-12
-- Description: Add SquareRefunds table to track refund transactions

-- Create SquareRefunds table
CREATE TABLE IF NOT EXISTS SquareRefunds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sale_id INTEGER NOT NULL,
    payment_id TEXT NOT NULL,
    square_refund_id TEXT NOT NULL UNIQUE,
    amount_money INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'PENDING',
    reason TEXT,
    refund_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    completed_at DATETIME,
    FOREIGN KEY (sale_id) REFERENCES Sales(sale_id) ON DELETE CASCADE
);

-- Create index on sale_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_square_refunds_sale_id ON SquareRefunds(sale_id);

-- Create index on square_refund_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_square_refunds_square_id ON SquareRefunds(square_refund_id);

-- Create index on payment_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_square_refunds_payment_id ON SquareRefunds(payment_id);

-- Create index on status for filtering
CREATE INDEX IF NOT EXISTS idx_square_refunds_status ON SquareRefunds(status);
