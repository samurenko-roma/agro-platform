package production

import (
	"context"
	"sync"

	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/allocation"
	repo "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type allocationRepository struct {
	mu    sync.RWMutex
	items map[vo.ID]*allocation.Allocation
}

func NewAllocationRepository() repo.AllocationRepository {
	return &allocationRepository{items: make(map[vo.ID]*allocation.Allocation)}
}

func (r *allocationRepository) ListActiveByProductionUnitID(ctx context.Context, productionUnitID vo.ID) ([]*allocation.Allocation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*allocation.Allocation, 0)
	for _, a := range r.items {
		if a.ProductionUnitID == productionUnitID && a.EndedAt == nil {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *allocationRepository) Save(ctx context.Context, a *allocation.Allocation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[a.ID] = a
	return nil
}

func (r *allocationRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*allocation.Allocation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.items[id]
	if !ok || a.FarmID != farmID {
		return nil, nil
	}
	return a, nil
}

func (r *allocationRepository) ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*allocation.Allocation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*allocation.Allocation, 0)
	for _, a := range r.items {
		if a.CycleID == cycleID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *allocationRepository) ListByProductionUnitID(ctx context.Context, productionUnitID vo.ID) ([]*allocation.Allocation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*allocation.Allocation, 0)
	for _, a := range r.items {
		if a.ProductionUnitID == productionUnitID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *allocationRepository) Delete(ctx context.Context, id vo.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}

var _ repo.AllocationRepository = (*allocationRepository)(nil)
