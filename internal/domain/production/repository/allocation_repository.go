package repository

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/allocation"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type AllocationRepository interface {
	Save(ctx context.Context, allocation *allocation.Allocation) error
	GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*allocation.Allocation, error)
	ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*allocation.Allocation, error)
	ListByProductionUnitID(ctx context.Context, productionUnitID vo.ID) ([]*allocation.Allocation, error)
	ListActiveByProductionUnitID(ctx context.Context, productionUnitID vo.ID) ([]*allocation.Allocation, error)
	Delete(ctx context.Context, id vo.ID) error
}
