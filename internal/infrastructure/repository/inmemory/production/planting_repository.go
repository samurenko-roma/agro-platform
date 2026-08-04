package production

import (
	"context"
	"sync"

	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/planting"
	repo "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type plantingRepository struct {
	mu    sync.RWMutex
	items map[vo.ID]*planting.Planting
}

func NewPlantingRepository() repo.PlantingRepository {
	return &plantingRepository{items: make(map[vo.ID]*planting.Planting)}
}

func (r *plantingRepository) Save(ctx context.Context, p *planting.Planting) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[p.ID] = p
	return nil
}

func (r *plantingRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*planting.Planting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.items[id]
	if !ok || p.FarmID != farmID {
		return nil, nil
	}
	return p, nil
}

func (r *plantingRepository) ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*planting.Planting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*planting.Planting, 0)
	for _, p := range r.items {
		if p.CycleID == cycleID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *plantingRepository) Delete(ctx context.Context, id vo.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}

var _ repo.PlantingRepository = (*plantingRepository)(nil)
