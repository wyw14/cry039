package service

import (
	"testing"
	"time"
)

func TestUndoWindowPipelineService(t *testing.T) {
	executed := time.Date(2026, 8, 17, 9, 0, 0, 0, time.FixedZone("CST", 8*3600))
	now := executed.In(time.UTC).Add(20 * time.Minute)
	if got := UndoRemaining(executed, now, 30*time.Minute); got != 10*time.Minute {
		t.Fatalf("remaining=%v, want 10m", got)
	}
}
