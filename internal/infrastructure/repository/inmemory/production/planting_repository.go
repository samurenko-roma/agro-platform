package production

import (
	"context"
	"sync"

	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/planting"
	repo "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type plantingRepository struct {
	mu sync.RWMutex
}

func (p plantingRepository) Save(ctx context.Context, planting *planting.Planting) error {
	//TODO implement me
	panic("implement me")
}

func (p plantingRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*planting.Planting, error) {
	panic("implement me")
}
func (p plantingRepository) ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*planting.Planting, error) {
	//TODO implement me
	panic("implement me")
}

func (p plantingRepository) Delete(ctx context.Context, id vo.ID) error {
	//TODO implement me
	panic("implement me")
}

func NewPlantingRepository() repo.PlantingRepository {
	return &plantingRepository{}
}

var _ repo.PlantingRepository = (*plantingRepository)(nil)
