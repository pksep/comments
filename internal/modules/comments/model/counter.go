package model

// ThreadCounter includes both root comments and replies, but never deleted comments.
type ThreadCounter struct {
	ThreadID string `json:"thread_id"`
	Total    int64  `json:"total"`
	Unread   int64  `json:"unread"`
}
