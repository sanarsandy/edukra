-- Add campaign_type column for flexible campaign behavior
-- Values: 'webinar_only', 'ecourse_only', 'webinar_ecourse'
-- Default: 'ecourse_only' for backward compatibility

ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS campaign_type VARCHAR(20) DEFAULT 'ecourse_only';

-- Add comment for documentation
COMMENT ON COLUMN campaigns.campaign_type IS 'webinar_only=webinar registration only, ecourse_only=course access only, webinar_ecourse=both';
