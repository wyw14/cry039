package domain

import (
	"errors"
	"testing"
	"time"
)

func approved() Migration {
	m := Migration{ID: "m", SourceArea: "a", TargetArea: "b"}
	_ = m.Approve("reviewer-1")
	_ = m.Approve("reviewer-2")
	return m
}
func TestApprovalRequiresDistinctPeople(t *testing.T) {
	m := Migration{}
	if err := m.Approve("same"); err != nil {
		t.Fatal(err)
	}
	if !errors.Is(m.Approve("same"), ErrSameReviewer) {
		t.Fatal("expected distinct reviewer guard")
	}
}
func TestMigrationPreservesSubmissionAndTimeline(t *testing.T) {
	submitted := time.Date(2026, 8, 1, 8, 0, 0, 0, time.UTC)
	prior := Event{At: submitted, Action: "submitted"}
	f := Feedback{ID: "f", AreaID: "a", SubmittedAt: submitted, Timeline: []Event{prior}}
	m := approved()
	moved, err := m.Execute([]Feedback{f}, Area{ID: "a", Environment: "open"}, Area{ID: "b", Environment: "open"}, "operator", submitted.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if !moved[0].SubmittedAt.Equal(submitted) || len(moved[0].Timeline) != 2 || moved[0].Timeline[0] != prior {
		t.Fatalf("history lost: %+v", moved[0])
	}
}
func TestUndoUsesElapsedDurationAcrossTimeZones(t *testing.T) {
	at := time.Date(2026, 8, 17, 9, 0, 0, 0, time.FixedZone("CST", 8*3600))
	m := approved()
	m.ExecutedAt = &at
	items := []Feedback{{ID: "f", AreaID: "b"}}
	now := at.In(time.UTC).Add(20 * time.Minute)
	back, err := m.Undo(items, "operator", now, 30*time.Minute)
	if err != nil || back[0].AreaID != "a" {
		t.Fatalf("%v %+v", err, back)
	}
}
func TestStatsFollowCurrentAreaOnly(t *testing.T) {
	items := []Feedback{{AreaID: "a", Status: "open", Level: 5}, {AreaID: "b", Status: "resolved", Level: 2}}
	a := Recalculate("a", items)
	b := Recalculate("b", items)
	if a.Total != 1 || a.Open != 1 || a.Severe != 1 || b.Total != 1 || b.Open != 0 {
		t.Fatalf("a=%+v b=%+v", a, b)
	}
}

func TestIdempotentMigrationPipelineDomain(t *testing.T) {
	now := time.Now()
	feedback := Feedback{ID: "f", AreaID: "b", Version: 1, Timeline: []Event{{At: now, Action: "area_migrated"}}}
	m := approved()
	moved, err := m.Execute([]Feedback{feedback}, Area{ID: "a", Environment: "open"}, Area{ID: "b", Environment: "open"}, "operator", now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if moved[0].Version != 1 || len(moved[0].Timeline) != 1 {
		t.Fatalf("already migrated feedback changed again: %+v", moved[0])
	}
}
