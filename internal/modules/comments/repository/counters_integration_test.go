package repository

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCommentCountersWithIsolatedPostgres(t *testing.T) {
	databaseURL := os.Getenv("COMMENTS_COUNTERS_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("isolated PostgreSQL is required")
	}
	address, err := url.Parse(databaseURL)
	if err != nil || !strings.HasPrefix(address.Path, "/comments_counter_test_") {
		t.Fatal("refusing a non-isolated database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	cfg.MaxConns = 2
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	files, err := filepath.Glob(filepath.Join("..", "..", "..", "..", "migrations", "*.up.sql"))
	if err != nil || len(files) != 6 {
		t.Fatalf("migration files unavailable: %d", len(files))
	}
	for _, file := range files {
		migration, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := pool.Exec(ctx, string(migration)); err != nil {
			t.Fatal(err)
		}
	}
	repo := &CommentRepo{db: pool}
	threadID := uuid.NewString()
	if _, err := pool.Exec(ctx, "INSERT INTO threads (id) VALUES ($1)", threadID); err != nil {
		t.Fatal(err)
	}
	ids := []string{uuid.NewString(), uuid.NewString(), uuid.NewString()}
	for index, author := range []string{"1", "2", "2"} {
		if _, err := pool.Exec(ctx, `INSERT INTO comments (id, author_id, content, thread_id, created_at, updated_at) VALUES ($1,$2,'test',$3,NOW(),NOW())`, ids[index], author, threadID); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := pool.Exec(ctx, "UPDATE comments SET answer_comment_id=$1 WHERE id=$2", ids[1], ids[2]); err != nil {
		t.Fatal(err)
	}
	initial, err := repo.Counters(ctx, "1", []string{threadID, threadID}, false)
	if err != nil || len(initial) != 1 || initial[0].Total != 3 || initial[0].Unread != 2 {
		t.Fatalf("root, reply or own-comment count wrong: %+v, %v", initial, err)
	}
	for count := 0; count < 2; count++ {
		seen, err := repo.Counters(ctx, "1", []string{threadID}, true)
		if err != nil || seen[0].Unread != 0 || seen[0].Total != 3 {
			t.Fatalf("mark failed: %+v, %v", seen, err)
		}
	}
	other, err := repo.Counters(ctx, "2", []string{threadID}, false)
	if err != nil || other[0].Unread != 1 {
		t.Fatal("one reader cleared another reader")
	}
	var receipts int
	if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM comment_reads WHERE user_id='1'").Scan(&receipts); err != nil || receipts != 2 {
		t.Fatal("idempotent marking created duplicate/own receipts")
	}

	// The new comment gets an old creation time and commits only after marking the current snapshot.
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	late := uuid.NewString()
	if _, err := tx.Exec(ctx, `INSERT INTO comments (id,author_id,content,thread_id,created_at,updated_at) VALUES ($1,'2','late',$2,'2020-01-01',NOW())`, late, threadID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Counters(ctx, "1", []string{threadID}, true); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	after, err := repo.Counters(ctx, "1", []string{threadID}, false)
	if err != nil || after[0].Total != 4 || after[0].Unread != 1 {
		t.Fatal("a concurrent late commit was incorrectly marked as read")
	}
	if _, err := pool.Exec(ctx, "UPDATE comments SET status='deleted',deleted_at=NOW() WHERE id=$1", late); err != nil {
		t.Fatal(err)
	}
	deleted, err := repo.Counters(ctx, "1", []string{threadID}, false)
	if err != nil || deleted[0].Total != 3 || deleted[0].Unread != 0 {
		t.Fatal("deleted comments remain in counters")
	}
	up, err := os.ReadFile(files[5])
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, string(up)); err != nil {
		t.Fatal("migration is not idempotent", err)
	}
}
