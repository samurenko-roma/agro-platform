package allocation

import (
	"time"

	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/aggregate"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type Allocation struct {
	ev.BaseAggregate
	ID               vo.ID
	FarmID           vo.ID
	CycleID          vo.ID
	ProductionUnitID vo.ID
	Area             float64
	StartedAt        *time.Time
	EndedAt          *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func New(farmID, cycleID, productionUnitID vo.ID, area float64, startedAt *time.Time) *Allocation {
	now := time.Now()
	a := &Allocation{
		ID:               vo.NewID(),
		FarmID:           farmID,
		CycleID:          cycleID,
		ProductionUnitID: productionUnitID,
		Area:             area,
		StartedAt:        startedAt,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	a.AddEvent(NewAllocationAllocated(a.ID, productionUnitID))
	return a
}

// Release завершает размещение цикла на физическом узле. Эмитит событие,
// на которое реагирует Spatial (освобождает узел через event handler).
func (a *Allocation) Release(releasedAt time.Time) error {
	if a.EndedAt != nil {
		return ErrAlreadyReleased
	}
	a.EndedAt = &releasedAt
	a.UpdatedAt = time.Now()
	a.AddEvent(NewAllocationReleased(a.ID, a.ProductionUnitID))
	return nil
}

// Reallocate переносит размещение на другой физический узел (или просто
// меняет площадь/даты на том же узле). Если узел действительно меняется —
// эмитит освобождение старого и занятие нового, оба события уже обработаны
// существующим event handler'ом в Spatial.
func (a *Allocation) Reallocate(newProductionUnitID vo.ID, newArea float64, startedAt, endedAt *time.Time) {
	oldProductionUnitID := a.ProductionUnitID

	a.ProductionUnitID = newProductionUnitID
	a.Area = newArea
	a.StartedAt = startedAt
	a.EndedAt = endedAt
	a.UpdatedAt = time.Now()

	if oldProductionUnitID != newProductionUnitID {
		a.AddEvent(NewAllocationReleased(a.ID, oldProductionUnitID))
		a.AddEvent(NewAllocationAllocated(a.ID, newProductionUnitID))
	}
}
