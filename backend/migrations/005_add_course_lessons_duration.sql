-- Migration: Add lessons_count and duration to courses table
-- This stores the estimated number of lessons and total duration

ALTER TABLE courses ADD COLUMN IF NOT EXISTS lessons_count INTEGER DEFAULT 0;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS duration VARCHAR(50) DEFAULT '';
ALTER TABLE courses ADD COLUMN IF NOT EXISTS category_id UUID REFERENCES categories(id) ON DELETE SET NULL;

-- Create index for category
CREATE INDEX IF NOT EXISTS idx_courses_category_id ON courses(category_id);

