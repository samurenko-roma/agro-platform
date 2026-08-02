package repository

import (
	"context"

	harvestbatch "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/harvest_batch"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type HarvestBatchRepository interface {
	Save(ctx context.Context, batch *harvestbatch.HarvestBatch) error
	GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*harvestbatch.HarvestBatch, error)
	ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*harvestbatch.HarvestBatch, error)
	Delete(ctx context.Context, id vo.ID) error
}
