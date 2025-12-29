-- Migration: Add GrapeJS fields to campaigns table
-- This migration adds columns for storing GrapeJS visual editor output

-- Add new columns for GrapeJS visual editor
ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS html_content TEXT;
ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS css_content TEXT;
ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS gjs_data JSONB;

-- Create index on gjs_data for faster queries
CREATE INDEX IF NOT EXISTS idx_campaigns_gjs_data ON campaigns USING gin (gjs_data);
