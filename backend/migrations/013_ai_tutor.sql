-- =============================================
-- AI TUTOR FEATURE - DATABASE MIGRATION
-- =============================================
-- Migration: 013_ai_tutor.sql
-- Description: Adds pgvector extension, embeddings storage, 
--              chat sessions, and AI usage tracking

-- 1. Enable pgvector extension for vector similarity search
CREATE EXTENSION IF NOT EXISTS vector;

-- 2. Course Content Embeddings
-- Stores chunked content with vector embeddings for RAG
CREATE TABLE IF NOT EXISTS course_embeddings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    lesson_id UUID REFERENCES lessons(id) ON DELETE CASCADE,
    
    -- Content source metadata
    content_type VARCHAR(30) NOT NULL, -- 'text', 'pdf', 'video_transcript', 'quiz'
    source_reference TEXT,             -- e.g., page number, timestamp, section title
    
    -- Chunked content
    chunk_index INT NOT NULL DEFAULT 0,
    chunk_text TEXT NOT NULL,
    chunk_tokens INT DEFAULT 0,
    
    -- Vector embedding (1536 dimensions for OpenAI ada-002)
    -- Can be adjusted for other embedding models
    embedding vector(1536),
    
    -- Additional metadata as JSON
    metadata JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_course_embeddings_course ON course_embeddings(course_id);
CREATE INDEX IF NOT EXISTS idx_course_embeddings_lesson ON course_embeddings(lesson_id);
CREATE INDEX IF NOT EXISTS idx_course_embeddings_type ON course_embeddings(content_type);

-- Vector similarity index using IVFFlat for approximate nearest neighbor search
-- Lists parameter should be sqrt(n) where n is expected number of rows
CREATE INDEX IF NOT EXISTS idx_course_embeddings_vector 
ON course_embeddings USING ivfflat (embedding vector_cosine_ops) 
WITH (lists = 100);

-- 3. Content Processing Status
-- Tracks the status of content processing for each lesson
CREATE TABLE IF NOT EXISTS content_processing_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    lesson_id UUID REFERENCES lessons(id) ON DELETE CASCADE,
    
    -- Processing status
    status VARCHAR(30) NOT NULL DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed'
    content_type VARCHAR(30) NOT NULL,
    
    -- Progress tracking
    total_chunks INT DEFAULT 0,
    processed_chunks INT DEFAULT 0,
    error_message TEXT,
    
    -- Timing
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_processing_status_course ON content_processing_status(course_id);
CREATE INDEX IF NOT EXISTS idx_processing_status_status ON content_processing_status(status);

-- 4. AI Chat Sessions
-- One active session per user per course
CREATE TABLE IF NOT EXISTS ai_chat_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    
    -- Session metadata
    title VARCHAR(255) DEFAULT 'Chat Baru',
    is_active BOOLEAN DEFAULT true,
    
    -- Analytics
    message_count INT DEFAULT 0,
    total_tokens_used INT DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ai_chat_sessions_user ON ai_chat_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_chat_sessions_course ON ai_chat_sessions(course_id);
-- Ensure only one active session per user per course
CREATE UNIQUE INDEX IF NOT EXISTS idx_ai_chat_sessions_active 
ON ai_chat_sessions(user_id, course_id) WHERE is_active = true;

-- 5. AI Chat Messages
-- Stores all messages in a chat session
CREATE TABLE IF NOT EXISTS ai_chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES ai_chat_sessions(id) ON DELETE CASCADE,
    
    -- Message content
    role VARCHAR(20) NOT NULL, -- 'user', 'assistant', 'system'
    content TEXT NOT NULL,
    
    -- For assistant messages: source references
    -- Array of {lesson_id, title, chunk_preview, relevance_score}
    sources JSONB DEFAULT '[]',
    
    -- Token usage for this message
    tokens_used INT DEFAULT 0,
    
    -- Provider info (for assistant messages)
    provider VARCHAR(30),  -- 'openai', 'claude', 'groq', 'gemini'
    model VARCHAR(100),    -- 'gpt-4-turbo', 'claude-3-sonnet', etc.
    
    -- Timestamp
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ai_chat_messages_session ON ai_chat_messages(session_id);
CREATE INDEX IF NOT EXISTS idx_ai_chat_messages_created ON ai_chat_messages(created_at);

-- 6. AI Usage Log
-- For rate limiting, analytics, and optional billing
CREATE TABLE IF NOT EXISTS ai_usage_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(id) ON DELETE SET NULL,
    
    -- Action details
    action_type VARCHAR(30) NOT NULL, -- 'chat', 'embedding', 'transcription'
    provider VARCHAR(30),
    model VARCHAR(100),
    
    -- Token usage
    tokens_input INT DEFAULT 0,
    tokens_output INT DEFAULT 0,
    
    -- Cost tracking (optional, in USD)
    estimated_cost DECIMAL(10, 6) DEFAULT 0,
    
    -- Timestamp
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ai_usage_log_user ON ai_usage_log(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_usage_log_date ON ai_usage_log(created_at);
-- For rate limiting queries (count per user per day)
CREATE INDEX IF NOT EXISTS idx_ai_usage_log_user_date 
ON ai_usage_log(user_id, DATE(created_at));

-- 7. Insert default AI settings
-- These will be stored in the existing settings table
INSERT INTO settings (key, value) VALUES 
    ('ai_enabled', 'false'),
    ('ai_provider', 'openai'),
    ('ai_model', 'gpt-4-turbo'),
    ('ai_api_key_openai', ''),
    ('ai_api_key_claude', ''),
    ('ai_api_key_groq', ''),
    ('ai_api_key_gemini', ''),
    ('ai_embedding_provider', 'openai'),
    ('ai_embedding_model', 'text-embedding-ada-002'),
    ('ai_max_tokens', '2048'),
    ('ai_temperature', '0.7'),
    ('ai_rate_limit_per_day', '50'),
    ('ai_system_prompt', 'Kamu adalah AI Tutor yang membantu siswa memahami materi kursus. Jawab pertanyaan berdasarkan materi yang tersedia. Jika tidak tahu jawabannya, katakan dengan jujur. Gunakan bahasa Indonesia yang baik dan benar.')
ON CONFLICT (key) DO NOTHING;
