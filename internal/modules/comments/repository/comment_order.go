package repository

import "github.com/pksep/comments/internal/modules/comments/model"

// commentComesFirst keeps pinned comments ordered by update, preserving ordinary reply chronology.
func commentComesFirst(left, right model.Comment, oldestFirst bool) bool {
	if left.IsPinned != right.IsPinned {
		return left.IsPinned
	}
	if left.IsPinned {
		return left.UpdatedAt.After(right.UpdatedAt)
	}
	if oldestFirst {
		return left.CreatedAt.Before(right.CreatedAt)
	}
	return left.CreatedAt.After(right.CreatedAt)
}
