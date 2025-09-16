CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    author_id TEXT NOT NULL,
    content TEXT NOT NULL,
    parent_id UUID NULL REFERENCES comments(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_comments_entity ON comments (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments (parent_id);

