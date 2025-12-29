-- Migration: Webinar Feature
-- Created: 2024-12-29

-- Add course_type to courses table
ALTER TABLE courses ADD COLUMN IF NOT EXISTS course_type VARCHAR(20) DEFAULT 'self_paced';
-- Values: 'self_paced', 'webinar', 'hybrid'

-- Create webinars table
CREATE TABLE IF NOT EXISTS webinars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    duration_minutes INTEGER DEFAULT 60,
    meeting_url TEXT,
    meeting_password VARCHAR(100),
    max_participants INTEGER,
    status VARCHAR(20) DEFAULT 'upcoming',
    -- Values: 'draft', 'upcoming', 'live', 'completed', 'cancelled'
    recording_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create webinar_registrations table (who registered for webinar)
CREATE TABLE IF NOT EXISTS webinar_registrations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webinar_id UUID REFERENCES webinars(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    registration_source VARCHAR(50) DEFAULT 'campaign',
    -- Values: 'campaign', 'direct', 'admin'
    attended BOOLEAN DEFAULT FALSE,
    attended_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(webinar_id, user_id)
);

-- Create webinar_reminders table (scheduled reminders)
CREATE TABLE IF NOT EXISTS webinar_reminders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webinar_id UUID REFERENCES webinars(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    reminder_type VARCHAR(20) NOT NULL,
    -- Values: '1_day', '3_hours', '30_min'
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE,
    channel VARCHAR(20) DEFAULT 'whatsapp',
    status VARCHAR(20) DEFAULT 'pending',
    -- Values: 'pending', 'sent', 'failed'
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create wa_notifications table (WhatsApp notification log)
CREATE TABLE IF NOT EXISTS wa_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    phone VARCHAR(20) NOT NULL,
    message_type VARCHAR(50) NOT NULL,
    -- Values: 'webinar_confirmation', 'webinar_reminder_1d', 
    --         'webinar_reminder_3h', 'webinar_reminder_30m', 'payment_success'
    template_data JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    -- Values: 'pending', 'sent', 'failed'
    sent_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_webinars_course ON webinars(course_id);
CREATE INDEX IF NOT EXISTS idx_webinars_scheduled ON webinars(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_webinars_status ON webinars(status);
CREATE INDEX IF NOT EXISTS idx_webinar_reg_user ON webinar_registrations(user_id);
CREATE INDEX IF NOT EXISTS idx_webinar_reg_webinar ON webinar_registrations(webinar_id);
CREATE INDEX IF NOT EXISTS idx_reminders_scheduled ON webinar_reminders(scheduled_at, status);
CREATE INDEX IF NOT EXISTS idx_reminders_pending ON webinar_reminders(status) WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_wa_notifications_status ON wa_notifications(status, created_at);
