-- =====================================================
-- Payment Module Migration
-- Add payment-related columns to transactions table
-- =====================================================

-- Add missing columns to transactions table for payment gateway integration
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS snap_token VARCHAR(255);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS payment_url TEXT;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS expired_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS order_id VARCHAR(100) UNIQUE;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS gross_amount DECIMAL(12,2);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS payment_type VARCHAR(50);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS fraud_status VARCHAR(50);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS transaction_time TIMESTAMP WITH TIME ZONE;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS settlement_time TIMESTAMP WITH TIME ZONE;

-- Create index for order_id lookups (used by payment gateway callbacks)
CREATE INDEX IF NOT EXISTS idx_transactions_order_id ON transactions(order_id);

-- Add payment settings to settings table
INSERT INTO settings (key, value)
VALUES 
    ('payment_enabled', 'false'),
    ('payment_provider', 'midtrans'),
    ('payment_midtrans_server_key', ''),
    ('payment_midtrans_client_key', ''),
    ('payment_midtrans_is_production', 'false')
ON CONFLICT (key) DO NOTHING;

