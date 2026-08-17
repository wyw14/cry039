package repository

import (
	"testing"

	"github.com/wyw14/cry039/internal/domain"
)

func TestStatsRecalculationPipelineRepository(t *testing.T) {
	repo := NewFeedbackMemory("migration-1", []domain.Feedback{{ID: "f", AreaID: "b", Status: "open", Level: 4}})
	stats := repo.Stats("migration-1", "b")
	if stats.AreaID != "b" || stats.Total != 1 || stats.Open != 1 || stats.Severe != 1 {
		t.Fatalf("repository statistics used the wrong scope: %+v", stats)
	}
}
