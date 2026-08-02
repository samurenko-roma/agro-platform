package production

import (
	"context"
	"sync"

	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/allocation"
	repo "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type allocationRepository struct {
	mu sync.RWMutex
}

func (a allocationRepository) ListActiveByProductionUnitID(ctx context.Context, productionUnitID vo.ID) ([]*allocation.Allocation, error) {
	//TODO implement me
	panic("implement me")
}

func (a allocationRepository) Save(ctx context.Context, allocation *allocation.Allocation) error {
	//TODO implement me
	panic("implement me")
}

func (a allocationRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*allocation.Allocation, error) {
	panic("implement me")
}

func (a allocationRepository) ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*allocation.Allocation, error) {
	//TODO implement me
	panic("implement me")
}

func (a allocationRepository) ListByProductionUnitID(ctx context.Context, productionUnitID vo.ID) ([]*allocation.Allocation, error) {
	//TODO implement me
	panic("implement me")
}

func (a allocationRepository) Delete(ctx context.Context, id vo.ID) error {
	//TODO implement me
	panic("implement me")
}

func NewAllocationRepository() repo.AllocationRepository {
	return &allocationRepository{}
}

var _ repo.AllocationRepository = (*allocationRepository)(nil)
