package application

import (
	"context"
	"errors"
	"github.com/wyw14/cry039/internal/domain"
	"github.com/wyw14/cry039/internal/repository"
	"testing"
	"time"
)

func TestExecuteIsIdempotent(t *testing.T) {
	now := time.Date(2026, 8, 17, 2, 0, 0, 0, time.UTC)
	repo := repository.NewFeedbackMemory("m", []domain.Feedback{{ID: "f", AreaID: "a", SubmittedAt: now}})
	svc := NewMover(repo, func() time.Time { return now })
	job := &domain.Migration{ID: "m", SourceArea: "a", TargetArea: "b"}
	_ = job.Approve("r1")
	_ = job.Approve("r2")
	source := domain.Area{ID: "a", Environment: "open"}
	target := domain.Area{ID: "b", Environment: "open"}
	first, err := svc.Execute(context.Background(), "k", "d", job, source, target, "op")
	if err != nil {
		t.Fatal(err)
	}
	second, err := svc.Execute(context.Background(), "k", "d", job, source, target, "op")
	if err != nil || len(second) != 1 || len(second[0].Timeline) != 1 || first[0].Version != second[0].Version {
		t.Fatalf("first=%+v second=%+v err=%v", first, second, err)
	}
	if _, err = svc.Execute(context.Background(), "k", "other", job, source, target, "op"); !errors.Is(err, ErrMigrationKeyConflict) {
		t.Fatal(err)
	}
}

func TestMigrationHistoryPipelineApplication(t *testing.T) {
	submitted := time.Date(2026, 8, 1, 8, 0, 0, 0, time.UTC)
	prior := domain.Event{At: submitted, Actor: "employee", Action: "submitted"}
	repo := repository.NewFeedbackMemory("m", []domain.Feedback{{ID: "f", AreaID: "a", SubmittedAt: submitted, Timeline: []domain.Event{prior}}})
	mover := NewMover(repo, func() time.Time { return submitted.Add(time.Hour) })
	job := &domain.Migration{ID: "m", SourceArea: "a", TargetArea: "b", Reviewers: []string{"r1", "r2"}}
	moved, err := mover.Execute(context.Background(), "k", "d", job,
		domain.Area{ID: "a", Environment: "open"}, domain.Area{ID: "b", Environment: "open"}, "operator")
	if err != nil {
		t.Fatal(err)
	}
	if !moved[0].SubmittedAt.Equal(submitted) || len(moved[0].Timeline) != 2 || moved[0].Timeline[0] != prior {
		t.Fatalf("application lost feedback history: %+v", moved[0])
	}
}
