package repository

import (
	"context"
	"testing"
	"time"

	"github.com/wyw14/cry039/internal/domain"
)

func TestMigrationHistoryPipelineRepository(t *testing.T) {
	now := time.Now()
	seed := []domain.Feedback{{ID: "f", Timeline: []domain.Event{
		{At: now, Action: "submitted"},
		{At: now.Add(time.Minute), Action: "triaged"},
	}}}
	repo := NewFeedbackMemory("m", seed)
	items, err := repo.LoadForMigration(context.Background(), "m")
	if err != nil {
		t.Fatal(err)
	}
	if len(items[0].Timeline) != 2 {
		t.Fatalf("repository truncated history: %+v", items[0].Timeline)
	}
	items[0].Timeline[0].Action = "tampered"
	again, err := repo.LoadForMigration(context.Background(), "m")
	if err != nil {
		t.Fatal(err)
	}
	if again[0].Timeline[0].Action != "submitted" {
		t.Fatalf("read result mutated stored audit: %+v", again[0].Timeline)
	}
}
