-- Campaign Landing Pages Migration
-- Created: 2024-12-29
-- Approach: Template-based system with JSONB blocks

-- Main campaigns table
CREATE TABLE IF NOT EXISTS campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- URL & Status
    slug VARCHAR(100) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT false,
    
    -- Link to Course (optional - for auto-pricing and checkout)
    course_id UUID REFERENCES courses(id) ON DELETE SET NULL,
    
    -- SEO & Meta
    title VARCHAR(255) NOT NULL,
    meta_description TEXT,
    og_image_url TEXT,
    
    -- Content Blocks (JSON - ordered array of blocks)
    -- Each block: { id, type, enabled, order, data }
    blocks JSONB DEFAULT '[]',
    
    -- Global Styles (optional future use)
    styles JSONB DEFAULT '{}',
    
    -- Campaign Period (for countdown timer & auto-deactivation)
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    
    -- Simple Analytics Counters
    view_count INT DEFAULT 0,
    click_count INT DEFAULT 0,
    conversion_count INT DEFAULT 0,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Detailed analytics table for tracking
CREATE TABLE IF NOT EXISTS campaign_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    
    -- Event type: 'view', 'click', 'conversion'
    event_type VARCHAR(20) NOT NULL,
    
    -- For conversions, link to transaction
    transaction_id UUID REFERENCES transactions(id) ON DELETE SET NULL,
    
    -- Visitor info
    visitor_ip VARCHAR(45),
    user_agent TEXT,
    referer TEXT,
    
    -- UTM tracking
    utm_source VARCHAR(100),
    utm_medium VARCHAR(100),
    utm_campaign VARCHAR(100),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_campaigns_slug ON campaigns(slug);
CREATE INDEX IF NOT EXISTS idx_campaigns_active ON campaigns(is_active, end_date);
CREATE INDEX IF NOT EXISTS idx_campaigns_course ON campaigns(course_id);
CREATE INDEX IF NOT EXISTS idx_campaign_analytics_campaign ON campaign_analytics(campaign_id);
CREATE INDEX IF NOT EXISTS idx_campaign_analytics_event ON campaign_analytics(campaign_id, event_type);
CREATE INDEX IF NOT EXISTS idx_campaign_analytics_created ON campaign_analytics(created_at);
