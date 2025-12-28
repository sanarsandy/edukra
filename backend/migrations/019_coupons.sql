-- Coupons Feature Migration
-- Created: 2024-12-28

-- Coupons Table
CREATE TABLE IF NOT EXISTS coupons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    
    -- Discount Type: 'percentage' or 'fixed'
    discount_type VARCHAR(20) NOT NULL CHECK (discount_type IN ('percentage', 'fixed')),
    discount_value DECIMAL(12,2) NOT NULL CHECK (discount_value > 0),
    max_discount DECIMAL(12,2), -- Max cap for percentage discounts (null = no cap)
    
    -- Scope (null = all courses)
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    instructor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Usage Limits
    usage_limit INT, -- null = unlimited
    usage_count INT DEFAULT 0,
    per_user_limit INT DEFAULT 1, -- How many times a single user can use this coupon
    
    -- Validity Period
    valid_from TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP WITH TIME ZONE,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Each code must be unique per tenant
    UNIQUE(tenant_id, code)
);

-- Track coupon usage per user for per_user_limit enforcement
CREATE TABLE IF NOT EXISTS coupon_usages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    coupon_id UUID REFERENCES coupons(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    transaction_id UUID REFERENCES transactions(id) ON DELETE SET NULL,
    discount_applied DECIMAL(12,2) NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- One coupon use per transaction
    UNIQUE(coupon_id, transaction_id)
);

-- Add discount tracking to transactions
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS coupon_id UUID REFERENCES coupons(id) ON DELETE SET NULL;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS original_amount DECIMAL(12,2);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS discount_amount DECIMAL(12,2) DEFAULT 0;

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_coupons_code ON coupons(tenant_id, code);
CREATE INDEX IF NOT EXISTS idx_coupons_course ON coupons(course_id) WHERE course_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_coupons_instructor ON coupons(instructor_id) WHERE instructor_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_coupons_active ON coupons(is_active, valid_until);
CREATE INDEX IF NOT EXISTS idx_coupon_usages_user ON coupon_usages(user_id);
CREATE INDEX IF NOT EXISTS idx_coupon_usages_coupon ON coupon_usages(coupon_id);
CREATE INDEX IF NOT EXISTS idx_transactions_coupon ON transactions(coupon_id) WHERE coupon_id IS NOT NULL;
