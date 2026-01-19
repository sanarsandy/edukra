-- Add webinar_id column to campaigns for direct webinar linking
-- Allows campaigns to link directly to webinars without going through courses

ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS webinar_id UUID REFERENCES webinars(id);

COMMENT ON COLUMN campaigns.webinar_id IS 'Direct link to webinar for webinar_only or webinar_ecourse campaigns';
