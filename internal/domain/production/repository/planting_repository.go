package repository

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/planting"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type PlantingRepository interface {
	Save(ctx context.Context, planting *planting.Planting) error
	GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*planting.Planting, error)
	ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*planting.Planting, error)
	Delete(ctx context.Context, id vo.ID) error
}
