package repository

import (
	"sort"
	"testing"
	"time"

	"github.com/pksep/comments/internal/modules/comments/model"
)

func TestPinnedCommentOrder(t *testing.T) {
	base := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	for _, oldestFirst := range []bool{false, true} {
		comments := []model.Comment{
			{ID: "ordinary-new", CreatedAt: base.Add(4 * time.Hour)},
			{ID: "pin-old", IsPinned: true, CreatedAt: base.Add(time.Hour), UpdatedAt: base.Add(2 * time.Hour)},
			{ID: "pin-new", IsPinned: true, CreatedAt: base, UpdatedAt: base.Add(3 * time.Hour)},
			{ID: "ordinary-old", CreatedAt: base},
		}
		sort.SliceStable(comments, func(i, j int) bool { return commentComesFirst(comments[i], comments[j], oldestFirst) })
		if comments[0].ID != "pin-new" || comments[1].ID != "pin-old" {
			t.Fatalf("pinned order incorrect: %v", comments)
		}
		expected := "ordinary-new"
		if oldestFirst {
			expected = "ordinary-old"
		}
		if comments[2].ID != expected {
			t.Fatalf("ordinary order changed: %v", comments)
		}
	}
}
