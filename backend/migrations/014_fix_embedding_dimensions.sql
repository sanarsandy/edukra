-- Migration to fix embedding dimensions for Gemini (768) vs OpenAI (1536)
-- This migration changes the embedding column to support 768 dimensions (Gemini)
-- If you need to support multiple dimensions, consider using a larger dimension or multiple tables

-- Drop the old vector index first
DROP INDEX IF EXISTS idx_course_embeddings_vector;

-- Drop and recreate the embedding column with 768 dimensions
ALTER TABLE course_embeddings DROP COLUMN IF EXISTS embedding;
ALTER TABLE course_embeddings ADD COLUMN embedding vector(768);

-- Recreate the vector index for 768 dimensions
CREATE INDEX IF NOT EXISTS idx_course_embeddings_vector 
ON course_embeddings USING ivfflat (embedding vector_cosine_ops) 
WITH (lists = 100);
