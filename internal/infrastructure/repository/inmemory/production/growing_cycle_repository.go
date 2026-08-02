package production

import (
	"context"
	"sync"

	gc "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"

	repo "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type growingCycleRepository struct {
	mu sync.RWMutex

	items map[vo.ID]*gc.GrowingCycle
}

func (g growingCycleRepository) Save(ctx context.Context, cycle *gc.GrowingCycle) error {
	//TODO implement me
	panic("implement me")
}

func (g growingCycleRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*gc.GrowingCycle, error) {
	panic("implement me")
}
func (g growingCycleRepository) GetByCode(ctx context.Context, code string, farmID vo.ID) (*gc.GrowingCycle, error) {
	panic("implement me")
}

func (g growingCycleRepository) List(ctx context.Context, filter repo.ListFilter) ([]*gc.GrowingCycle, error) {
	//TODO implement me
	panic("implement me")
}

func (g growingCycleRepository) Delete(ctx context.Context, id vo.ID) error {
	//TODO implement me
	panic("implement me")
}

func NewGrowingCycleRepository() repo.GrowingCycleRepository {
	return &growingCycleRepository{
		items: make(map[vo.ID]*gc.GrowingCycle),
	}
}

var _ repo.GrowingCycleRepository = (*growingCycleRepository)(nil)
