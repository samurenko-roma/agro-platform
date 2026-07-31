package spatial

import (
	"context"
	"sync"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	ls "github.com/samurenkoroma/agro-platform/internal/domain/spatial/entity/layout_snapshot"
	"github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
)

type layoutSnapshotRepository struct {
	mu    sync.RWMutex
	items map[vo.ID]*ls.Aggregate
}

func NewLayoutSnapshotRepository() repository.LayoutSnapshotRepository {
	return &layoutSnapshotRepository{items: make(map[vo.ID]*ls.Aggregate)}
}

func (r *layoutSnapshotRepository) Save(ctx context.Context, snapshot *ls.Aggregate) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[snapshot.Root.ID] = snapshot
	return nil
}

func (r *layoutSnapshotRepository) Get(ctx context.Context, id vo.ID) (*ls.Aggregate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.items[id], nil
}

func (r *layoutSnapshotRepository) GetLatest(ctx context.Context, farmID vo.ID) (*ls.Aggregate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var latest *ls.Aggregate
	for _, s := range r.items {
		if s.Root.FarmID != farmID {
			continue
		}
		if latest == nil || s.Root.Version > latest.Root.Version {
			latest = s
		}
	}
	return latest, nil
}

var _ repository.LayoutSnapshotRepository = (*layoutSnapshotRepository)(nil)
