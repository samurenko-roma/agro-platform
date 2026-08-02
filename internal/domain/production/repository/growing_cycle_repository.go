package repository

import (
	"context"

	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type GrowingCycleRepository interface {
	Save(ctx context.Context, cycle *growingcycle.GrowingCycle) error
	GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*growingcycle.GrowingCycle, error)
	GetByCode(ctx context.Context, code string, farmID vo.ID) (*growingcycle.GrowingCycle, error)
	List(ctx context.Context, filter ListFilter) ([]*growingcycle.GrowingCycle, error)
	Delete(ctx context.Context, id vo.ID) error
}

type ListFilter struct {
	FarmID *vo.ID
	CropID *vo.ID
	Status *growingcycle.CycleStatus
	Limit  int
	Offset int
	Code   *string
	Search *string
}
