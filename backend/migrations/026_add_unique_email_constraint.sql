-- Add unique constraint on email to prevent duplicate users
-- First, we need to handle any existing duplicates

-- Find duplicate emails and keep only the first created user
-- For each duplicate, we identify the user to keep (earliest created_at, or if same, smallest id)
WITH duplicate_groups AS (
    SELECT 
        email,
        id,
        created_at,
        ROW_NUMBER() OVER (PARTITION BY LOWER(email) ORDER BY created_at ASC, id ASC) as rn
    FROM users
    WHERE email IS NOT NULL AND email != ''
),
users_to_keep AS (
    SELECT id FROM duplicate_groups WHERE rn = 1
),
users_to_delete AS (
    SELECT dg.id as delete_id, 
           (SELECT id FROM duplicate_groups WHERE LOWER(email) = LOWER(dg.email) AND rn = 1) as keep_id
    FROM duplicate_groups dg 
    WHERE rn > 1
)
-- First update enrollments
UPDATE enrollments e
SET user_id = utd.keep_id
FROM users_to_delete utd
WHERE e.user_id = utd.delete_id
  AND NOT EXISTS (
    SELECT 1 FROM enrollments e2 WHERE e2.user_id = utd.keep_id AND e2.course_id = e.course_id
  );

-- Also update transactions to point to the kept user
WITH duplicate_groups AS (
    SELECT 
        email,
        id,
        created_at,
        ROW_NUMBER() OVER (PARTITION BY LOWER(email) ORDER BY created_at ASC, id ASC) as rn
    FROM users
    WHERE email IS NOT NULL AND email != ''
),
users_to_delete AS (
    SELECT dg.id as delete_id, 
           (SELECT id FROM duplicate_groups d2 WHERE LOWER(d2.email) = LOWER(dg.email) AND d2.rn = 1) as keep_id
    FROM duplicate_groups dg 
    WHERE rn > 1
)
UPDATE transactions t
SET user_id = utd.keep_id
FROM users_to_delete utd
WHERE t.user_id = utd.delete_id;

-- Delete duplicate users (keeping the first one by created_at)
WITH duplicate_groups AS (
    SELECT 
        email,
        id,
        ROW_NUMBER() OVER (PARTITION BY LOWER(email) ORDER BY created_at ASC, id ASC) as rn
    FROM users
    WHERE email IS NOT NULL AND email != ''
)
DELETE FROM users
WHERE id IN (SELECT id FROM duplicate_groups WHERE rn > 1);

-- Now add unique constraint on email (case-insensitive using lower)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON users (LOWER(email)) WHERE email IS NOT NULL AND email != '';
