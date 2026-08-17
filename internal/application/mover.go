package application

import (
	"context"
	"errors"
	"github.com/wyw14/cry039/internal/domain"
	"sync"
	"time"
)

var ErrMigrationKeyConflict = errors.New("migration idempotency key conflict")

type FeedbackStore interface {
	LoadForMigration(context.Context, string) ([]domain.Feedback, error)
	ReplaceMigrationSet(context.Context, string, []domain.Feedback) error
}
type Mover struct {
	store FeedbackStore
	clock func() time.Time
	mu    sync.Mutex
	keys  map[string]string
}

func NewMover(s FeedbackStore, clock func() time.Time) *Mover {
	return &Mover{store: s, clock: clock, keys: map[string]string{}}
}
func (m *Mover) Execute(ctx context.Context, key, digest string, job *domain.Migration, source, target domain.Area, actor string) ([]domain.Feedback, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !job.HasReviewQuorum() {
		return nil, domain.ErrReviewIncomplete
	}
	if prior, ok := m.keys[key]; ok {
		if prior != digest {
			return nil, ErrMigrationKeyConflict
		}
		return m.store.LoadForMigration(ctx, job.ID)
	}
	items, err := m.store.LoadForMigration(ctx, job.ID)
	if err != nil {
		return nil, err
	}
	moved, err := job.Execute(items, source, target, actor, m.clock())
	if err != nil {
		return nil, err
	}
	if err = m.store.ReplaceMigrationSet(ctx, job.ID, moved); err != nil {
		return nil, err
	}
	m.keys[key] = digest
	return moved, nil
}
