package spatial

import (
	"context"
	"sync"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
	"github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
)

type productionUnitRepository struct {
	mu    sync.RWMutex
	items map[vo.ID]*pu.ProductionUnit
}

func NewProductionUnitRepository() repository.ProductionUnitRepository {
	return &productionUnitRepository{items: make(map[vo.ID]*pu.ProductionUnit)}
}

func (r *productionUnitRepository) Save(ctx context.Context, unit *pu.ProductionUnit) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[unit.ID] = unit
	return nil
}

func (r *productionUnitRepository) GetByID(ctx context.Context, id vo.ID, orgId vo.ID) (*pu.ProductionUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	unit, ok := r.items[id]
	if !ok || unit.OwnerID != orgId {
		return nil, nil
	}
	return unit, nil
}

func (r *productionUnitRepository) GetChildren(ctx context.Context, parentID vo.ID) ([]*pu.ProductionUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*pu.ProductionUnit, 0)
	for _, unit := range r.items {
		if unit.ParentID != nil && *unit.ParentID == parentID {
			result = append(result, unit)
		}
	}
	return result, nil
}

func (r *productionUnitRepository) ListByOwner(ctx context.Context, ownerID vo.ID) ([]*pu.ProductionUnit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*pu.ProductionUnit, 0)
	for _, unit := range r.items {
		if unit.OwnerID == ownerID {
			result = append(result, unit)
		}
	}
	return result, nil
}

func (r *productionUnitRepository) Exists(ctx context.Context, id vo.ID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.items[id]
	return ok, nil
}

func (r *productionUnitRepository) GetNextSequence(ctx context.Context, orgID vo.ID, parentID *vo.ID, unitType pu.ProductionUnitType) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	max := 0
	for _, unit := range r.items {
		if unit.OwnerID != orgID || unit.Type != unitType {
			continue
		}
		sameParent := (unit.ParentID == nil && parentID == nil) ||
			(unit.ParentID != nil && parentID != nil && *unit.ParentID == *parentID)
		if !sameParent {
			continue
		}
		if unit.Sequence > max {
			max = unit.Sequence
		}
	}
	return max + 1, nil
}

var _ repository.ProductionUnitRepository = (*productionUnitRepository)(nil)
