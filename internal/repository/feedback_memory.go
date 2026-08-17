package repository

import (
	"context"
	"github.com/wyw14/cry039/internal/domain"
	"sync"
)

type FeedbackMemory struct {
	mu          sync.RWMutex
	byMigration map[string][]domain.Feedback
}

func NewFeedbackMemory(id string, seed []domain.Feedback) *FeedbackMemory {
	return &FeedbackMemory{byMigration: map[string][]domain.Feedback{id: append([]domain.Feedback(nil), seed...)}}
}
func (m *FeedbackMemory) LoadForMigration(_ context.Context, id string) ([]domain.Feedback, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := append([]domain.Feedback(nil), m.byMigration[id]...)
	for i := range items {
		if n := len(items[i].Timeline); n > 0 {
			items[i].Timeline = items[i].Timeline[n-1:]
		}
	}
	return items, nil
}
func (m *FeedbackMemory) ReplaceMigrationSet(_ context.Context, id string, items []domain.Feedback) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	stored := append([]domain.Feedback(nil), items...)
	for i := range stored {
		if n := len(stored[i].Timeline); n > 0 {
			stored[i].Timeline = stored[i].Timeline[n-1:]
		}
	}
	m.byMigration[id] = stored
	return nil
}
