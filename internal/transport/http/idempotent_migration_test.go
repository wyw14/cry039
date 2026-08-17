package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/wyw14/cry039/internal/application"
	"github.com/wyw14/cry039/internal/domain"
	"github.com/wyw14/cry039/internal/repository"
	"go.uber.org/zap"
)

func TestIdempotentMigrationPipelineHTTP(t *testing.T) {
	repo := repository.NewFeedbackMemory("m", []domain.Feedback{{ID: "f", AreaID: "a"}})
	router := Server(application.NewMover(repo, time.Now), zap.NewNop())
	body := `{"idempotency_key":"k","digest":"d","migration_id":"m","source":{"ID":"a","Environment":"open"},"target":{"ID":"b","Environment":"open"},"actor":"operator","reviewers":["r1","r2"]}`
	for attempt := 0; attempt < 2; attempt++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/migrations/execute", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"migrated":1`) {
			t.Fatalf("attempt=%d status=%d body=%s", attempt, res.Code, res.Body.String())
		}
	}
}
