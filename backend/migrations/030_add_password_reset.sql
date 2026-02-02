-- Migration: Add Password Reset Columns to Users Table
-- Description: Adds columns to store password reset tokens and their expiration time

-- Add reset_password_token column if it doesn't exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_password_token VARCHAR(255);

-- Add reset_password_expires_at column if it doesn't exist
ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_password_expires_at TIMESTAMP WITH TIME ZONE;

-- Create index for faster token lookup
CREATE INDEX IF NOT EXISTS idx_users_reset_password_token ON users(reset_password_token);
