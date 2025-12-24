-- =====================================================
-- Migration: Instructor Workflow
-- Adds course status for review workflow
-- =====================================================

-- Add status column for course review workflow
-- Status: draft, pending_review, approved, rejected, published
ALTER TABLE courses ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'draft';

-- Add review-related columns
ALTER TABLE courses ADD COLUMN IF NOT EXISTS review_notes TEXT;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS submitted_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL;

-- Create index for status filtering
CREATE INDEX IF NOT EXISTS idx_courses_status ON courses(status);
CREATE INDEX IF NOT EXISTS idx_courses_instructor_status ON courses(instructor_id, status);

-- Update existing courses to have proper status based on is_published
UPDATE courses SET status = CASE 
    WHEN is_published = true THEN 'published'
    ELSE 'draft'
END WHERE status IS NULL OR status = '';

-- Notifications table for instructor/admin communication
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    type VARCHAR(50) NOT NULL, -- course_submitted, course_approved, course_rejected, etc
    title VARCHAR(255) NOT NULL,
    message TEXT,
    reference_id UUID,
    reference_type VARCHAR(50), -- course, lesson, etc
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications(created_at DESC);


