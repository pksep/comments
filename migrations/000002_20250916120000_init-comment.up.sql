CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    author_id TEXT NOT NULL,
    content TEXT NOT NULL,
    thread_id UUID NULL REFERENCES threads(id) ON DELETE CASCADE,
    answer_comment_id UUID NULL REFERENCES comments(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
