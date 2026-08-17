package domain

import (
	"errors"
	"sort"
	"time"
)

var (
	ErrSameReviewer     = errors.New("two distinct reviewers are required")
	ErrReviewIncomplete = errors.New("migration needs two approvals")
	ErrUndoExpired      = errors.New("migration undo window expired")
	ErrAreaMismatch     = errors.New("target area is not compatible")
)

type Area struct {
	ID, Floor, Environment string
	Capacity               int
	Archived               bool
}
type Feedback struct {
	ID, AreaID, Tag, Status, Remark string
	Level                           int
	SubmittedAt                     time.Time
	Timeline                        []Event
	Version                         uint64
}
type Event struct {
	At                              time.Time
	Actor, Action, FromArea, ToArea string
}
type Migration struct {
	ID, SourceArea, TargetArea, Filter string
	Reviewers                          []string
	AffectedIDs                        []string
	ExecutedAt                         *time.Time
	UndoneAt                           *time.Time
	Digest                             string
}

func (m *Migration) Approve(reviewer string) error {
	for _, r := range m.Reviewers {
		if r == reviewer {
			return ErrSameReviewer
		}
	}
	m.Reviewers = append(m.Reviewers, reviewer)
	return nil
}
func ValidateAreas(source, target Area) error {
	if source.ID == target.ID || target.Archived || source.Environment != target.Environment {
		return ErrAreaMismatch
	}
	return nil
}
func (m *Migration) Execute(items []Feedback, source, target Area, actor string, now time.Time) ([]Feedback, error) {
	if len(m.Reviewers) < 2 {
		return nil, ErrReviewIncomplete
	}
	if err := ValidateAreas(source, target); err != nil {
		return nil, err
	}
	out := make([]Feedback, len(items))
	var affected []Feedback
	for i, f := range items {
		out[i] = f
		if f.AreaID != source.ID {
			continue
		}
		out[i].AreaID = target.ID
		out[i].Version++
		out[i].Timeline = append(append([]Event(nil), f.Timeline...), Event{At: now, Actor: actor, Action: "area_migrated", FromArea: source.ID, ToArea: target.ID})
		affected = append(affected, out[i])
	}
	m.AffectedIDs = IDs(affected)
	m.ExecutedAt = &now
	return out, nil
}
func (m *Migration) Undo(items []Feedback, actor string, now time.Time, window time.Duration) ([]Feedback, error) {
	if m.ExecutedAt == nil {
		return nil, ErrUndoExpired
	}
	if now.Sub(*m.ExecutedAt) > window {
		return nil, ErrUndoExpired
	}
	affected := map[string]struct{}{}
	for _, id := range m.AffectedIDs {
		affected[id] = struct{}{}
	}
	out := make([]Feedback, len(items))
	for i, f := range items {
		out[i] = f
		if f.AreaID != m.TargetArea {
			continue
		}
		if len(affected) > 0 {
			if _, ok := affected[f.ID]; !ok {
				continue
			}
		}
		out[i].AreaID = m.SourceArea
		out[i].Version++
		out[i].Timeline = append(append([]Event(nil), f.Timeline...), Event{At: now, Actor: actor, Action: "migration_undone", FromArea: m.TargetArea, ToArea: m.SourceArea})
	}
	m.UndoneAt = &now
	return out, nil
}

type AreaStats struct {
	AreaID              string
	Total, Open, Severe int
}

func Recalculate(areaID string, items []Feedback) AreaStats {
	s := AreaStats{AreaID: areaID}
	for _, f := range items {
		if f.AreaID != areaID {
			continue
		}
		s.Total++
		if f.Status != "resolved" {
			s.Open++
		}
		if f.Level >= 4 {
			s.Severe++
		}
	}
	return s
}
func IDs(items []Feedback) []string {
	out := make([]string, len(items))
	for i, f := range items {
		out[i] = f.ID
	}
	sort.Strings(out)
	return out
}
