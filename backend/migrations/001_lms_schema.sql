-- LMS Database Schema Migration
-- Based on docs/design_database.md

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Table: Tenants (Whitelabel Control)
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(100) UNIQUE,
    custom_domain VARCHAR(255) UNIQUE,
    
    -- Whitelabel Configuration (Colors, Logo, Font)
    ui_config JSONB DEFAULT '{
        "primary_color": "#4F46E5",
        "secondary_color": "#10B981",
        "logo_url": null,
        "favicon_url": null,
        "font_family": "Inter"
    }',
    
    -- Feature Flags
    feature_config JSONB DEFAULT '{
        "enable_quiz": true,
        "enable_certificate": true,
        "enable_forum": false,
        "drm_protection_level": "basic"
    }',
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Table: Users (Students, Instructors, Tenant Admins)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash TEXT,
    role VARCHAR(50) DEFAULT 'student', -- admin, instructor, student
    full_name VARCHAR(255),
    avatar_url TEXT,
    google_id VARCHAR(255),
    auth_provider VARCHAR(50) DEFAULT 'email', -- email, google, both
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, email) -- Email unique per tenant
);

-- Create indexes for users
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;

-- 3. Table: Courses
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    instructor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255),
    description TEXT,
    thumbnail_url TEXT,
    price DECIMAL(12,2) DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'IDR',
    is_published BOOLEAN DEFAULT false,
    is_featured BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for courses
CREATE INDEX IF NOT EXISTS idx_courses_tenant_id ON courses(tenant_id);
CREATE INDEX IF NOT EXISTS idx_courses_instructor_id ON courses(instructor_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_courses_tenant_slug ON courses(tenant_id, slug);

-- 4. Table: Lessons (Modules & Materials)
CREATE TABLE IF NOT EXISTS lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INT DEFAULT 0,
    content_type VARCHAR(50), -- 'video', 'pdf', 'quiz', 'text'
    
    -- Video & Protection Metadata
    video_url TEXT, -- Path to .m3u8 if DRM
    video_duration INT DEFAULT 0, -- Duration in seconds
    security_level VARCHAR(50) DEFAULT 'standard', -- 'public', 'signed_url', 'aes_128', 'full_drm'
    is_preview BOOLEAN DEFAULT false, -- Can be viewed without purchase (marketing)
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for lessons
CREATE INDEX IF NOT EXISTS idx_lessons_course_id ON lessons(course_id);
CREATE INDEX IF NOT EXISTS idx_lessons_order ON lessons(course_id, order_index);

-- 5. Table: Content Keys (For Content Security)
CREATE TABLE IF NOT EXISTS content_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id UUID REFERENCES lessons(id) ON DELETE CASCADE,
    encryption_key VARCHAR(255) NOT NULL, -- AES-128 Key
    iv VARCHAR(255), -- Initialization Vector
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for content_keys
CREATE INDEX IF NOT EXISTS idx_content_keys_lesson_id ON content_keys(lesson_id);

-- 6. Table: Transactions (Payment Agnostic)
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(id) ON DELETE SET NULL,
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'IDR',
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'success', 'failed', 'refunded'
    payment_gateway VARCHAR(50), -- 'midtrans', 'stripe', 'xendit'
    payment_gateway_ref VARCHAR(255), -- Reference ID from gateway
    payment_method VARCHAR(100), -- 'bank_transfer', 'credit_card', 'e_wallet'
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for transactions
CREATE INDEX IF NOT EXISTS idx_transactions_tenant_id ON transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_course_id ON transactions(course_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);

-- 7. Table: Enrollments (User-Course relationship)
CREATE TABLE IF NOT EXISTS enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    transaction_id UUID REFERENCES transactions(id) ON DELETE SET NULL,
    progress_percentage INT DEFAULT 0,
    completed_at TIMESTAMP WITH TIME ZONE,
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, course_id)
);

-- Create indexes for enrollments
CREATE INDEX IF NOT EXISTS idx_enrollments_user_id ON enrollments(user_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_course_id ON enrollments(course_id);

-- 8. Table: Lesson Progress (Track user progress per lesson)
CREATE TABLE IF NOT EXISTS lesson_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    lesson_id UUID REFERENCES lessons(id) ON DELETE CASCADE,
    is_completed BOOLEAN DEFAULT false,
    watch_time INT DEFAULT 0, -- In seconds
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, lesson_id)
);

-- Create indexes for lesson_progress
CREATE INDEX IF NOT EXISTS idx_lesson_progress_user_id ON lesson_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_lesson_progress_lesson_id ON lesson_progress(lesson_id);
