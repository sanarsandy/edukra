-- Migration: Add order_id column to transactions table
-- This fixes the "Failed to create transaction" error on production

-- Add order_id column if not exists
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS order_id VARCHAR(100);

-- Create index for faster lookups by order_id
CREATE INDEX IF NOT EXISTS idx_transactions_order_id ON transactions(order_id);

-- Add comment for documentation
COMMENT ON COLUMN transactions.order_id IS 'External order ID sent to payment gateway (e.g., LMS-abc123-12345)';
