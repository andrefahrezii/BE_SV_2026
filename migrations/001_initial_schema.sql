CREATE SCHEMA IF NOT EXISTS sv_portal;

-- Users
CREATE TABLE IF NOT EXISTS sv_portal.users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user' CHECK (role IN ('admin','user')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Posts
CREATE TABLE IF NOT EXISTS sv_portal.posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    created_date TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_date TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(100) NOT NULL DEFAULT 'draft' CHECK (status IN ('publish','draft','thrash')),
    author_id INT NOT NULL REFERENCES sv_portal.users(id)
);

CREATE INDEX IF NOT EXISTS idx_posts_status ON sv_portal.posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_category ON sv_portal.posts(category);
CREATE INDEX IF NOT EXISTS idx_posts_created_date ON sv_portal.posts(created_date DESC);

-- Full-text search vector (category-only search kept lightweight)
CREATE INDEX IF NOT EXISTS idx_posts_category_fts ON sv_portal.posts USING GIN (to_tsvector('indonesian', coalesce(category,'') || ' ' || coalesce(title,'')));

-- Audit logs
CREATE TABLE IF NOT EXISTS sv_portal.audit_logs (
    id SERIAL PRIMARY KEY,
    actor_user_id INT REFERENCES sv_portal.users(id),
    action VARCHAR(50) NOT NULL CHECK (action IN ('create','update','delete','login','logout')),
    resource_type VARCHAR(100) NOT NULL,
    resource_id INT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    old_values JSONB,
    new_values JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON sv_portal.audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_actor ON sv_portal.audit_logs(actor_user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON sv_portal.audit_logs(resource_type, resource_id);


