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
	return append([]domain.Feedback(nil), m.byMigration[id]...), nil
}
func (m *FeedbackMemory) ReplaceMigrationSet(_ context.Context, id string, items []domain.Feedback) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.byMigration[id] = append([]domain.Feedback(nil), items...)
	return nil
}

func (m *FeedbackMemory) Stats(id, areaID string) domain.AreaStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return domain.Recalculate(id, m.byMigration[id])
}
