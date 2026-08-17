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

func TestUndoWindowPipelineApplication(t *testing.T) {
	executed := time.Date(2026, 8, 17, 9, 0, 0, 0, time.FixedZone("CST", 8*3600))
	current := executed
	repo := repository.NewFeedbackMemory("m", []domain.Feedback{{ID: "moved", AreaID: "a"}, {ID: "preexisting", AreaID: "b"}})
	mover := NewMover(repo, func() time.Time { return current })
	job := &domain.Migration{ID: "m", SourceArea: "a", TargetArea: "b", Reviewers: []string{"r1", "r2"}}
	if _, err := mover.Execute(context.Background(), "k", "d", job,
		domain.Area{ID: "a", Environment: "open"}, domain.Area{ID: "b", Environment: "open"}, "operator"); err != nil {
		t.Fatal(err)
	}
	current = executed.In(time.UTC).Add(20 * time.Minute)
	back, err := mover.Undo(context.Background(), job, "operator", 30*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if len(back) != 2 || back[0].AreaID != "a" || back[1].AreaID != "b" {
		t.Fatalf("undo changed unrelated feedback: %+v", back)
	}
}
