package repository

import (
	"context"
	"testing"

	"github.com/wyw14/cry039/internal/domain"
)

func TestIdempotentMigrationPipelineRepository(t *testing.T) {
	repo := NewFeedbackMemory("m", []domain.Feedback{{ID: "before", AreaID: "a"}})
	replacement := []domain.Feedback{{ID: "after", AreaID: "b"}}
	if err := repo.ReplaceMigrationSet(context.Background(), "m", replacement); err != nil {
		t.Fatal(err)
	}
	if err := repo.ReplaceMigrationSet(context.Background(), "m", replacement); err != nil {
		t.Fatal(err)
	}
	stored, err := repo.LoadForMigration(context.Background(), "m")
	if err != nil {
		t.Fatal(err)
	}
	if len(stored) != 1 || stored[0].ID != "after" {
		t.Fatalf("replace appended duplicate rows: %+v", stored)
	}
}
