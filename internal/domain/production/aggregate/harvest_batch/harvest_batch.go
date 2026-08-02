package harvestbatch

import (
	"time"

	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/aggregate"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type HarvestBatch struct {
	ev.BaseAggregate
	ID          vo.ID
	FarmID      vo.ID
	CycleID     vo.ID
	HarvestedAt time.Time
	Quantity    float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func New(farmID, cycleID vo.ID, harvestedAt time.Time, quantity float64) *HarvestBatch {
	now := time.Now()

	return &HarvestBatch{
		ID:          vo.NewID(),
		FarmID:      farmID,
		CycleID:     cycleID,
		HarvestedAt: harvestedAt,
		Quantity:    quantity,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
