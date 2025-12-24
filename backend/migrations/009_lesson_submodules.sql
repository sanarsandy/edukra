-- Migration: Add support for lesson submodules/hierarchy
-- This enables multi-level nesting of lessons within courses

-- Add parent_id column for self-referencing hierarchy
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES lessons(id) ON DELETE CASCADE;

-- Add is_container flag to distinguish folders from content
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS is_container BOOLEAN DEFAULT false;

-- Index for efficient tree queries
CREATE INDEX IF NOT EXISTS idx_lessons_parent_id ON lessons(parent_id);
CREATE INDEX IF NOT EXISTS idx_lessons_course_parent ON lessons(course_id, parent_id);

-- Add composite index for ordering within parent
CREATE INDEX IF NOT EXISTS idx_lessons_parent_order ON lessons(course_id, parent_id, order_index);
