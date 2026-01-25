-- Migration: Add encryption keys table for HLS video protection
-- This stores the AES-128 keys used for HLS segment encryption

CREATE TABLE IF NOT EXISTS video_encryption_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    encryption_key BYTEA NOT NULL,  -- 16-byte AES-128 key
    iv BYTEA NOT NULL,              -- 16-byte initialization vector
    hls_path VARCHAR(500),          -- Path to HLS manifest in MinIO
    status VARCHAR(20) DEFAULT 'pending', -- pending, processing, ready, failed
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT unique_lesson_key UNIQUE(lesson_id)
);

-- Index for fast lookups by lesson
CREATE INDEX IF NOT EXISTS idx_video_keys_lesson ON video_encryption_keys(lesson_id);

-- Index for finding pending/processing jobs
CREATE INDEX IF NOT EXISTS idx_video_keys_status ON video_encryption_keys(status);

COMMENT ON TABLE video_encryption_keys IS 'Stores encryption keys for HLS protected video content';
