package allocation

import (
	"time"

	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/event"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

const (
	EventAllocated = "production.allocation.allocated"
	EventReleased  = "production.allocation.released"
)

type AllocationAllocated struct {
	ev.BaseEvent
	ProductionUnitID vo.ID
}

func NewAllocationAllocated(allocationID vo.ID, prodUnitID vo.ID) AllocationAllocated {
	return AllocationAllocated{BaseEvent: ev.NewBaseEvent(allocationID, EventAllocated), ProductionUnitID: prodUnitID}
}

type AllocationReleased struct {
	ev.BaseEvent
	ProductionUnitID vo.ID
	ReleasedAt       time.Time
}

func NewAllocationReleased(allocationID vo.ID, prodUnitID vo.ID) AllocationReleased {
	return AllocationReleased{
		BaseEvent:        ev.NewBaseEvent(allocationID, EventReleased),
		ProductionUnitID: prodUnitID,
		ReleasedAt:       time.Now(),
	}
}
