package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wyw14/cry039/internal/application"
	"github.com/wyw14/cry039/internal/domain"
	"github.com/wyw14/cry039/internal/repository"
	"go.uber.org/zap"
)

func TestReviewerQuorumPipelineHTTP(t *testing.T) {
	repo := repository.NewFeedbackMemory("m", []domain.Feedback{{ID: "f", AreaID: "a"}})
	app := application.NewMover(repo, time.Now)
	router := Server(app, zap.NewNop())
	body := `{"idempotency_key":"k","digest":"d","migration_id":"m","source":{"ID":"a","Environment":"open"},"target":{"ID":"b","Environment":"open"},"actor":"operator","reviewers":["Alice"," alice "]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/migrations/execute", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != http.StatusConflict {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
}
