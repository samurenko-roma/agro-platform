package production

import (
	"context"
	"sync"

	harvestbatch "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/harvest_batch"
	repo "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type harvestBatchRepository struct {
	mu    sync.RWMutex
	items map[vo.ID]*harvestbatch.HarvestBatch
}

func NewHarvestBatchRepository() repo.HarvestBatchRepository {
	return &harvestBatchRepository{items: make(map[vo.ID]*harvestbatch.HarvestBatch)}
}

func (r *harvestBatchRepository) Save(ctx context.Context, batch *harvestbatch.HarvestBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[batch.ID] = batch
	return nil
}

func (r *harvestBatchRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*harvestbatch.HarvestBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.items[id]
	if !ok || b.FarmID != farmID {
		return nil, nil
	}
	return b, nil
}

func (r *harvestBatchRepository) ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*harvestbatch.HarvestBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*harvestbatch.HarvestBatch, 0)
	for _, b := range r.items {
		if b.CycleID == cycleID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (r *harvestBatchRepository) Delete(ctx context.Context, id vo.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}

var _ repo.HarvestBatchRepository = (*harvestBatchRepository)(nil)
