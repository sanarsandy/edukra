-- Add Duitku payment settings
INSERT INTO settings (key, value)
VALUES 
    ('payment_duitku_merchant_code', ''),
    ('payment_duitku_merchant_key', ''),
    ('payment_duitku_is_production', 'false')
ON CONFLICT (key) DO NOTHING;
