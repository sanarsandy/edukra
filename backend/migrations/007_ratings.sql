-- Ratings Migration
-- Course ratings and reviews system

CREATE TABLE IF NOT EXISTS course_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, course_id) -- 1 rating per user per course
);

CREATE INDEX IF NOT EXISTS idx_ratings_course_id ON course_ratings(course_id);
CREATE INDEX IF NOT EXISTS idx_ratings_user_id ON course_ratings(user_id);
