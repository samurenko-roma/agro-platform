package production

import (
	"context"
	"sync"

	gc "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	repo "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type growingCycleRepository struct {
	mu    sync.RWMutex
	items map[vo.ID]*gc.GrowingCycle
}

func NewGrowingCycleRepository() repo.GrowingCycleRepository {
	return &growingCycleRepository{items: make(map[vo.ID]*gc.GrowingCycle)}
}

func (r *growingCycleRepository) Save(ctx context.Context, cycle *gc.GrowingCycle) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[cycle.ID] = cycle
	return nil
}

func (r *growingCycleRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*gc.GrowingCycle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.items[id]
	if !ok || c.FarmID != farmID {
		return nil, nil
	}
	return c, nil
}

func (r *growingCycleRepository) GetByCode(ctx context.Context, code string, farmID vo.ID) (*gc.GrowingCycle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.items {
		if c.Code == code && c.FarmID == farmID {
			return c, nil
		}
	}
	return nil, nil
}

func (r *growingCycleRepository) List(ctx context.Context, filter repo.ListFilter) ([]*gc.GrowingCycle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*gc.GrowingCycle, 0)
	for _, c := range r.items {
		if filter.FarmID != nil && c.FarmID != *filter.FarmID {
			continue
		}
		if filter.CropID != nil && c.CropID != *filter.CropID {
			continue
		}
		if filter.Status != nil && c.Status != *filter.Status {
			continue
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *growingCycleRepository) Delete(ctx context.Context, id vo.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}

var _ repo.GrowingCycleRepository = (*growingCycleRepository)(nil)
