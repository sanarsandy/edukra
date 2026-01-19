-- Add tracking IDs to campaigns table
ALTER TABLE campaigns ADD COLUMN gtm_id VARCHAR(50);
ALTER TABLE campaigns ADD COLUMN facebook_pixel_id VARCHAR(50);
