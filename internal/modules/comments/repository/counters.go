package repository

import (
	"context"

	"github.com/pksep/comments/internal/modules/comments/model"
)

// Counters reads metadata only. Receipts mark the committed snapshot, not a future timestamp.
func (r *CommentRepo) Counters(ctx context.Context, userID string, threadIDs []string, markRead bool) ([]model.ThreadCounter, error) {
	result := make([]model.ThreadCounter, 0)
	if len(threadIDs) == 0 {
		return result, nil
	}

	if markRead {
		_, err := r.db.Exec(ctx, `
			INSERT INTO comment_reads (user_id, comment_id)
			SELECT $1, id FROM comments
			WHERE thread_id = ANY($2::uuid[]) AND status <> 'deleted' AND author_id <> $1
			ON CONFLICT DO NOTHING
		`, userID, threadIDs)
		if err != nil {
			return nil, err
		}
	}

	// A comment committed after the receipt statement remains unread, even with an older created_at.
	rows, err := r.db.Query(ctx, `
		SELECT c.thread_id::text, COUNT(*),
		       COUNT(*) FILTER (WHERE c.author_id <> $1 AND seen.comment_id IS NULL)
		FROM comments c
		LEFT JOIN comment_reads seen ON seen.comment_id = c.id AND seen.user_id = $1
		WHERE c.thread_id = ANY($2::uuid[]) AND c.status <> 'deleted'
		GROUP BY c.thread_id ORDER BY c.thread_id
	`, userID, threadIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.ThreadCounter
		if err := rows.Scan(&item.ThreadID, &item.Total, &item.Unread); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}
