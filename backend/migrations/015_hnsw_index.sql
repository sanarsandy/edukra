-- Migration: Switch from IVFFlat to HNSW index for better vector search performance
-- HNSW (Hierarchical Navigable Small World) is faster and more accurate than IVFFlat,
-- especially for larger datasets. It doesn't require training and has O(log n) complexity.

-- Drop existing IVFFlat index
DROP INDEX IF EXISTS idx_course_embeddings_ivfflat;
DROP INDEX IF EXISTS idx_course_embeddings_vector;

-- Create HNSW index with cosine distance operator
-- Parameters:
--   m = 16: number of connections per layer (default, good balance)
--   ef_construction = 64: size of dynamic candidate list during build (default)
CREATE INDEX idx_course_embeddings_hnsw 
ON course_embeddings 
USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- Add comment for documentation
COMMENT ON INDEX idx_course_embeddings_hnsw IS 'HNSW index for fast vector similarity search with cosine distance';
