-- Migration: Add is_free_webinar flag to campaigns
-- NULL = use course.price (default behavior)
-- TRUE = force free registration regardless of course price
-- FALSE = force paid even if course is free

ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS is_free_webinar BOOLEAN DEFAULT NULL;
