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
		items[i].Timeline = append([]domain.Event(nil), items[i].Timeline...)
	}
	return items, nil
}
func (m *FeedbackMemory) ReplaceMigrationSet(_ context.Context, id string, items []domain.Feedback) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	stored := append([]domain.Feedback(nil), items...)
	for i := range stored {
		stored[i].Timeline = append([]domain.Event(nil), stored[i].Timeline...)
	}
	m.byMigration[id] = stored
	return nil
}
