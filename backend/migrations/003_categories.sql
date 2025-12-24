-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    course_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for tenant lookup
CREATE INDEX IF NOT EXISTS idx_categories_tenant ON categories(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_slug_tenant ON categories(tenant_id, slug);

-- Insert default categories
INSERT INTO categories (name, slug, description, icon) VALUES
    ('Web Development', 'web-development', 'Learn to build modern web applications', '💻'),
    ('Design', 'design', 'Master UI/UX and graphic design', '🎨'),
    ('Data Science', 'data-science', 'Explore data analysis and machine learning', '📊'),
    ('Mobile Development', 'mobile-development', 'Build iOS and Android applications', '📱'),
    ('Business', 'business', 'Business skills and entrepreneurship', '💼'),
    ('Marketing', 'marketing', 'Digital marketing and growth strategies', '📢')
ON CONFLICT DO NOTHING;
