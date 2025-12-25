-- Settings table for platform configuration
-- This table stores key-value pairs for platform settings

CREATE TABLE IF NOT EXISTS settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_settings_key ON settings(key);

-- Insert default values
INSERT INTO settings (key, value) VALUES 
    ('site_name', 'LearnHub'),
    ('site_description', 'Platform pembelajaran online'),
    ('contact_email', 'admin@learnhub.id'),
    ('currency', 'IDR'),
    ('language', 'id'),
    ('logo_url', ''),
    ('theme', 'default')
ON CONFLICT (key) DO NOTHING;
