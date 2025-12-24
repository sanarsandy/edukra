-- =====================================================
-- Seed: Create default admin user
-- Password: Admin123!
-- =====================================================

-- Insert admin user if not exists
INSERT INTO users (email, password_hash, full_name, role, auth_provider, is_active)
SELECT 
    'admin@lms.local',
    '$2a$10$TRWAnpwu0uo4pwLtC9.ahuJXasgUb5xHLbPAZArKyek9vaW7Fdzlm', -- Admin123!
    'System Administrator',
    'admin',
    'email',
    true
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'admin@lms.local'
);

-- Verify admin was created
SELECT id, email, full_name, role, is_active, created_at 
FROM users 
WHERE email = 'admin@lms.local';
