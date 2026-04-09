CREATE TABLE IF NOT EXISTS comment_media (
    id BIGSERIAL PRIMARY KEY,
    comment_id UUID NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    original_name TEXT NULL,
    type TEXT NOT NULL,
    size BIGINT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_comment_media_comment_id
    ON comment_media(comment_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_comment_media_comment_id_name
    ON comment_media(comment_id, name);
