-- Add category_id to courses table
ALTER TABLE courses ADD COLUMN IF NOT EXISTS category_id UUID REFERENCES categories(id) ON DELETE SET NULL;

-- Create index for category_id
CREATE INDEX IF NOT EXISTS idx_courses_category_id ON courses(category_id);
