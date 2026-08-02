package planting

import (
	"time"

	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/aggregate"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type Planting struct {
	ev.BaseAggregate
	ID        vo.ID
	FarmID    vo.ID
	CycleID   vo.ID
	PlantedAt time.Time
	Quantity  float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(farmID, cycleID vo.ID, plantedAt time.Time, quantity float64) *Planting {
	now := time.Now()

	return &Planting{
		ID:        vo.NewID(),
		FarmID:    farmID,
		CycleID:   cycleID,
		PlantedAt: plantedAt,
		Quantity:  quantity,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
