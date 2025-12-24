-- Migration: Add Metadata to Users
-- Description: Add JSONB metadata column to users table for storing extra info like instructor bio and specialty

ALTER TABLE users ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}';

-- Create index for metadata (optional, but good for future querying)
CREATE INDEX IF NOT EXISTS idx_users_metadata ON users USING gin (metadata);
