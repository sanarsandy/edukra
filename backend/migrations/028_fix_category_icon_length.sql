-- Fix category icon column to support longer SVG paths
-- Previous VARCHAR(100) was too small for SVG path data (some paths are 150+ chars)

ALTER TABLE categories ALTER COLUMN icon TYPE TEXT;
