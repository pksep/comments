CREATE TABLE IF NOT EXISTS comment_reads (
    user_id TEXT NOT NULL,
    comment_id UUID NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    CONSTRAINT comment_reads_pkey PRIMARY KEY (user_id, comment_id)
);

CREATE INDEX IF NOT EXISTS comments_visible_thread_idx
    ON comments (thread_id) WHERE status <> 'deleted';
