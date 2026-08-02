package production

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/harvest_batch"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type harvestBatchRepository struct {
}

func (h harvestBatchRepository) Save(ctx context.Context, batch *harvestbatch.HarvestBatch) error {
	//TODO implement me
	panic("implement me")
}

func (h harvestBatchRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*harvestbatch.HarvestBatch, error) {
	panic("implement me")
}

func (h harvestBatchRepository) ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*harvestbatch.HarvestBatch, error) {
	//TODO implement me
	panic("implement me")
}

func (h harvestBatchRepository) Delete(ctx context.Context, id vo.ID) error {
	//TODO implement me
	panic("implement me")
}

func NewHarvestBatchRepository() repository.HarvestBatchRepository {
	return &harvestBatchRepository{}
}
