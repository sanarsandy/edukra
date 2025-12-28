-- Add discount fields to courses table
-- Created: 2024-12-28

-- Discount price (optional - shows strikethrough on original price)
ALTER TABLE courses ADD COLUMN IF NOT EXISTS discount_price DECIMAL(12,2);

-- Discount validity period (optional - auto-expires)
ALTER TABLE courses ADD COLUMN IF NOT EXISTS discount_valid_until TIMESTAMP WITH TIME ZONE;

-- Index for finding active discounts
CREATE INDEX IF NOT EXISTS idx_courses_discount ON courses(discount_price) WHERE discount_price IS NOT NULL;
