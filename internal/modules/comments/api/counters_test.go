package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pksep/comments/internal/modules/comments/model"
	"github.com/pksep/comments/internal/modules/comments/repository"
	"github.com/pksep/comments/internal/modules/comments/service"
)

type counterRepoStub struct {
	repository.CommentRepoInterface
	calls    int
	userID   string
	markRead bool
}

func (r *counterRepoStub) Counters(_ context.Context, userID string, ids []string, markRead bool) ([]model.ThreadCounter, error) {
	r.calls++
	r.userID = userID
	r.markRead = markRead
	return []model.ThreadCounter{{ThreadID: ids[0], Total: 2, Unread: 1}}, nil
}

func TestCounterGateway(t *testing.T) {
	gin.SetMode(gin.TestMode)
	key := strings.Repeat("isolated-key", 4)
	threadID := uuid.NewString()
	for _, scenario := range []struct {
		name       string
		key        string
		configured bool
		body       counterRequest
		status     int
	}{
		{"read", key, true, counterRequest{UserID: "7", ThreadIDs: []string{threadID}}, 200},
		{"mark", key, true, counterRequest{UserID: "7", ThreadIDs: []string{threadID}, MarkRead: true}, 200},
		{"missing_key", "", true, counterRequest{UserID: "7", ThreadIDs: []string{threadID}}, 401},
		{"wrong_key", "incorrect", true, counterRequest{UserID: "7", ThreadIDs: []string{threadID}}, 401},
		{"not_configured", key, false, counterRequest{UserID: "7", ThreadIDs: []string{threadID}}, 503},
		{"foreign_identity_format", key, true, counterRequest{UserID: "-1", ThreadIDs: []string{threadID}}, 400},
		{"invalid_thread", key, true, counterRequest{UserID: "7", ThreadIDs: []string{"invalid"}}, 400},
		{"unbounded_request", key, true, counterRequest{UserID: "7", ThreadIDs: make([]string, 501)}, 400},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			configured := ""
			if scenario.configured {
				configured = key
			}
			t.Setenv("COMMENTS_COUNTERS_KEY", configured)
			repo := &counterRepoStub{}
			router := gin.New()
			handler := NewCommentHandler(service.NewCommentService(repo))
			handler.RegisterRoutes(router.Group("/api"))
			body, err := json.Marshal(scenario.body)
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest(http.MethodPost, "/api/comments/counters", strings.NewReader(string(body)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-Comments-Counters-Key", scenario.key)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != scenario.status {
				t.Fatalf("status %d, expected %d: %s", response.Code, scenario.status, response.Body)
			}
			if scenario.status == 200 {
				if repo.calls != 1 || repo.userID != "7" || repo.markRead != scenario.body.MarkRead {
					t.Fatalf("wrong reader or operation: %+v", repo)
				}
			} else if repo.calls != 0 {
				t.Fatal("rejected request reached repository")
			}
		})
	}
}
