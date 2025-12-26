-- Migration: Add Full Text Search support for Hybrid Search
-- Combines vector similarity with keyword search for better RAG accuracy

-- Add tsvector column for full text search
ALTER TABLE course_embeddings 
ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Create GIN index for fast full text search
CREATE INDEX IF NOT EXISTS idx_course_embeddings_fts 
ON course_embeddings USING gin(search_vector);

-- Create function to generate tsvector from chunk_text
-- Uses 'indonesian' config if available, fallback to 'simple'
CREATE OR REPLACE FUNCTION generate_search_vector(text_content TEXT)
RETURNS tsvector AS $$
BEGIN
    -- Try Indonesian config first, fallback to simple
    BEGIN
        RETURN to_tsvector('indonesian', COALESCE(text_content, ''));
    EXCEPTION WHEN undefined_object THEN
        RETURN to_tsvector('simple', COALESCE(text_content, ''));
    END;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Update existing rows to populate search_vector
UPDATE course_embeddings 
SET search_vector = generate_search_vector(chunk_text)
WHERE search_vector IS NULL;

-- Create trigger to auto-update search_vector on insert/update
CREATE OR REPLACE FUNCTION update_search_vector()
RETURNS trigger AS $$
BEGIN
    NEW.search_vector := generate_search_vector(NEW.chunk_text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_update_search_vector ON course_embeddings;
CREATE TRIGGER trg_update_search_vector
BEFORE INSERT OR UPDATE OF chunk_text ON course_embeddings
FOR EACH ROW EXECUTE FUNCTION update_search_vector();

-- Add comment
COMMENT ON COLUMN course_embeddings.search_vector IS 'Full Text Search vector for hybrid search (keyword + semantic)';
